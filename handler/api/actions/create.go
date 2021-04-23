package actions

import (
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

type CreateRequest struct {
	// follow id to track tx (uuid)
	FollowID string `json:"follow_id,omitempty" format:"uuid" valid:"uuid,required"`
	// tx parameters
	Parameters []string `json:"parameters,omitempty" valid:"required"`
	// payment asset id (optional)
	AssetID string `json:"asset_id,omitempty" format:"uuid"`
	// payment amount (optional)
	Amount decimal.Decimal `json:"amount,omitempty" swaggertype:"number"`
}

type CreateResponse struct {
	// payment memo
	Memo string `json:"memo,omitempty"`
	// multisig payment code
	Code string `json:"code,omitempty"`
	// multisig payment code url
	CodeURL string `json:"code_url,omitempty"`
}

// HandleCreate godoc
// @Summary request payment code
// @Description
// @Tags actions
// @Accept  json
// @Produce  json
// @Param request body CreateRequest false "request payments"
// @Success 200 {object} CreateResponse
// @Router /actions [post]
func HandleCreate(walletz core.WalletService, system *core.System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body CreateRequest
		if err := param.Binding(r, &body); err != nil {
			render.Error(w, err)
			return
		}

		follow, _ := uuid.FromString(body.FollowID)
		data, err := types.EncodeWithTypes(body.Parameters...)
		if err == nil {
			data, _ = core.TransactionAction{
				FollowID: follow.Bytes(),
				Body:     data,
			}.Encode()

			key := mixin.GenerateEd25519Key()
			data, err = mtg.Encrypt(data, key, system.PublicKey)
		}

		if err != nil {
			logger.FromContext(ctx).WithError(err).Debugln("EncodeWithTypes", body.Parameters)
			render.Error(w, twirp.InvalidArgumentError("actions", "encode failed"))
			return
		}

		memo := base64.StdEncoding.EncodeToString(data)
		resp := CreateResponse{Memo: memo}

		if body.AssetID != "" && body.Amount.Truncate(8).IsPositive() {
			transfer := &core.Transfer{
				TraceID:   uuid.New(),
				AssetID:   body.AssetID,
				Amount:    body.Amount.Truncate(8),
				Memo:      memo,
				Threshold: system.Threshold,
				Opponents: system.Members,
			}

			code, err := walletz.ReqTransfer(ctx, transfer)
			if err != nil {
				render.BadRequest(w, err)
				return
			}

			resp.Code = code
			resp.CodeURL = mixin.URL.Codes(code)
		}

		render.JSON(w, resp)
	}
}
