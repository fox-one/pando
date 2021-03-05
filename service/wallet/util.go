package wallet

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
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

func extractSender(key ed25519.PrivateKey, memo string) (string, string) {
	b, err := base64.StdEncoding.DecodeString(memo)
	if err != nil {
		b, err = base64.RawURLEncoding.DecodeString(memo)
	}

	if err != nil {
		return memo, ""
	}

	data, err := mtg.Decrypt(b, key)
	if err != nil {
		return memo, ""
	}

	var userID, followID uuid.UUID
	data, err = mtg.Scan(data, &userID, &followID)
	if err != nil {
		return memo, ""
	}

	newData, _ := mtg.Encode(followID)
	newData = append(newData, data...)

	newKey := mixin.GenerateEd25519Key()
	encryptedData, err := mtg.Encrypt(newData, newKey, key.Public().(ed25519.PublicKey))
	if err != nil {
		return memo, ""
	}

	return base64.StdEncoding.EncodeToString(encryptedData), userID.String()
}
