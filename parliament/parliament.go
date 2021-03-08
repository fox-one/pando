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
	"github.com/shopspring/decimal"
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
	trace, _ := uuid.FromString(proposal.TraceID)
	uid, _ := uuid.FromString(s.system.ClientID)
	memo, err := mtg.Encode(uid, core.ActionProposalVote, trace)
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
		Opponents: s.system.Members,
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

func (s *parliament) fetchCatName(ctx context.Context, id string) string {
	c, err := s.collaterals.Find(ctx, id)
	if err != nil {
		return "NULL"
	}

	return c.Name
}

func (s *parliament) Created(ctx context.Context, p *core.Proposal) error {
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

	switch p.Action {
	case core.ActionCatCreate:
		var (
			gem, dai uuid.UUID
			name     string
		)

		_, _ = mtg.Scan(data, &gem, &dai, &name)

		view.Meta = []Item{
			{
				Key:   "name",
				Value: name,
			},
			{
				Key:    "gem",
				Value:  s.fetchAssetSymbol(ctx, gem.String()),
				Action: assetAction(gem.String()),
			},
			{
				Key:    "dai",
				Value:  s.fetchAssetSymbol(ctx, dai.String()),
				Action: assetAction(dai.String()),
			},
		}
	case core.ActionCatEdit:
		var id uuid.UUID
		data, err := mtg.Scan(data, &id)

		view.Meta = []Item{
			{
				Key:   "cat",
				Value: s.fetchCatName(ctx, id.String()),
			},
		}

		for {
			var item Item
			if data, err = mtg.Scan(data, &item.Key, &item.Value); err != nil {
				break
			}

			view.Meta = append(view.Meta, item)
		}
	case core.ActionOracleFeed:
		var (
			id    uuid.UUID
			price decimal.Decimal
			ts    int64
		)

		_, _ = mtg.Scan(data, &id, &price, &ts)

		view.Meta = []Item{
			{
				Key:    "asset",
				Value:  s.fetchAssetSymbol(ctx, id.String()),
				Action: assetAction(id.String()),
			},
			{
				Key:   "price",
				Value: number.Humanize(price),
			},
			{
				Key:   "date",
				Value: time.Unix(ts, 0).Format(time.RFC3339),
			},
		}
	case core.ActionSysWithdraw:
		var (
			asset    uuid.UUID
			amount   decimal.Decimal
			opponent uuid.UUID
		)

		_, _ = mtg.Scan(data, &asset, &amount, &opponent)

		view.Meta = []Item{
			{
				Key:    "asset",
				Value:  fmt.Sprintf("%s %s", number.Humanize(amount), s.fetchAssetSymbol(ctx, asset.String())),
				Action: assetAction(asset.String()),
			},
			{
				Key:    "opponent",
				Value:  s.fetchUserName(ctx, opponent.String()),
				Action: userAction(opponent.String()),
			},
		}
	case core.ActionFlipOpt:
		var (
			beg      decimal.Decimal
			ttl, tau int64
		)

		_, _ = mtg.Scan(data, &beg, &ttl, &tau)

		view.Meta = []Item{
			{
				Key:   "beg",
				Value: beg.String(),
			},
			{
				Key:   "ttl",
				Value: (time.Duration(ttl) * time.Second).String(),
			},
			{
				Key:   "Tau",
				Value: (time.Duration(tau) * time.Second).String(),
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

func (s *parliament) Approved(ctx context.Context, p *core.Proposal) error {
	by := p.Votes[len(p.Votes)-1]

	view := Proposal{
		ApprovedCount: len(p.Votes),
		ApprovedBy:    s.fetchUserName(ctx, by),
	}

	post := renderApprovedBy(view)

	var messages []*core.Message
	for _, admin := range s.system.Admins {
		quote := uuid.Modify(p.TraceID, s.system.ClientID+admin)
		msg := &mixin.MessageRequest{
			RecipientID:    admin,
			ConversationID: mixin.UniqueConversationID(s.system.ClientID, admin),
			MessageID:      uuid.Modify(quote, "Approved By "+by),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString(post),
			QuoteMessageID: quote,
		}

		messages = append(messages, core.BuildMessage(msg))
	}

	return s.messages.Create(ctx, messages)
}

func (s *parliament) Passed(ctx context.Context, proposal *core.Proposal) error {
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
