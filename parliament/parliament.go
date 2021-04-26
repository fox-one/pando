package parliament

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pando/service/asset"
)

func New(
	messages core.MessageStore,
	userz core.UserService,
	assetz core.AssetService,
	walletz core.WalletService,
	collaterals core.CollateralStore,
	system *core.System,
) core.Parliament {
	return &parliament{
		messages:    messages,
		userz:       userz,
		assetz:      asset.Cache(assetz),
		walletz:     walletz,
		collaterals: collaterals,
		system:      system,
	}
}

type parliament struct {
	messages    core.MessageStore
	userz       core.UserService
	assetz      core.AssetService
	walletz     core.WalletService
	collaterals core.CollateralStore
	system      *core.System
}

func (s *parliament) requestVoteAction(ctx context.Context, proposal *core.Proposal) (string, error) {
	id, _ := uuid.FromString(proposal.TraceID)
	body, err := mtg.Encode(core.ActionProposalVote, id)
	if err != nil {
		return "", err
	}

	data, err := core.TransactionAction{
		Body: body,
	}.Encode()
	if err != nil {
		return "", err
	}

	memo, err := mtg.Encrypt(data, mixin.GenerateEd25519Key(), s.system.PublicKey)
	if err != nil {
		return "", err
	}

	transfer := &core.Transfer{
		TraceID:   uuid.Modify(proposal.TraceID, s.system.ClientID),
		AssetID:   s.system.GasAssetID,
		Amount:    s.system.GasAmount,
		Memo:      base64.StdEncoding.EncodeToString(memo),
		Threshold: s.system.Threshold,
		Opponents: s.system.Members,
	}

	code, err := s.walletz.ReqTransfer(ctx, transfer)
	if err != nil {
		return "", err
	}

	return paymentAction(code), nil
}

func (s *parliament) ProposalCreated(ctx context.Context, p *core.Proposal) error {
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
				Value:  s.fetchUserName(ctx, p.Creator),
				Action: userAction(p.Creator),
			},
			{
				Key:    "pay",
				Value:  fmt.Sprintf("%s %s", number.Humanize(p.Amount), s.fetchAssetSymbol(ctx, p.AssetID)),
				Action: assetAction(p.AssetID),
			},
		},
	}

	data, _ := base64.StdEncoding.DecodeString(p.Data)
	view.Meta = s.renderProposalItems(ctx, p.Action, data)

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
	post := execute("proposal_created", view)

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

func (s *parliament) ProposalApproved(ctx context.Context, p *core.Proposal) error {
	by := p.Votes[len(p.Votes)-1]

	view := Proposal{
		ApprovedCount: len(p.Votes),
		ApprovedBy:    s.fetchUserName(ctx, by),
	}

	post := execute("proposal_approved", view)

	var messages []*core.Message
	for _, admin := range s.system.Admins {
		quote := uuid.Modify(p.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "ProposalApproved By "+by),
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

	post := execute("proposal_passed", nil)

	for _, admin := range s.system.Admins {
		quote := uuid.Modify(proposal.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "ProposalPassed"),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}

func (s *parliament) FlipCreated(ctx context.Context, flip *core.Flip) error {
	gem, dai := s.fetchCatGemDai(ctx, flip.CollateralID)

	view := Flip{
		Number: flip.ID,
		Info: []Item{
			{
				Key:   "id",
				Value: flip.TraceID,
			},
			{
				Key:   "vault",
				Value: flip.VaultID,
			},
			{
				Key:   "lot",
				Value: fmt.Sprintf("%s %s", number.Humanize(flip.Lot), gem),
			},
			{
				Key:   "tab",
				Value: fmt.Sprintf("%s %s", number.Humanize(flip.Tab), dai),
			},
			{
				Key:    "kicker",
				Value:  s.fetchUserName(ctx, flip.Guy),
				Action: userAction(flip.Guy),
			},
		},
	}

	buttons := generateButtons(view.Info)
	if len(buttons) > 6 {
		buttons = buttons[len(buttons)-6:]
	}

	buttonsData, _ := json.Marshal(buttons)
	post := execute("flip_create", view)

	var messages []*core.Message
	for _, admin := range s.system.Admins {
		// post
		postMsg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(flip.TraceID, s.system.ClientID+admin),
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

func (s *parliament) buildFlipStat(ctx context.Context, flip *core.Flip) FlipStat {
	gem, dai := s.fetchCatGemDai(ctx, flip.CollateralID)

	return FlipStat{
		Lot: number.Humanize(flip.Lot),
		Bid: number.Humanize(flip.Bid),
		Tab: number.Humanize(flip.Tab),
		Gem: gem,
		Dai: dai,
	}
}

func (s *parliament) FlipBid(ctx context.Context, flip *core.Flip, _ *core.FlipEvent) error {
	var messages []*core.Message

	stat := s.buildFlipStat(ctx, flip)
	post := execute("flip_bid", stat)

	for _, admin := range s.system.Admins {
		quote := uuid.Modify(flip.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "FlipBid"),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}

func (s *parliament) FlipDeal(ctx context.Context, flip *core.Flip) error {
	var messages []*core.Message

	stat := s.buildFlipStat(ctx, flip)
	post := execute("flip_deal", stat)

	for _, admin := range s.system.Admins {
		quote := uuid.Modify(flip.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "FlipDeal"),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}
