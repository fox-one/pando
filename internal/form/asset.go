package form

import (
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/text/columnize"
)

func Assets(assets []*core.Asset) *columnize.Form {
	f := &columnize.Form{}

	for _, asset := range assets {
		f.Append(asset.Symbol, asset.Name, asset.ID)
	}

	return f
}
