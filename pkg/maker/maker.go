package maker

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/shopspring/decimal"
)

type HandlerFunc func(r *Request) error

type Request struct {
	Now      time.Time
	Version  int64
	TraceID  string
	Sender   string
	FollowID string
	AssetID  string
	Amount   decimal.Decimal
	Action   core.Action
	Body     []byte
	ctx      context.Context
	values   []interface{}

	// 权限
	Gov    bool // 来自节点投票
	Oracle bool // 来自价格预言机

	Next *Request
}

func (r *Request) Scan(dest ...interface{}) error {
	b, err := mtg.Scan(r.Body, dest...)
	if err != nil {
		return err
	}

	r.Body = b
	r.values = append(r.values, dest...)

	return nil
}

func (r *Request) Values() []interface{} {
	if r.values == nil {
		return []interface{}{}
	}

	return r.values[:]
}

func (r *Request) Copy() *Request {
	r2 := new(Request)
	*r2 = *r
	return r2
}

func (r *Request) WithBody(values ...interface{}) *Request {
	b, err := mtg.Encode(values...)
	if err != nil {
		panic(err)
	}

	r2 := r.Copy()
	r2.Body = b
	r2.values = nil

	return r2
}

func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}

	r2 := r.Copy()
	r2.ctx = ctx
	return r2
}

func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}

	return context.Background()
}

func (r *Request) WithProposal(p *core.Proposal) *Request {
	r2 := r.Copy()
	r2.TraceID = uuid.Modify(r.TraceID, p.TraceID)
	r2.Sender = p.Creator
	r2.AssetID = p.AssetID
	r2.FollowID = p.TraceID
	r2.Amount = p.Amount
	r2.Action = p.Action
	r2.Body, _ = base64.StdEncoding.DecodeString(p.Data)
	r2.values = nil
	r2.Gov = true

	return r2
}

func (r *Request) Tx() *core.Transaction {
	return &core.Transaction{
		CreatedAt: r.Now,
		TraceID:   r.TraceID,
		Version:   r.Version,
		UserID:    r.Sender,
		FollowID:  r.FollowID,
		AssetID:   r.AssetID,
		Amount:    r.Amount,
		Action:    r.Action,
		Status:    core.TransactionStatusPending,
	}
}
