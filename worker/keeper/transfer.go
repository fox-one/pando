package keeper

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pkg/logger"
)

func (w *Keeper) handleTransfer(ctx context.Context, trace string, values ...interface{}) error {
	if err := w.walletz.HandleTransfer(ctx, &core.Transfer{
		TraceID:   trace,
		AssetID:   w.system.GasAssetID,
		Amount:    w.system.GasAmount,
		Threshold: w.system.Threshold,
		Opponents: w.system.Members,
		Memo:      buildMemo(w.system.PublicKey, values...),
	}); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("walletz.HandleTransfer")
		return err
	}

	return nil
}

func buildMemo(pk ed25519.PublicKey, values ...interface{}) string {
	body, err := mtg.Encode(values...)
	if err != nil {
		panic(err)
	}

	action := core.TransactionAction{
		Body: body,
	}

	data, err := action.Encode()
	if err != nil {
		panic(err)
	}

	key := mixin.GenerateEd25519Key()
	encryptedData, err := mtg.Encrypt(data, key, pk)
	if err != nil {
		panic(err)
	}

	memo := base64.StdEncoding.EncodeToString(encryptedData)
	if len(memo) > 200 {
		memo = base64.StdEncoding.EncodeToString(data)
	}

	return memo
}
