package views

import (
	"github.com/fox-one/pando/core"
	api "github.com/fox-one/pando/handler/rpc/pando"
)

func Stat(stat *core.Stat) *api.Stat {
	return &api.Stat{
		CollateralId: stat.CollateralID,
		Date:         Time(&stat.Date),
		Timestamp:    stat.Date.Unix(),
		Gem:          stat.Gem,
		Dai:          stat.Dai,
		Ink:          stat.Ink.String(),
		Debt:         stat.Debt.String(),
		GemPrice:     stat.GemPrice.String(),
		DaiPrice:     stat.DaiPrice.String(),
	}
}

func AggregatedStat(stat core.AggregatedStat) *api.AggregatedStat {
	return &api.AggregatedStat{
		Date:      Time(&stat.Date),
		Timestamp: stat.Date.Unix(),
		GemValue:  stat.GemValue.String(),
		DaiValue:  stat.DaiValue.String(),
	}
}
