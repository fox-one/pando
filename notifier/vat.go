package notifier

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/spf13/cast"
)

func (n *notifier) handleVatTx(ctx context.Context, tx *core.Transaction, user *core.User, data *TxData) error {
	vatID := tx.TraceID

	if tx.Action != core.ActionVatOpen {
		var parameters []interface{}
		_ = tx.Parameters.Unmarshal(&parameters)

		if len(parameters) > 0 {
			vatID = cast.ToString(parameters[0])
		}
	}

	vat, err := n.vats.Find(ctx, vatID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("vats.Find")
		return err
	}

	if vat.ID == 0 {
		return nil
	}

	cat, err := n.cats.Find(ctx, vat.CollateralID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("cats.Find")
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

	event, err := n.vats.FindEvent(ctx, vatID, tx.Version)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("vats.FindEvent")
		return err
	}

	args := map[string]interface{}{
		"Name": cat.Name,
		"ID":   vat.ID,
		"Dink": number.Humanize(event.Dink),
		"Debt": number.Humanize(event.Debt),
		"Gem":  gem.Symbol,
		"Dai":  dai.Symbol,
	}

	data.AddLine(n.localize("vat_name", user.Lang, args))

	if event.Dink.IsPositive() {
		data.AddLine(n.localize("vat_deposit", user.Lang, args))
	} else if event.Dink.IsNegative() {
		data.AddLine(n.localize("vat_withdraw", user.Lang, args))
	}

	if event.Debt.IsPositive() {
		data.AddLine(n.localize("vat_generate", user.Lang, args))
	} else if event.Debt.IsNegative() {
		data.AddLine(n.localize("vat_payback", user.Lang, args))
	}

	return nil
}
