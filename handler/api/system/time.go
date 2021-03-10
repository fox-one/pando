package system

import (
	"net/http"
	"time"

	"github.com/fox-one/pando/handler/render"
)

type TimeResponse struct {
	ISO   string `json:"iso,omitempty"`
	Epoch int64  `json:"epoch,omitempty"`
}

// ShowSystemTime godoc
// @Summary Show server time
// @Description
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} TimeResponse
// @Router /time [get]
func HandleTime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		render.JSON(w, TimeResponse{
			ISO:   t.Format(time.RFC3339),
			Epoch: t.Unix(),
		})
	}
}
