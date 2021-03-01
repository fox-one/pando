package view

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/pkg/maker/flip"
)

func Flip(flip *core.Flip) *api.Flip {
	now := time.Now()

	return &api.Flip{
		Id:           flip.TraceID,
		CreatedAt:    Time(&flip.CreatedAt),
		Tic:          Time(&flip.Tic),
		End:          Time(&flip.End),
		Ended:        now.After(flip.Tic) || now.After(flip.End),
		Settled:      flip.Action == core.ActionFlipDeal,
		Bid:          flip.Bid.String(),
		Lot:          flip.Lot.String(),
		Tab:          flip.Tab.String(),
		CollateralId: flip.CollateralID,
		VaultId:      flip.VaultID,
		Guy:          flip.Guy,
	}
}

func FlipOption(opt *flip.Option) *api.FlipOption {
	return &api.FlipOption{
		Beg: opt.Beg.String(),
		Ttl: opt.TTL.Seconds(),
		Tau: opt.Tau.Seconds(),
	}
}
