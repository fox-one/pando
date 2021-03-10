package views

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

func Milliseconds(t *time.Time) float64 {
	if t == nil {
		return 0
	}

	return float64(t.UnixNano() / int64(time.Millisecond))
}

func Unix(t *time.Time) float64 {
	if t == nil {
		return 0
	}

	return float64(t.Unix())
}

func Time(t *time.Time) *tspb.Timestamp {
	if t == nil {
		return nil
	}

	pt, _ := ptypes.TimestampProto(*t)
	return pt
}
