package main

import (
	"context"
	"flag"
	"os"

	"github.com/drone/signal"
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/server"
	"github.com/sirupsen/logrus"
)

var (
	debug = flag.Bool("debug", false, "debug mode")
	port  = flag.Int("port", 7778, "server port")

	version, commit string
)

func main() {
	flag.Parse()

	version = os.Getenv("PANDO_VERSION")
	commit = os.Getenv("PANDO_VERSION")

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Infof("pando server %s(%s) launched at port %d", version, commit, *port)

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

	ctx := signal.WithContext(context.Background())
	if err := svr.ListenAndServe(ctx); err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("program terminated")
	}
}

type app struct {
	server *server.Server
}
