package oracle

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

type (
	Signer struct {
		Index     int    `json:"index,omitempty"`
		VerifyKey string `json:"verify_key,omitempty"`
	}

	Receiver struct {
		Members   []string `json:"members,omitempty"`
		Threshold uint8    `json:"threshold"`
	}

	PriceRequest struct {
		AssetID   string   `json:"asset_id"`
		TraceID   string   `json:"trace_id,omitempty"`
		Receiver  Receiver `json:"receiver,omitempty"`
		Signers   []Signer `json:"signers,omitempty"`
		Threshold int64    `json:"threshold"`
	}
)

func HandleScanRequests(
	oracles core.OracleStore,
	system *core.System,
) http.HandlerFunc {
	nonce := uuid.New()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		feeds, err := oracles.ListFeeds(ctx)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("api: cannot list oracle feeds")
			render.Error(w, err)
		}

		signers := make([]Signer, len(feeds))
		for idx, feed := range feeds {
			signers[idx] = Signer{
				Index:     idx + 1,
				VerifyKey: feed.PublicKey,
			}
		}

		list, err := oracles.List(ctx)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("api: cannot list oracles")
			render.Error(w, err)
		}

		resp := make([]PriceRequest, 0, len(list))
		for _, oracle := range list {
			next := oracle.NextPeekAt()
			if time.Until(next) > 0 {
				continue
			}

			resp = append(resp, PriceRequest{
				AssetID: oracle.AssetID,
				TraceID: uuid.Modify(nonce, fmt.Sprintf("%s-%s", oracle.AssetID, next.Format(time.RFC3339Nano))),
				Receiver: Receiver{
					Members:   system.Members,
					Threshold: system.Threshold,
				},
				Signers:   signers,
				Threshold: oracle.Threshold,
			})
		}

		render.JSON(w, resp)
	}
}
