package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/fox-one/pando/cmd/pando-worker/config"
	"github.com/fox-one/pando/server"
	"github.com/fox-one/pando/worker"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	_flag struct {
		notify  bool
		keeper  bool
		debug   bool
		port    int
		cfgFile string

		cashier struct {
			batch    int
			capacity int64
		}

		datadog struct {
			interval time.Duration
		}
	}

	// go build -ldflags -X main.version={{.version}}
	version string
	commit  string
	embed   string
)

func init() {
	flag.BoolVar(&_flag.notify, "notify", false, "enable notifier")
	flag.BoolVar(&_flag.keeper, "keeper", false, "run keeper")
	flag.BoolVar(&_flag.debug, "debug", false, "debug mode")
	flag.IntVar(&_flag.port, "port", 7777, "server port")
	flag.StringVar(&_flag.cfgFile, "config", "", "config filename")

	// worker.cashier.Config
	flag.IntVar(&_flag.cashier.batch, "cashier.batch", 100, "custom batch for worker cashier")
	flag.Int64Var(&_flag.cashier.capacity, "cashier.capacity", 1, "custom capacity for worker cashier")

	// worker.datadog.Config
	flag.DurationVar(&_flag.datadog.interval, "datadog.interval", 5*time.Minute, "custom datadog trigger interval")
}

func main() {
	flag.Parse()

	logrus.Infof("pando worker %s(%s) launched at port %d", version, commit, _flag.port)

	if _flag.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg, err := config.Viperion(_flag.cfgFile, embed)
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
