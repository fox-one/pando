package parliament

import (
	"fmt"
	"testing"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/uuid"
)

func TestRenderProposal(t *testing.T) {
	p := Proposal{
		Number: 1,
		Action: core.ProposalActionWithdraw.String(),
		Info: []Item{
			{
				Key:   "id",
				Value: uuid.New(),
			},
			{
				Key:    "user",
				Value:  uuid.New(),
				Action: fmt.Sprintf("mixin://users/%s", uuid.New()),
			},
			{
				Key:    "creator",
				Value:  uuid.New(),
				Action: fmt.Sprintf("mixin://users/%s", uuid.New()),
			},
		},
	}

	view := renderProposal(p)
	fmt.Println(string(view))
}
