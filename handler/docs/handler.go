package docs

import (
	"embed"
	"net/http"
	"strings"

	"github.com/fox-one/pando/handler/render"
)

//go:embed swagger.*
var contents embed.FS

func Handler() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var (
			b   []byte
			err error
		)

		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "yaml") {
			w.Header().Set("Content-type", "application/x-yaml")
			b, err = contents.ReadFile("swagger.yaml")
		} else {
			w.Header().Set("Content-type", "application/json")
			b, err = contents.ReadFile("swagger.json")
		}

		if err != nil {
			render.Error(w, err)
			return
		}

		_, _ = w.Write(b)
	}

	return http.HandlerFunc(f)
}
