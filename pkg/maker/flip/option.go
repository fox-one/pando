package flip

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/shopspring/decimal"
)

type Option struct {
	Beg decimal.Decimal `json:"beg,omitempty"`
	TTL time.Duration   `json:"ttl,omitempty"`
	Tau time.Duration   `json:"tau,omitempty"`
}

const (
	optionsKey = "flip_options"
)

func SaveOptions(ctx context.Context, properties property.Store, opt Option) error {
	b, _ := json.Marshal(opt)
	return properties.Save(ctx, optionsKey, b)
}

func ReadOptions(ctx context.Context, properties property.Store) (*Option, error) {
	v, err := properties.Get(ctx, optionsKey)
	if err != nil {
		return nil, err
	}

	opt := Option{
		Beg: number.Decimal("0.05"),
		TTL: time.Minute * 15,
		Tau: time.Hour * 3,
	}

	if b := []byte(v.String()); len(b) > 0 {
		_ = json.Unmarshal(b, &opt)
	}

	return &opt, nil
}

func HandleOpt(
	properties property.Store,
) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		var (
			beg      decimal.Decimal
			ttl, tau int64
		)

		if err := require(r.Scan(&beg, &ttl, &tau) == nil, "bad-data"); err != nil {
			return err
		}

		opt := Option{
			Beg: beg,
			TTL: time.Duration(ttl) * time.Second,
			Tau: time.Duration(tau) * time.Second,
		}

		if err := SaveOptions(ctx, properties, opt); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("SaveOptions")
			return err
		}

		return nil
	}
}
