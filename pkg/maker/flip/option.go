package flip

import (
	"time"

	"github.com/shopspring/decimal"
)

type Option struct {
	Beg decimal.Decimal `json:"beg,omitempty"`
	TTL time.Duration   `json:"ttl,omitempty"`
	Tau time.Duration   `json:"tau,omitempty"`
}
