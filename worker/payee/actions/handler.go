package actions

import (
	"context"
	"errors"

	"github.com/fox-one/pando/pkg/maker"
)

var (
	ErrBuildAbort = errors.New("actions: build-abort")
)

type Handler interface {
	Build(ctx context.Context, tx *maker.Tx, body []byte) ([]byte, error)
	Apply(ctx context.Context, tx *maker.Tx, body []byte) error
}
