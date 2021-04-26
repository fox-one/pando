package parliament

import (
	"bytes"
	"embed"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed files/*.tmpl
var files embed.FS

var T *template.Template

func init() {
	T = template.New("_").Funcs(template.FuncMap{
		"upper": strings.ToUpper,
	})

	dir := "files"
	entries, err := files.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		name := entry.Name()
		b, err := files.ReadFile(path.Join(dir, name))
		if err != nil {
			panic(err)
		}

		ext := filepath.Ext(name)
		name = name[0 : len(name)-len(ext)]

		T = template.Must(
			T.New(name).Parse(string(b)),
		)
	}
}

func execute(name string, data interface{}) []byte {
	var b bytes.Buffer

	if err := T.ExecuteTemplate(&b, name, data); err != nil {
		panic(err)
	}

	return bytes.TrimSpace(b.Bytes())
}
