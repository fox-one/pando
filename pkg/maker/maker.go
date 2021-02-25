package maker

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

type HandlerFunc func(ctx context.Context, r *Request) error

type Request struct {
	UTXO   *core.Output
	Action core.Action
	Body   []byte
	Gov    bool

	UserID   string
	FollowID string
}

func (r *Request) Scan(dest ...interface{}) error {
	b, err := mtg.Scan(r.Body, dest...)
	if err != nil {
		return err
	}

	r.Body = b
	return nil
}

func (r *Request) WithBody(values ...interface{}) *Request {
	b, err := mtg.Encode(values...)
	if err != nil {
		panic(err)
	}

	r2 := new(Request)
	*r2 = *r
	r2.Body = b

	return r2
}

func (r *Request) BindUser() error {
	var id uuid.UUID
	if err := r.Scan(&id); err != nil {
		return err
	}

	r.UserID = id.String()
	return nil
}

func (r *Request) BindFollow() error {
	var id uuid.UUID
	if err := r.Scan(&id); err != nil {
		return err
	}

	r.FollowID = id.String()
	return nil
}

func (r *Request) Version() int64 {
	return r.UTXO.ID
}

func (r *Request) Now() time.Time {
	return r.UTXO.CreatedAt
}

func (r *Request) TraceID() string {
	return r.UTXO.TraceID
}

func (r *Request) Payment() (string, decimal.Decimal) {
	return r.UTXO.AssetID, r.UTXO.Amount
}

func (r *Request) Tx() *core.Transaction {
	asset, amount := r.Payment()
	return &core.Transaction{
		CreatedAt: r.Now(),
		TraceID:   r.TraceID(),
		AssetID:   asset,
		Amount:    amount,
		Action:    r.Action,
		UserID:    r.UserID,
		FollowID:  r.FollowID,
	}
}
