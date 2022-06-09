package views

import (
	"github.com/fox-one/pando/core"
	api "github.com/fox-one/pando/handler/rpc/pando"
)

func Proposal(p *core.Proposal, items ...core.ProposalItem) *api.Proposal {
	v := &api.Proposal{
		Id:        p.TraceID,
		CreatedAt: Time(&p.CreatedAt),
		PassedAt:  Time(&p.PassedAt.Time),
		Creator:   p.Creator,
		AssetId:   p.AssetID,
		Amount:    p.Amount.String(),
		Action:    api.Action(p.Action),
		Data:      p.Data,
		Votes:     p.Votes,
	}

	if len(items) > 0 {
		v.Items = make([]*api.Proposal_Item, len(items))
	}

	for idx, item := range items {
		v.Items[idx] = &api.Proposal_Item{
			Key:    item.Key,
			Value:  item.Value,
			Hint:   item.Hint,
			Action: item.Action,
		}
	}

	return v
}
