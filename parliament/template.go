package parliament

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func codeBlock(tag string) string {
	return fmt.Sprintf("```%s", tag)
}

func toUpper(s string) string {
	return strings.ToUpper(s)
}

const proposalTpl = `
## #{{.Number}} NEW PROPOSAL "{{.Action}}"

### INFO

{{ codeBlock "yaml" }}
{{ range .Info -}}
{{.Key}}: {{.Value}}
{{ end }}
{{- codeBlock ""}}

### {{ .Action | upper }}

{{ codeBlock "yaml" }}
{{ range .Meta -}}
{{.Key}}: {{.Value}}
{{ end }}
{{- codeBlock ""}}
`

func renderProposal(p Proposal) []byte {
	t, err := template.New("-").Funcs(template.FuncMap{
		"codeBlock": codeBlock,
		"upper":     toUpper,
	}).Parse(proposalTpl)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, p); err != nil {
		panic(err)
	}

	return bytes.TrimSpace(b.Bytes())
}

const approvedByTpl = `
✅ Approved By {{.ApprovedBy}}

({{.ApprovedCount}} Votes In Total)
`

func renderApprovedBy(p Proposal) []byte {
	t, err := template.New("-").Parse(approvedByTpl)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, p); err != nil {
		panic(err)
	}

	return bytes.TrimSpace(b.Bytes())

}

const passedTpl = "🎉 Proposal Passed"
