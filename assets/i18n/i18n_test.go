package i18n

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

const tx_ok = `
✅ {{.Action}}成功

{{ range .Lines -}}
- {{.}}
{{end}}
`

func TestTxOK(t *testing.T) {
	tmp, err := template.New("-").Parse(tx_ok)
	require.Nil(t, err)

	b := bytes.Buffer{}
	_ = tmp.Execute(&b, map[string]interface{}{
		"Action": "action",
		"Lines": []string{
			"aaaaaa",
			"bbbbbb",
		},
	})

	t.Log(b.String())
}
