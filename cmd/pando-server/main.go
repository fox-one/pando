package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/server"
	"github.com/sirupsen/logrus"
)

var (
	debug   = flag.Bool("debug", false, "debug mode")
	port    = flag.Int("port", 7778, "server port")
	cfgFile = flag.String("config", "", "config filename")

	version, commit, embed string
)

func main() {
	flag.Parse()

	logrus.Infof("pando server %s(%s) launched at port %d", version, commit, *port)

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg, err := config.Viperion(*cfgFile, embed)
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	svr, err := buildServer(cfg)
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

	if err := svr.ListenAndServe(ctx); err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("program terminated")
	}
}

type app struct {
	server *server.Server
}
