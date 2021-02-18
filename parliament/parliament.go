package parliament

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/core/proposal"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pkg/uuid"
)

func New(
	messages core.MessageStore,
	userz core.UserService,
	assetz core.AssetService,
	walletz core.WalletService,
	system *core.System,
) core.Parliament {
	return &parliament{
		messages: messages,
		userz:    userz,
		assetz:   assetz,
		walletz:  walletz,
		system:   system,
	}
}

type parliament struct {
	messages core.MessageStore
	userz    core.UserService
	assetz   core.AssetService
	walletz  core.WalletService
	system   *core.System
}

func (s *parliament) requestVoteAction(ctx context.Context, proposal *core.Proposal) (string, error) {
	trace, _ := uuid.FromString(proposal.TraceID)
	uid, _ := uuid.FromString(s.system.ClientID)
	memo, err := mtg.Encode(uid, trace, int(core.ProposalActionVote))
	if err != nil {
		return "", err
	}

	sign := mtg.Sign(memo, s.system.SignKey)
	memo = mtg.Pack(memo, sign)

	transfer := &core.Transfer{
		TraceID:   uuid.Modify(proposal.TraceID, s.system.ClientID),
		AssetID:   s.system.VoteAsset,
		Amount:    s.system.VoteAmount,
		Memo:      base64.StdEncoding.EncodeToString(memo),
		Threshold: s.system.Threshold,
		Opponents: s.system.MemberIDs(),
	}

	code, err := s.walletz.ReqTransfer(ctx, transfer)
	if err != nil {
		return "", err
	}

	return paymentAction(code), nil
}

func (s *parliament) fetchAssetSymbol(ctx context.Context, assetID string) string {
	asset, err := s.assetz.Find(ctx, assetID)
	if err != nil {
		return "NULL"
	}

	return asset.Symbol
}

func (s *parliament) fetchUserName(ctx context.Context, userID string) string {
	user, err := s.userz.Find(ctx, userID)
	if err != nil {
		return "NULL"
	}

	return user.Name
}

func (s *parliament) ProposalCreated(ctx context.Context, p *core.Proposal, by *core.Member) error {
	view := Proposal{
		Number: p.ID,
		Action: p.Action.String(),
		Info: []Item{
			{
				Key:   "action",
				Value: p.Action.String(),
			},
			{
				Key:   "id",
				Value: p.TraceID,
			},
			{
				Key:   "date",
				Value: p.CreatedAt.Format(time.RFC3339),
			},
			{
				Key:    "creator",
				Value:  s.fetchUserName(ctx, by.ClientID),
				Action: userAction(by.ClientID),
			},
			{
				Key:    "asset",
				Value:  fmt.Sprintf("%s %s", p.Amount, s.fetchAssetSymbol(ctx, p.AssetID)),
				Action: assetAction(p.AssetID),
			},
		},
	}

	switch p.Action {
	case core.ProposalActionWithdraw:
		var content proposal.Withdraw
		_ = p.Content.Unmarshal(&content)

		view.Meta = []Item{
			{
				Key:    "asset",
				Value:  fmt.Sprintf("%s %s", content.Amount, s.fetchAssetSymbol(ctx, content.Asset)),
				Action: assetAction(content.Asset),
			},
			{
				Key:    "recipient",
				Value:  s.fetchUserName(ctx, content.Opponent),
				Action: userAction(content.Opponent),
			},
		}
	case core.ProposalActionSetProperty:
		var content proposal.SetProperty
		_ = p.Content.Unmarshal(&content)

		view.Meta = []Item{
			{
				Key:   "key",
				Value: content.Key,
			},
			{
				Key:   "Value",
				Value: content.Value,
			},
		}
	}

	items := append(view.Info, view.Meta...)

	voteAction, err := s.requestVoteAction(ctx, p)
	if err != nil {
		return err
	}

	items = append(items, Item{
		Key:    "Vote",
		Value:  "Vote",
		Action: voteAction,
	})

	buttons := generateButtons(items)
	if len(buttons) > 6 {
		buttons = buttons[len(buttons)-6:]
	}

	buttonsData, _ := json.Marshal(buttons)
	post := renderProposal(view)

	var messages []*core.Message
	for _, admin := range s.system.Admins {
		// post
		postMsg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(p.TraceID, s.system.ClientID+admin),
			Category:       mixin.MessageCategoryPlainPost,
			Data:           base64.StdEncoding.EncodeToString(post),
		}

		// buttons
		buttonMsg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(postMsg.MessageID, "buttons"),
			Category:       mixin.MessageCategoryAppButtonGroup,
			Data:           base64.StdEncoding.EncodeToString(buttonsData),
		}

		messages = append(messages, core.BuildMessage(postMsg), core.BuildMessage(buttonMsg))
	}

	return s.messages.Create(ctx, messages)
}

func (s *parliament) ProposalApproved(ctx context.Context, p *core.Proposal, by *core.Member) error {
	view := Proposal{
		ApprovedCount: len(p.Votes),
		ApprovedBy:    s.fetchUserName(ctx, by.ClientID),
	}

	post := renderApprovedBy(view)

	var messages []*core.Message
	for _, admin := range s.system.Admins {
		quote := uuid.Modify(p.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "Approved By "+by.ClientID),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}

func (s *parliament) ProposalPassed(ctx context.Context, proposal *core.Proposal) error {
	var messages []*core.Message

	post := []byte(passedTpl)
	for _, admin := range s.system.Admins {
		quote := uuid.Modify(proposal.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "Proposal Passed"),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}
