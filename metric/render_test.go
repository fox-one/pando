package metric

import (
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

	out := Render(groups)
	fmt.Println(out)
}
