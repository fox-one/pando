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
		// FlipCreated called when a new flip created
		FlipCreated(ctx context.Context, flip *Flip) error
		// FlipBid called when a flip bid created
		FlipBid(ctx context.Context, flip *Flip, event *FlipEvent) error
		// FlipDeal called when a flip completed
		FlipDeal(ctx context.Context, flip *Flip) error
	}
)
