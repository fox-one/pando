package main

import (
	"github.com/fox-one/pando/notifier"
	"github.com/google/wire"
)

var notifierSet = wire.NewSet(
	wire.Value(notifier.Config{}),
	notifier.New,
)
