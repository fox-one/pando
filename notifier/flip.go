package notifier

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/internal/color"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/spf13/cast"
)

func (n *notifier) handleFlipTx(ctx context.Context, tx *core.Transaction, user *core.User, data *TxData) error {
	flipID := tx.TraceID

	if tx.Action != core.ActionFlipKick {
		var parameters []interface{}
		_ = tx.Parameters.Unmarshal(&parameters)

		if len(parameters) > 0 {
			flipID = cast.ToString(parameters[0])
		}
	}

	flip, err := n.flips.Find(ctx, flipID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("flips.Find")
		return err
	}

	if flip.ID == 0 {
		return nil
	}

	vat, err := n.vats.Find(ctx, flip.VaultID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("vats.Find")
		return err
	}

	cat, err := n.cats.Find(ctx, vat.CollateralID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("cats.Find")
		return err
	}

	event, err := n.flips.FindEvent(ctx, flipID, tx.Version)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("flips.FindEvent")
		return err
	}

	gem, err := n.assetz.Find(ctx, cat.Gem)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("assetz.Find")
		return err
	}

	dai, err := n.assetz.Find(ctx, cat.Dai)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("assetz.Find")
		return err
	}

	args := map[string]interface{}{
		"Name":    cat.Name,
		"VaultID": vat.ID,
		"Lot":     number.Humanize(event.Lot),
		"Bid":     number.Humanize(event.Bid),
		"Gem":     gem.Symbol,
		"Dai":     dai.Symbol,
		"Price":   number.Humanize(event.Bid.Div(event.Lot).Truncate(8)),
	}

	data.AddLine(n.localize("flip_lot", user.Lang, args))
	data.AddLine(n.localize("flip_bid", user.Lang, args))
	data.AddLine(n.localize("flip_bid_price", user.Lang, args))

	if action, err := n.executeLink("flip_detail", map[string]string{
		"flip_id": flipID,
	}); err == nil {
		label := n.localize("flip_button", user.Lang)
		data.AddButton(label, action)
	}

	return n.notifyFlipParticipant(ctx, tx, vat, flip, args)
}

func (n *notifier) notifyFlipParticipant(ctx context.Context, tx *core.Transaction, vat *core.Vault, flip *core.Flip, args interface{}) error {
	var (
		user  *core.User
		topic string
		err   error
	)

	switch tx.Action {
	case core.ActionFlipKick:
		user, err = n.users.Find(ctx, vat.UserID)
		topic = "vat_kicked"
	case core.ActionFlipDeal:
		user, err = n.users.Find(ctx, flip.Guy)
		topic = "flip_win"
	default:
		return nil
	}

	if err != nil {
		return err
	}

	var messages []*core.Message

	{
		msg := n.localize(topic, user.Lang, args)
		req := &mixin.MessageRequest{
			ConversationID: mixin.UniqueConversationID(n.system.ClientID, user.MixinID),
			RecipientID:    user.MixinID,
			MessageID:      uuid.Modify(vat.TraceID, tx.TraceID+topic),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString([]byte(msg)),
		}

		messages = append(messages, core.BuildMessage(req))
	}

	if action, err := n.executeLink("flip_detail", map[string]string{
		"flip_id": flip.TraceID,
	}); err == nil {
		label := n.localize("flip_button", user.Lang)
		buttons, _ := json.Marshal(mixin.AppButtonGroupMessage{mixin.AppButtonMessage{
			Label:  label,
			Action: action,
			Color:  color.Random(),
		}})

		req := &mixin.MessageRequest{
			ConversationID: mixin.UniqueConversationID(n.system.ClientID, user.MixinID),
			RecipientID:    user.MixinID,
			MessageID:      uuid.Modify(vat.TraceID, tx.TraceID+topic+"buttons"),
			Category:       mixin.MessageCategoryAppButtonGroup,
			Data:           base64.StdEncoding.EncodeToString(buttons),
		}

		messages = append(messages, core.BuildMessage(req))
	}

	return n.messages.Create(ctx, messages)
}
