package core

import "context"

type (
	// Parliament is a proposal version notifier to mtg member admins
	Parliament interface {
		// ProposalCreated Called when a new proposal created
		ProposalCreated(ctx context.Context, proposal *Proposal) error
		// ProposalApproved called when a proposal has a new vote
		ProposalApproved(ctx context.Context, proposal *Proposal) error
		// ProposalPassed called when a proposal is passed
		ProposalPassed(ctx context.Context, proposal *Proposal) error
	}
)
