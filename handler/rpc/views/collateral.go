package views

import (
	"github.com/fox-one/pando/core"
	api "github.com/fox-one/pando/handler/rpc/pando"
)

func Collateral(cat *core.Collateral) *api.Collateral {
	return &api.Collateral{
		Id:        cat.TraceID,
		CreatedAt: Time(&cat.CreatedAt),
		Name:      cat.Name,
		Gem:       cat.Gem,
		Dai:       cat.Dai,
		Ink:       cat.Ink.String(),
		Art:       cat.Art.String(),
		Rate:      cat.Rate.String(),
		Rho:       Time(&cat.Rho),
		Debt:      cat.Debt.String(),
		Line:      cat.Line.String(),
		Dust:      cat.Dust.String(),
		Price:     cat.Price.String(),
		Mat:       cat.Mat.String(),
		Duty:      cat.Duty.String(),
		Chop:      cat.Chop.String(),
		Dunk:      cat.Dunk.String(),
		Beg:       cat.Beg.String(),
		Ttl:       int32(cat.TTL),
		Tau:       int32(cat.Tau),
		Box:       cat.Box.String(),
		Litter:    cat.Litter.String(),
		Live:      cat.Live > 0,
		Supply:    cat.Supply.String(),
	}
}
