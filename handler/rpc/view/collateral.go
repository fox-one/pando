package view

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
)

func Collateral(cat *core.Collateral) *api.Collateral {
	return &api.Collateral{
		Id:        cat.TraceID,
		CreatedAt: Time(&cat.CreatedAt),
		Name:      cat.Name,
		Gem:       cat.Gem,
		Dai:       cat.Dai,
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
		Live:      cat.Live > 0,
	}
}
