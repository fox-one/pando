package parliament

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/fox-one/pkg/text/columnize"
)

type Item struct {
	Key    string
	Value  string
	Action string
}

type Proposal struct {
	Number int64
	Action string
	Info   []Item
	Meta   []Item

	ApprovedCount int
	ApprovedBy    string
}

func AlignKey(items []Item) {
	var form columnize.Form

	for _, item := range items {
		form.Append(item.Key, item.Value)
	}

	var b bytes.Buffer
	_ = form.Fprint(&b)

	s := bufio.NewScanner(&b)
	for idx, item := range items {
		if !s.Scan() {
			break
		}

		line := s.Text()
		if fields := strings.Fields(line); len(fields) >= 2 {
			pos := strings.Index(line, fields[1])
			item.Key = line[:pos]
			items[idx] = item
		}
	}
}
