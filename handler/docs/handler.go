package docs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/alecthomas/template"
	"github.com/fox-one/pando/handler/render"
)

func Handler(version string) http.Handler {
	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)

	if err != nil {
		panic(err)
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		info := SwaggerInfo
		info.Version = version
		info.Host = r.Host
		info.Description = strings.Replace(info.Description, "\n", "\\n", -1)

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, info); err != nil {
			render.Error(w, err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		_, _ = io.Copy(w, &tpl)
	}

	return http.HandlerFunc(f)
}
