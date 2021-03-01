package pay

import (
	"context"
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/qrcode"
	"github.com/shopspring/decimal"
)

const (
	CNB = "965e5c6e-434c-3fa9-b780-c50f43cd955c"
)

var store *Keystore

type Keystore struct {
	mixin.Keystore
	Pin string `json:"pin,omitempty"`
}

func UseKeystore(s *Keystore) {
	store = s
}

func Request(ctx context.Context, assetID string, amount decimal.Decimal, memo string) error {
	input := mixin.TransferInput{
		AssetID: assetID,
		Amount:  amount,
		TraceID: uuid.New(),
		Memo:    memo,
	}
	input.OpponentMultisig.Receivers = cfg.GetGroupMembers()
	input.OpponentMultisig.Threshold = uint8(cfg.GetGroupThreshold())

	if store != nil {
		client, err := mixin.NewFromKeystore(&store.Keystore)
		if err != nil {
			return err
		}

		tx, err := client.Transaction(ctx, &input, store.Pin)
		if err != nil {
			return err
		}

		fmt.Printf("✅ transfer done with trace %s", tx.TraceID)
		return nil
	}

	token := cfg.GetAuthToken()
	client := mixin.NewFromAccessToken(token)

	payment, err := client.VerifyPayment(ctx, input)
	if err != nil {
		return err
	}

	url := mixin.URL.Codes(payment.CodeID)
	qrcode.Print(url)

	return nil
}
