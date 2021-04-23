package main

import (
	"github.com/fox-one/pando/worker"
	"github.com/fox-one/pando/worker/assigner"
	"github.com/fox-one/pando/worker/cashier"
	"github.com/fox-one/pando/worker/events"
	"github.com/fox-one/pando/worker/keeper"
	"github.com/fox-one/pando/worker/messenger"
	"github.com/fox-one/pando/worker/payee"
	"github.com/fox-one/pando/worker/pricesync"
	"github.com/fox-one/pando/worker/spentsync"
	"github.com/fox-one/pando/worker/syncer"
	"github.com/fox-one/pando/worker/txsender"
	"github.com/google/wire"
)

var workerSet = wire.NewSet(
	wire.Value(cashier.Config{
		Batch:    *_cashierBatch,
		Capacity: *_cashierCapacity,
	}),
	cashier.New,
	messenger.New,
	payee.New,
	pricesync.New,
	spentsync.New,
	syncer.New,
	txsender.New,
	events.New,
	keeper.New,
	assigner.New,
	provideWorkers,
)

func provideWorkers(
	a *cashier.Cashier,
	b *messenger.Messenger,
	c *payee.Payee,
	d *pricesync.Sync,
	e *spentsync.SpentSync,
	f *txsender.Sender,
	g *syncer.Syncer,
	h *events.Events,
	i *keeper.Keeper,
	j *assigner.Assigner,
) []worker.Worker {
	workers := []worker.Worker{a, b, c, d, e, f, g, h, j}

	if *_keeper {
		workers = append(workers, i)
	}

	return workers
}
