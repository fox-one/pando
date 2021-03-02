package events

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
)

const (
	checkpoint = "tx_notify_checkpoint"
)

func New(
	transactions core.TransactionStore,
	notifier core.Notifier,
	properties property.Store,
) *Events {
	return &Events{
		transactions: transactions,
		notifier:     notifier,
		properties:   properties,
	}
}

type Events struct {
	transactions core.TransactionStore
	notifier     core.Notifier
	properties   property.Store
}

func (w *Events) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "notify")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 300 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Events) run(ctx context.Context) error {
	v, err := w.properties.Get(ctx, checkpoint)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorf("properties.Get(%q)", checkpoint)
		return err
	}

	from := v.Int64()
	const Limit = 100

	transactions, err := w.transactions.List(ctx, from, Limit)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("transactions.List")
		return err
	}

	if len(transactions) == 0 {
		return errors.New("EOF")
	}

	for _, tx := range transactions {
		if err := w.notifier.Transaction(ctx, tx); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("notifier.Transaction")
			return err
		}

		from = tx.ID
	}

	if err := w.properties.Save(ctx, checkpoint, from); err != nil {
		logger.FromContext(ctx).WithError(err).Errorf("properties.Save(%q,%v)", checkpoint, from)
		return err
	}

	return nil
}
