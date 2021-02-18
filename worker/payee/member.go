package payee

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/core/proposal"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
	"github.com/gofrs/uuid"
)

func (w *Payee) handleMemberAction(ctx context.Context, output *core.Output, member *core.Member, body []byte) error {
	log := logger.FromContext(ctx)

	var (
		trace  uuid.UUID
		action int
	)

	body, err := mtg.Scan(body, &trace, &action)
	if err != nil {
		log.WithError(err).Debugln("scan proposal trace & action failed")
		return nil
	}

	log.WithField("trace", trace.String()).Debugf("handle member action %s", core.ProposalAction(action).String())

	if core.ProposalAction(action) == core.ProposalActionVote {
		return w.voteProposal(ctx, output, member, trace.String())
	}

	// new proposal
	p := &core.Proposal{
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.CreatedAt,
		TraceID:   trace.String(),
		Creator:   member.ClientID,
		AssetID:   output.AssetID,
		Amount:    output.Amount,
		Action:    core.ProposalAction(action),
	}

	var content interface{}

	switch p.Action {
	case core.ProposalActionAddPair:
		content = &proposal.AddPair{}
	case core.ProposalActionWithdraw:
		content = &proposal.Withdraw{}
	case core.ProposalActionSwapMethod:
		content = &proposal.SwapMethod{}
	case core.ProposalActionSetProperty:
		content = &proposal.SetProperty{}
	default:
		log.Panicf("unknown proposal action %s", p.Action.String())
	}

	if _, err := mtg.Scan(body, content); err != nil {
		log.WithError(err).Debugln("decode proposal content failed")
	}

	p.Content, _ = json.Marshal(content)

	if err := w.proposals.Create(ctx, p); err != nil {
		log.WithError(err).Errorln("proposals.Create")
		return err
	}

	if err := w.parliament.ProposalCreated(ctx, p, member); err != nil {
		log.WithError(err).Errorln("notifier.ProposalVoted")
		return err
	}

	return nil
}

func (w *Payee) voteProposal(ctx context.Context, output *core.Output, member *core.Member, traceID string) error {
	log := logger.FromContext(ctx).WithField("proposal", traceID)

	p, err := w.proposals.Find(ctx, traceID)
	if err != nil {
		// 如果 proposal 不存在，直接跳过
		if store.IsErrNotFound(err) {
			log.WithError(err).Debugln("proposal not found")
			return nil
		}

		log.WithError(err).Errorln("proposals.Find")
		return err
	}

	passed := p.PassedAt.Valid

	if !passed && !govalidator.IsIn(member.ClientID, p.Votes...) {
		p.Votes = append(p.Votes, member.ClientID)
		log.Infof("Proposal Voted by %s", member.ClientID)

		if err := w.parliament.ProposalApproved(ctx, p, member); err != nil {
			log.WithError(err).Errorln("notifier.ProposalVoted")
			return err
		}

		if passed = len(p.Votes) >= int(w.system.Threshold); passed {
			p.PassedAt = sql.NullTime{
				Time:  output.CreatedAt,
				Valid: true,
			}

			log.Infof("Proposal Approved")
			if err := w.parliament.ProposalPassed(ctx, p); err != nil {
				log.WithError(err).Errorln("notifier.ProposalApproved")
				return err
			}
		}

		if err := w.proposals.Update(ctx, p); err != nil {
			log.WithError(err).Errorln("proposals.Update")
			return err
		}
	}

	if passed {
		return w.handlePassedProposal(ctx, p)
	}

	return nil
}
