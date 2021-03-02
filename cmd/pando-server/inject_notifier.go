package main

import (
	"github.com/fox-one/pando/notifier"
	"github.com/google/wire"
)

var notifierSet = wire.NewSet(
	notifier.New,
)
