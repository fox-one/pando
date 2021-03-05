package syncer

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
)

const checkpointKey = "sync_checkpoint"

func New(
	wallets core.WalletStore,
	walletz core.WalletService,
	property property.Store,
) *Syncer {
	return &Syncer{
		wallets:  wallets,
		walletz:  walletz,
		property: property,
	}
}

type Syncer struct {
	wallets  core.WalletStore
	walletz  core.WalletService
	property property.Store
}

func (w *Syncer) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "syncer")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (w *Syncer) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	v, err := w.property.Get(ctx, checkpointKey)
	if err != nil {
		log.WithError(err).Errorln("property.Get", checkpointKey)
		return err
	}

	offset := v.Time()

	var (
		outputs   = make([]*core.Output, 0, 8)
		positions = make(map[string]int)
		pos       = 0
	)

	const Limit = 500

	for {
		batch, err := w.walletz.Pull(ctx, offset, Limit)
		if err != nil {
			log.WithError(err).Errorln("walletz.Pull")
			return err
		}

		for _, u := range batch {
			offset = u.UpdatedAt

			p, ok := positions[u.TraceID]
			if ok {
				outputs[p] = u
				continue
			}

			outputs = append(outputs, u)
			positions[u.TraceID] = pos
			pos += 1
		}

		if len(batch) < Limit {
			break
		}
	}

	if len(outputs) == 0 {
		return errors.New("EOF")
	}

	core.SortOutputs(outputs)
	if err := w.wallets.Save(ctx, outputs); err != nil {
		log.WithError(err).Errorln("wallets.Save")
		return err
	}

	if err := w.property.Save(ctx, checkpointKey, offset); err != nil {
		log.WithError(err).Errorln("property.Save", checkpointKey)
		return err
	}

	return nil
}
