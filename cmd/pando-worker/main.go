package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/fox-one/pando/cmd/pando-worker/config"
	"github.com/fox-one/pando/server"
	"github.com/fox-one/pando/worker"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	_notify  = flag.Bool("notify", false, "enable notifier")
	_keeper  = flag.Bool("keeper", false, "run keeper")
	_debug   = flag.Bool("debug", false, "debug mode")
	_port    = flag.Int("port", 7777, "server port")
	_cfgFile = flag.String("config", "", "config filename")

	// worker.cashier.Config
	_cashierBatch    = flag.Int("cashier.batch", 100, "custom batch for worker cashier")
	_cashierCapacity = flag.Int64("cashier.capacity", 1, "custom capacity for worker cashier")
)

// build embed
var (
	version string
	commit  string
	embed   string
)

func main() {
	flag.Parse()

	logrus.Infof("pando worker %s(%s) launched at port %d", version, commit, *_port)

	if *_debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg, err := config.Viperion(*_cfgFile, embed)
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	app, err := buildApp(cfg)
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: cannot initialize worker")
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	for i := range app.workers {
		w := app.workers[i]
		g.Go(func() error {
			return w.Run(ctx)
		})
	}

	g.Go(func() error {
		return app.server.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("program terminated")
	}
}

type app struct {
	workers []worker.Worker
	server  *server.Server
}
