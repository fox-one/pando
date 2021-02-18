package vat

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

type FrobData struct {
	Dink decimal.Decimal `json:"dink,omitempty"`
	Dart decimal.Decimal `json:"dart,omitempty"`
	Debt decimal.Decimal `json:"debt,omitempty"`
}

// Frob modify a Vault
func Frob(tx *maker.Tx, cat *core.Collateral, vault *core.Vault, data *FrobData) error {
	if !cat.Live {
		return ErrVatNotLive
	}

	if tx.Sender != vault.UserID {
		return ErrVatNotAllowed
	}

	if data.Dink.IsPositive() {
		if ok := tx.AssetID == cat.Gem && tx.Amount.Equal(data.Dink); !ok {
			return ErrVatInkNotMatch
		}
	}

	if data.Debt.IsNegative() {
		if ok := tx.AssetID == cat.Dai && tx.Amount.Equal(data.Debt.Neg()); !ok {
			return ErrVatArtNotMatch
		}
	}

	data.Dart = data.Debt.Div(cat.Rate)

	totalArt := cat.Art.Add(data.Dart)
	totalDebt := cat.Rate.Mul(totalArt).Truncate(8)

	if data.Dart.IsPositive() && totalDebt.GreaterThan(cat.Line) {
		return ErrVatCeilingExceeded
	}

	ink, art := vault.Ink.Add(data.Dink), vault.Art.Add(data.Dart)
	tab := cat.Rate.Mul(art).Truncate(8)
	if ink.Mul(cat.Price).LessThan(tab.Mul(cat.Mat)) {
		return ErrVatNotSafe
	}

	if !tab.IsZero() && tab.LessThan(cat.Dust) {
		return ErrVatDust
	}

	if data.Dink.IsNegative() {
		memo := maker.EncodeMemo(module, cat.TraceID, "Withdraw")
		tx.Transfer(
			uuid.Modify(tx.TraceID, memo),
			cat.Gem,
			vault.UserID,
			memo,
			data.Dink.Abs(),
		)
	}

	if data.Debt.IsPositive() {
		memo := maker.EncodeMemo(module, cat.TraceID, "Generate")
		tx.Transfer(
			uuid.Modify(tx.TraceID, memo),
			cat.Dai,
			vault.UserID,
			memo,
			data.Debt.Abs(),
		)
	}

	return nil
}

func ApplyFrob(tx *maker.Tx, cat *core.Collateral, vault *core.Vault, data FrobData) {
	// cat
	cat.Art = cat.Art.Add(data.Dart)
	cat.Debt = cat.Debt.Add(data.Debt)

	// vault
	vault.Art = vault.Art.Add(data.Dart)
	vault.Ink = vault.Ink.Add(data.Dink)
}

func createVault(tx *maker.Tx, cat *core.Collateral) *core.Vault {
	return &core.Vault{
		CreatedAt:    tx.Now,
		TraceID:      tx.TraceID,
		UserID:       tx.Sender,
		Version:      tx.Version,
		CollateralID: cat.TraceID,
	}
}

func Open(tx *maker.Tx, cat *core.Collateral, data *FrobData) error {
	vault := createVault(tx, cat)
	data.Dink = tx.Amount
	return Frob(tx, cat, vault, data)
}

func ApplyOpen(tx *maker.Tx, cat *core.Collateral, data FrobData) *core.Vault {
	vault := createVault(tx, cat)
	ApplyFrob(tx, cat, vault, data)
	return vault
}
