package metric

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed render.tmpl
var tmpl string

var t = template.Must(
	template.New("-").Parse(tmpl),
)

func Render(w io.Writer, groups []Group) {
	if err := t.Execute(w, groups); err != nil {
		panic(err)
	}
}
