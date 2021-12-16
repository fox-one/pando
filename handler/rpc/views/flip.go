package views

import (
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/rpc/api"
)

func Flip(flip *core.Flip) *api.Flip {
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
		Guy:          flip.Guy,
		Action:       api.Action(flip.Action),
	}
}

func FlipEvent(event *core.FlipEvent, me string) *api.Flip_Event {
	return &api.Flip_Event{
		FlipId:    event.FlipID,
		CreatedAt: Time(&event.CreatedAt),
		Action:    api.Action(event.Action),
		Bid:       event.Bid.String(),
		Lot:       event.Lot.String(),
		IsMe:      event.Guy == me,
	}
}
