package notifier

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/big"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/number"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

func mixinRawTransactionTraceId(hash string, index uint8) string {
	h := md5.New()
	_, _ = io.WriteString(h, hash)
	b := new(big.Int).SetInt64(int64(index))
	h.Write(b.Bytes())
	s := h.Sum(nil)
	s[6] = (s[6] & 0x0f) | 0x30
	s[8] = (s[8] & 0x3f) | 0x80
	sid, err := uuid.FromBytes(s)
	if err != nil {
		panic(err)
	}

	return sid.String()
}

func getDebt(cat *core.Collateral, vat *core.Vault) decimal.Decimal {
	return number.Ceil(cat.Rate.Mul(vat.Art), 8)
}

func getCollateralRate(cat *core.Collateral, vat *core.Vault) string {
	ink := vat.Ink
	debt := getDebt(cat, vat)

	rate := decimal.Zero
	if debt.IsPositive() {
		rate = ink.Mul(cat.Price).Div(debt).Truncate(4)
	}

	return fmt.Sprintf(`%s%`, rate.Shift(2))
}
