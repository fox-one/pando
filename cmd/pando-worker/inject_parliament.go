package main

import (
	"github.com/fox-one/pando/parliament"
	"github.com/google/wire"
)

var parliamentSet = wire.NewSet(
	parliament.New,
)
