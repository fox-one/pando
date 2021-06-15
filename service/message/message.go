package message

import (
	"context"
	"encoding/json"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
)

func New(client *mixin.Client) core.MessageService {
	return &messageService{c: client}
}

type messageService struct {
	c *mixin.Client
}

func parseSendMessageErr(err error) error {
	if err != nil && mixin.IsErrorCodes(err, 10002) {
		return nil
	}

	return err
}

func (s *messageService) Send(ctx context.Context, messages []*core.Message, batch bool) error {
	raws := make([]json.RawMessage, 0, len(messages))
	for _, msg := range messages {
		if msg.UserID == s.c.ClientID {
			continue
		}

		raws = append(raws, json.RawMessage(msg.Raw))
	}

	if batch {
		return parseSendMessageErr(s.c.SendRawMessages(ctx, raws))
	}

	for _, raw := range raws {
		if err := parseSendMessageErr(s.c.SendRawMessage(ctx, raw)); err != nil {
			return err
		}
	}

	return nil
}

func (s *messageService) Meet(ctx context.Context, userID string) error {
	if userID == s.c.ClientID {
		return nil
	}

	if _, err := s.c.CreateContactConversation(ctx, userID); err != nil && !mixin.IsErrorCodes(err, 10002) {
		return err
	}

	return nil
}
