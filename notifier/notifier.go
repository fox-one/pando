package notifier

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

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
	s := n.i18n.LocalizeOr(id, id, args...)
	return strings.TrimSpace(s)
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

	action := n.localize("Action" + tx.Action.String())
	data := TxData{
		Action:  action,
		Message: tx.Message,
	}

	id := "tx_abort"
	if tx.Status == core.TransactionStatusOk {
		switch tx.Action {
		case core.ActionVatOpen, core.ActionVatDeposit, core.ActionVatWithdraw, core.ActionVatPayback, core.ActionVatGenerate:
			if err := n.handleVatTx(ctx, tx, &data); err != nil {
				return err
			}
		}

		id = "tx_ok"
	}

	msg := n.localize(id, data)
	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(n.system.ClientID, tx.UserID),
		RecipientID:    tx.UserID,
		MessageID:      uuid.Modify(tx.TraceID, "notify"),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(msg)),
	}

	return n.messages.Create(ctx, []*core.Message{core.BuildMessage(req)})
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

	coin, err := n.assetz.Find(ctx, transfer.AssetID)
	if err != nil {
		return err
	}

	card := mixin.AppCardMessage{
		AppID:       n.system.ClientID,
		IconURL:     coin.Logo,
		Title:       transfer.Amount.String(),
		Description: coin.Symbol,
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
