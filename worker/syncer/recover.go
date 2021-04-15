package syncer

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/schollz/progressbar/v3"
)

const recoverCheckpointKey = "4swap_recover_sync_checkpoint"

func (w *Syncer) recover(ctx context.Context, offset time.Time) error {
	log := logger.FromContext(ctx)

	v, err := w.property.Get(ctx, recoverCheckpointKey)
	if err != nil {
		log.WithError(err).Errorln("property.Get", recoverCheckpointKey)
		return err
	}

	if v.Time().After(offset) {
		offset = v.Time()
	}

	const LIMIT = 500

	for {
		outputs, err := w.walletz.Pull(ctx, offset, LIMIT)
		if err != nil {
			log.WithError(err).Errorln("walletz.Pull")
			return err
		}

		if len(outputs) == 0 {
			break
		}

		offset = outputs[len(outputs)-1].UpdatedAt

		if err := w.wallets.Save(ctx, outputs, true); err != nil {
			log.WithError(err).Errorln("outputs.Save")
			return err
		}

		if err := w.property.Save(ctx, recoverCheckpointKey, offset); err != nil {
			log.WithError(err).Errorln("property.Save")
			return err
		}

		fmt.Printf("\rPulling /multisig/outputs %s %d", offset.UTC().Format(time.RFC3339), len(outputs))

		if len(outputs) < LIMIT {
			break
		}
	}

	log.Infoln("Pull /multisig/outputs done")

	if err := w.recoverOutputs(ctx); err != nil {
		return err
	}

	if err := w.property.Save(ctx, checkpointKey, offset); err != nil {
		log.WithError(err).Errorln("property.Save")
		return err
	}

	log.Infoln("Finish recovery mode")
	return nil
}

func (w *Syncer) recoverOutputs(ctx context.Context) error {
	log := logger.FromContext(ctx)

	count, err := w.wallets.CountRecovery(ctx)
	if err != nil {
		log.WithError(err).Errorln("outputs.Count")
		return err
	}

	if count == 0 {
		return nil
	}

	bar := progressbar.Default(count, "recover outputs")

	offset := time.Time{}
	const LIMIT = 1000

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		outputs, err := w.wallets.ListRecovery(ctx, offset, LIMIT)
		if err != nil {
			log.WithError(err).Errorln("outputs.List")
			return err
		}

		if len(outputs) == 0 {
			break
		}

		_ = bar.Add(len(outputs))
		offset = outputs[len(outputs)-1].CreatedAt

		if err := w.wallets.Save(ctx, outputs, false); err != nil {
			log.WithError(err).Errorln("wallets.Save")
			return err
		}
	}

	_ = bar.Finish()
	return nil
}
