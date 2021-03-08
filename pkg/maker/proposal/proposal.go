package proposal

import (
	"encoding/base64"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/mtg"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return maker.Require(condition, "proposal/"+msg)
}

func From(r *maker.Request, proposals core.ProposalStore) (*core.Proposal, error) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	p, err := proposals.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("proposals.Find")
		return nil, err
	}

	if err := require(p.ID > 0, "not init"); err != nil {
		return nil, err
	}

	return p, nil
}

func handleProposal(r *maker.Request, walletz core.WalletService, system *core.System, action core.Action, p *core.Proposal) error {
	uid, _ := uuid.FromString(system.ClientID)
	pid, _ := uuid.FromString(p.TraceID)
	data, _ := mtg.Encode(action, pid)
	data, _ = core.TransactionAction{
		UserID: uid.Bytes(),
		Body:   data,
	}.Encode()
	data, _ = mtg.Encrypt(data, mixin.GenerateEd25519Key(), system.PublicKey)
	memo := base64.StdEncoding.EncodeToString(data)

	ctx := r.Context()
	if err := walletz.HandleTransfer(ctx, &core.Transfer{
		TraceID:   uuid.Modify(r.TraceID, p.TraceID+system.ClientID),
		AssetID:   system.VoteAsset,
		Amount:    system.VoteAmount,
		Threshold: system.Threshold,
		Opponents: system.Members,
		Memo:      memo,
	}); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("wallets.HandleTransfer")
		return err
	}

	return nil
}
