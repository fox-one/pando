package views

import (
	"time"

	"github.com/fox-one/pando/core"
	api "github.com/fox-one/pando/handler/rpc/pando"
)

func Flip(flip *core.Flip, tags ...api.Flip_Tag) *api.Flip {
	tic := time.Unix(flip.Tic, 0)
	end := time.Unix(flip.End, 0)

	return &api.Flip{
		Id:           flip.TraceID,
		CreatedAt:    Time(&flip.CreatedAt),
		Tic:          Time(&tic),
		End:          Time(&end),
		Bid:          flip.Bid.String(),
		Lot:          flip.Lot.String(),
		Tab:          flip.Tab.String(),
		Art:          flip.Art.String(),
		CollateralId: flip.CollateralID,
		VaultId:      flip.VaultID,
		// Guy:          flip.Guy,
		Action: api.Action(flip.Action),
		Tags:   tags,
	}
}

func FlipEvent(event *core.FlipEvent, me string) *api.Flip_Event {
	return &api.Flip_Event{
		FlipId:    event.FlipID,
		CreatedAt: Time(&event.CreatedAt),
		Action:    api.Action(event.Action),
		Bid:       event.Bid.String(),
		Lot:       event.Lot.String(),
		IsMe:      event.Guy != "" && event.Guy == me,
	}
}
