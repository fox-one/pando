package parliament

import (
	"context"

	"github.com/fox-one/pando/core"
)

func (s *parliament) renderProposalItems(ctx context.Context, p *core.Proposal) []Item {
	items, _ := s.proposalz.ListItems(ctx, p)

	results := make([]Item, len(items))
	for idx, item := range items {
		results[idx] = Item{
			Key:    item.Key,
			Value:  item.Value,
			Action: item.Action,
		}

		if item.Hint != "" {
			results[idx].Value = item.Hint
		}
	}

	return results
}
