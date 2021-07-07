package notifier

import (
	"context"
	"encoding/base64"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
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

	switch tx.Action {
	case core.ActionFlipKick:
		msg := n.localize("vat_kicked", user.Lang, args)
		req := &mixin.MessageRequest{
			ConversationID: mixin.UniqueConversationID(n.system.ClientID, vat.UserID),
			RecipientID:    vat.UserID,
			MessageID:      uuid.Modify(vat.TraceID, "kick by "+tx.TraceID),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString([]byte(msg)),
		}

		if err := n.messages.Create(ctx, []*core.Message{core.BuildMessage(req)}); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("messages.Create")
			return err
		}
	case core.ActionFlipDeal:
		data.cc(flip.Guy)
	}

	return nil
}
