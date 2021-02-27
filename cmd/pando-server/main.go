package main

import (
	"context"
	"flag"

	"github.com/drone/signal"
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/server"
	"github.com/sirupsen/logrus"
)

var (
	debug = flag.Bool("debug", false, "debug mode")
	port  = flag.Int("port", 7778, "server port")
)

func main() {
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg, err := config.Viperion()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	svr, err := buildServer(cfg)
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: cannot initialize worker")
	}

	logrus.Infof("pando mtg api server with version %q launched at port %d!", "v0.0.1", *port)

	ctx := signal.WithContext(context.Background())
	if err := svr.ListenAndServe(ctx); err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("program terminated")
	}
}

type app struct {
	server *server.Server
}
