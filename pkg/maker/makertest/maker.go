package makertest

import (
	"time"

	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pando/pkg/uuid"
)

var req = &maker.Request{
	Now:      time.Now(),
	Version:  1,
	TraceID:  uuid.New(),
	Sender:   uuid.New(),
	FollowID: uuid.New(),
	AssetID:  uuid.New(),
	Amount:   number.One,
}

func step() {
	req.Now = req.Now.Add(time.Second)
	req.Version += 1
	req.TraceID = uuid.New()
}

func Next() *maker.Request {
	step()
	return req.Copy()
}
