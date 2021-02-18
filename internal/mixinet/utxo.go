package mixinet

import (
	"bytes"
	"sort"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
)

func cmpUTXO(a, b *mixin.MultisigUTXO) int {
	if dur := a.CreatedAt.Sub(b.CreatedAt); dur > 0 {
		return 1
	} else if dur < 0 {
		return -1
	}

	if r := bytes.Compare(a.TransactionHash[:], b.TransactionHash[:]); r != 0 {
		return r
	}

	if i, j := a.OutputIndex, b.OutputIndex; i == j {
		return 0
	} else if i > j {
		return 1
	}

	return -1
}

func SortOutputs(outputs []*core.Output) {
	sort.Slice(outputs, func(i, j int) bool {
		ui, uj := outputs[i].UTXO, outputs[j].UTXO
		return cmpUTXO(ui, uj) < 0
	})
}
