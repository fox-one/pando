package param

import (
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBindingParams(t *testing.T) {
	r, _ := http.NewRequest("GET", "https://api.fox.one?symbol=BOX&amount=0.1&hide=1&foo=bar", nil)

	var params struct {
		Symbol string          `json:"symbol,omitempty"`
		Amount decimal.Decimal `json:"amount,omitempty"`
		Hide   bool            `json:"hide,omitempty"`
	}

	if err := Binding(r, &params); assert.Nil(t, err) {
		assert.Equal(t, "BOX", params.Symbol)
		assert.Equal(t, "0.1", params.Amount.String())
		assert.True(t, params.Hide)
	}
}
