package system

import (
	"net/http"

	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pkg/property"
)

func HandleProperty(property property.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, err := property.List(r.Context())
		if err != nil {
			render.Error(w, err)
			return
		}

		render.JSON(w, values)
	}
}
