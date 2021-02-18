package flip

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/shopspring/decimal"
)

func Kick(tx *maker.Tx, cat *core.Collateral, vault *core.Vault, opt Option, data *Data) error {
	if !cat.Live {
		return ErrFlipNotLive
	}

	if vault.Ink.Mul(cat.Price).GreaterThanOrEqual(urn.Art.Mul(cat.Rate).Mul(cat.Mat)) {
		return ErrFlipNotUnsafe
	}

	dart := decimal.Min(
		cat.Dunk.Div(cat.Rate).Div(cat.Chop),
		urn.Art,
	)

	dink := vault.Ink
	if dart.LessThan(urn.Art) {
		dink = dart.Div(urn.Art).Mul(dink).Truncate(8)
	}

	if !dart.IsPositive() || !dink.IsPositive() {
		return ErrFlipNullAuction
	}

	data.Lot = dink
	data.Dink = dink.Neg()
	data.Dart = dart.Neg()

	return nil
}

func ApplyKick(tx *maker.Tx, cat *core.Collateral, vault *core.Vault, opt Option, data Data) *core.Flip {
	// cat
	cat.Art = cat.Art.Add(data.Dart)

	// vault
	urn.Art = vault.Art.Add(data.Dart)
	urn.Ink = vault.Ink.Add(data.Dink)

	tab := data.Dart.Mul(cat.Rate).Mul(cat.Chop).Truncate(8)

	return &core.Flip{
		CreatedAt: tx.Now,
		UpdatedAt: tx.Now,
		Version:   tx.Version,
		TraceID:   tx.TraceID,
		VaultID:   vault.TraceID,
		Action:    core.ActionFlipKick,
		Tic:       time.Time{},
		End:       tx.Now.Add(opt.Tau),
		Bid:       data.Bid,
		Lot:       data.Lot,
		Tab:       tab,
		Guy:       tx.Sender,
	}
}
