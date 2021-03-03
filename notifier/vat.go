package notifier

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker/vat"
	"github.com/fox-one/pando/pkg/uuid"
)

type VatData struct {
	CatName   string
	GemSymbol string
	DaiSymbol string
	vat.Data
	Msg string
}

func (n *notifier) handleVatTx(ctx context.Context, tx *core.Transaction) error {
	var data VatData
	_ = tx.Data.Unmarshal(&data)

	v, err := n.vats.Find(ctx, tx.TargetID)
	if err != nil {
		return fmt.Errorf("vats.Find(%q) %w", tx.TargetID, err)
	}

	c, err := n.cats.Find(ctx, v.CollateralID)
	if err != nil {
		return fmt.Errorf("cats.Find(%q) %w", v.CollateralID, err)
	}

	data.CatName = c.Name

	gem, err := n.assetz.Find(ctx, c.Gem)
	if err != nil {
		return fmt.Errorf("assetz.Find(%q) %w", c.Gem, err)
	}

	data.GemSymbol = gem.Symbol

	dai, err := n.assetz.Find(ctx, c.Dai)
	if err != nil {
		return fmt.Errorf("assetz.Find(%q) %w", c.Dai, err)
	}

	data.DaiSymbol = dai.Symbol

	id := "vat"
	switch tx.Action {
	case core.ActionVatOpen:
		id = id + "_open"
	case core.ActionVatDeposit:
		id = id + "_deposit"
	case core.ActionVatWithdraw:
		id = id + "_withdraw"
	case core.ActionVatPayback:
		id = id + "_payback"
	case core.ActionVatGenerate:
		id = id + "_generate"
	default:
		return nil
	}

	if tx.Status == core.TxStatusSuccess {
		id = id + "_success"
	} else {
		id = id + "_failed"
	}

	msg := n.localize(id, data)
	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(n.system.ClientID, tx.UserID),
		RecipientID:    tx.UserID,
		MessageID:      uuid.Modify(tx.TraceID, "notifier"),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(msg)),
	}

	return n.messages.Create(ctx, []*core.Message{core.BuildMessage(req)})
}
