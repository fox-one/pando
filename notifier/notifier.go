package notifier

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pando/service/asset"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/text/localizer"
)

func New(
	system *core.System,
	assetz core.AssetService,
	messages core.MessageStore,
	vats core.VaultStore,
	cats core.CollateralStore,
	i18n *localizer.Localizer,
) core.Notifier {
	return &notifier{
		system:   system,
		assetz:   asset.Cache(assetz),
		messages: messages,
		vats:     vats,
		cats:     cats,
		i18n:     i18n,
	}
}

type notifier struct {
	system   *core.System
	assetz   core.AssetService
	messages core.MessageStore
	vats     core.VaultStore
	cats     core.CollateralStore
	i18n     *localizer.Localizer
}

func (n *notifier) localize(id string, args ...interface{}) string {
	return n.i18n.LocalizeOr(id, id, args...)
}

func (n *notifier) Auth(ctx context.Context, user *core.User) error {
	msg := n.localize("login_done")
	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(n.system.ClientID, user.MixinID),
		RecipientID:    user.MixinID,
		MessageID:      uuid.Modify(user.MixinID, user.AccessToken),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(msg)),
	}

	return n.messages.Create(ctx, []*core.Message{core.BuildMessage(req)})
}

func (n *notifier) Transaction(ctx context.Context, tx *core.Transaction) error {
	if tx.UserID == "" {
		return nil
	}

	switch tx.Action {
	case core.ActionVatOpen, core.ActionVatDeposit, core.ActionVatWithdraw, core.ActionVatPayback, core.ActionVatGenerate:
		return n.handleVatTx(ctx, tx)
	}

	return nil
}

func (n *notifier) Snapshot(ctx context.Context, transfer *core.Transfer, signedTx string) error {
	log := logger.FromContext(ctx)

	tx, err := mixin.TransactionFromRaw(signedTx)
	if err != nil {
		log.WithError(err).Debugln("decode transaction from raw tx failed")
		return nil
	}

	hash, err := tx.TransactionHash()
	if err != nil {
		return nil
	}

	traceID := mixinRawTransactionTraceId(hash.String(), 0)

	if len(transfer.Opponents) != 1 {
		log.Debugln("transfer opponents is not 1")
		return nil
	}

	asset, err := n.assetz.Find(ctx, transfer.AssetID)
	if err != nil {
		return err
	}

	card := mixin.AppCardMessage{
		AppID:       n.system.ClientID,
		IconURL:     asset.Logo,
		Title:       transfer.Amount.String(),
		Description: asset.Symbol,
		Action:      mixin.URL.Snapshots("", traceID),
	}
	data, _ := json.Marshal(card)

	recipientID := transfer.Opponents[0]
	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(n.system.ClientID, recipientID),
		RecipientID:    recipientID,
		MessageID:      traceID,
		Category:       mixin.MessageCategoryAppCard,
		Data:           base64.StdEncoding.EncodeToString(data),
	}

	return n.messages.Create(ctx, []*core.Message{core.BuildMessage(req)})
}
