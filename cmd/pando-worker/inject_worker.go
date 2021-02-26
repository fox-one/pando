package main

import (
	"github.com/fox-one/pando/worker"
	"github.com/fox-one/pando/worker/cashier"
	"github.com/fox-one/pando/worker/messenger"
	"github.com/fox-one/pando/worker/payee"
	"github.com/fox-one/pando/worker/pricesync"
	"github.com/fox-one/pando/worker/spentsync"
	"github.com/fox-one/pando/worker/syncer"
	"github.com/fox-one/pando/worker/txsender"
	"github.com/google/wire"
)

var workerSet = wire.NewSet(
	cashier.New,
	messenger.New,
	payee.New,
	pricesync.New,
	spentsync.New,
	syncer.New,
	txsender.New,
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
) []worker.Worker {
	return []worker.Worker{a, b, c, d, e, f, g}
}
