package metric

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRender(t *testing.T) {
	groups := []Group{
		{
			Name: "Properties",
			Entries: []Entry{
				{
					Name:  "foo",
					Value: "bar",
				},
				{
					Name:  "x",
					Value: "y",
				},
			},
		},
		{
			Name: "System",
			Entries: []Entry{
				{
					Name:  "uptime",
					Value: "30m",
				},
			},
		},
	}

	var b bytes.Buffer
	Render(&b, groups)
	fmt.Println(b.String())
}
