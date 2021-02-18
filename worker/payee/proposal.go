package payee

import (
	"context"
	"encoding/json"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/core/proposal"
	"github.com/fox-one/pando/pkg/number"
	"github.com/fox-one/pkg/logger"
)

func (w *Payee) handlePassedProposal(ctx context.Context, p *core.Proposal) error {
	switch p.Action {
	case core.ProposalActionWithdraw:
		var action proposal.Withdraw
		_ = json.Unmarshal(p.Content, &action)
		return w.withdraw(ctx, p, action)
	case core.ProposalActionSetProperty:
		var action proposal.SetProperty
		_ = json.Unmarshal(p.Content, &action)
		return w.setProperty(ctx, p, action)
	}

	return nil
}

func (w *Payee) withdraw(ctx context.Context, p *core.Proposal, action proposal.Withdraw) error {
	log := logger.FromContext(ctx)

	if _, err := w.findAsset(ctx, action.Asset); err != nil {
		log.WithError(err).Errorln("find asset by id")
		if err == core.ErrAssetNotExist {
			return nil
		}

		return err
	}

	amount := number.Decimal(action.Amount).Truncate(8)
	if !amount.IsPositive() {
		return nil
	}

	transfer := &core.Transfer{
		TraceID:   p.TraceID,
		AssetID:   action.Asset,
		Amount:    amount,
		Threshold: 1,
		Opponents: []string{action.Opponent},
	}

	if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{transfer}); err != nil {
		log.WithError(err).Errorln("wallets.CreateTransfers")
		return err
	}

	return nil
}

func (w *Payee) swapMethod(ctx context.Context, _ *core.Proposal, action proposal.SwapMethod) error {
	log := logger.FromContext(ctx)

	pair, err := w.pairs.Find(ctx, action.BaseAsset, action.QuoteAsset)
	if err != nil {
		log.WithError(err).Errorln("pairs.Find")
		return err
	}

	// pair not exist
	if pair.ID == 0 {
		return nil
	}

	if pair.SwapMethod == action.Method {
		return nil
	}

	if err := swap.ValidateMethod(action.Method); err != nil {
		return nil
	}

	pair.SwapMethod = action.Method
	if err := w.pairs.Update(ctx, pair); err != nil {
		log.WithError(err).Errorln("pairs.Update")
		return err
	}

	return nil
}

func (w *Payee) setProperty(ctx context.Context, _ *core.Proposal, action proposal.SetProperty) error {
	log := logger.FromContext(ctx)

	if action.Key == "" {
		return nil
	}

	if err := w.property.Save(ctx, action.Key, action.Value); err != nil {
		log.WithError(err).Errorln("update property", action.Key, action.Value)
		return err
	}

	return nil
}
