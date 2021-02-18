package vat

import (
	"github.com/fox-one/pando/pkg/maker"
)

const (
	module = "vat"
)

var (
	ErrVatCeilingExceeded = maker.RegisterErr(1001, module, "ceiling-exceeded")
	ErrVatNotSafe         = maker.RegisterErr(1002, module, "not-safe")
	ErrVatDust            = maker.RegisterErr(1003, module, "dust")
	ErrVatInkNotMatch     = maker.RegisterErr(1004, module, "ink-not-match")
	ErrVatArtNotMatch     = maker.RegisterErr(1005, module, "art-not-match")
	ErrVatNotAllowed      = maker.RegisterErr(1006, module, "not-allowed")
	ErrVatNotLive         = maker.RegisterErr(1007, module, "not-live")
	ErrVatValidateFailed  = maker.RegisterErr(1008, module, "validate-failed")
)
