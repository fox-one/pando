package actions

import (
	"crypto/ed25519"
	"encoding/base64"
	"net/http"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/param"
	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/mtg/types"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/twitchtv/twirp"
)

func HandleCreate(walletz core.WalletService, system *core.System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body struct {
			Actions []string        `json:"actions,omitempty"`
			AssetID string          `json:"asset_id,omitempty"`
			Amount  decimal.Decimal `json:"amount,omitempty"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, err)
			return
		}

		data, err := types.EncodeWithTypes(body.Actions...)
		if err == nil {
			key := mixin.GenerateEd25519Key()
			pub := system.PrivateKey.Public().(ed25519.PublicKey)
			data, err = mtg.Encrypt(data, key, pub)
		}

		if err != nil {
			logger.FromContext(ctx).WithError(err).Debugln("EncodeWithTypes", body.Actions)
			render.Error(w, twirp.InvalidArgumentError("actions", "encode failed"))
			return
		}

		memo := base64.StdEncoding.EncodeToString(data)
		view := render.H{
			"memo": memo,
		}

		if body.AssetID != "" && body.Amount.Truncate(8).IsPositive() {
			transfer := &core.Transfer{
				TraceID:   uuid.New(),
				AssetID:   body.AssetID,
				Amount:    body.Amount.Truncate(8),
				Memo:      memo,
				Threshold: system.Threshold,
				Opponents: system.MemberIDs(),
			}

			code, err := walletz.ReqTransfer(ctx, transfer)
			if err != nil {
				render.BadRequest(w, err)
				return
			}

			view["code"] = code
			view["code_url"] = mixin.URL.Codes(code)
		}
	}
}
