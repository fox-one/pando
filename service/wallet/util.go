package wallet

import (
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
)

func convertToOutput(utxo *mixin.MultisigUTXO) *core.Output {
	return &core.Output{
		CreatedAt:       utxo.CreatedAt,
		UpdatedAt:       utxo.UpdatedAt,
		TraceID:         utxo.UTXOID,
		AssetID:         utxo.AssetID,
		Amount:          utxo.Amount,
		Memo:            utxo.Memo,
		State:           utxo.State,
		TransactionHash: utxo.TransactionHash.String(),
		OutputIndex:     utxo.OutputIndex,
		SignedTx:        utxo.SignedTx,
	}
}

func convertToUTXO(output *core.Output, members []string, threshold uint8) *mixin.MultisigUTXO {
	hash, err := mixin.HashFromString(output.TransactionHash)
	if err != nil {
		panic(fmt.Errorf("hash from %q failed: %w", output.TransactionHash, err))
	}

	return &mixin.MultisigUTXO{
		UTXOID:          output.TraceID,
		AssetID:         output.AssetID,
		TransactionHash: hash,
		OutputIndex:     output.OutputIndex,
		Amount:          output.Amount,
		Threshold:       threshold,
		Members:         members,
		Memo:            output.Memo,
		State:           output.State,
		SignedTx:        output.SignedTx,
	}
}
