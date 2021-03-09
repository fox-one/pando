package notifier

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
	"github.com/spf13/cast"
)

func (n *notifier) handleVatTx(ctx context.Context, tx *core.Transaction, data *TxData) error {
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

	data.Lines = append(data.Lines, n.localize("vat_name", "Name", cat.Name))

	if event.Dink.IsPositive() {
		data.Lines = append(data.Lines, n.localize("vat_deposit", "Dink", event.Dink.String(), "Gem", gem.Symbol))
	} else if event.Dink.IsNegative() {
		data.Lines = append(data.Lines, n.localize("vat_withdraw", "Dink", event.Dink.Abs().String(), "Gem", gem.Symbol))
	}

	if event.Debt.IsPositive() {
		data.Lines = append(data.Lines, n.localize("vat_generate", "Debt", event.Debt.String(), "Dai", dai.Symbol))
	} else {
		data.Lines = append(data.Lines, n.localize("vat_payback", "Debt", event.Debt.Abs().String(), "Dai", dai.Symbol))
	}

	return nil
}
