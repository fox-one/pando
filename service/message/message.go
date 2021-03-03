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

func (s *messageService) Send(ctx context.Context, messages []*core.Message) error {
	raws := make([]json.RawMessage, 0, len(messages))
	for _, msg := range messages {
		if msg.UserID == s.c.ClientID {
			continue
		}

		raws = append(raws, json.RawMessage(msg.Raw))
	}

	err := s.c.SendRawMessages(ctx, raws)

	// 如果 message.UserID 是机器人创建出来的账号，
	// 或者 conversation id 没有创建，发消息会报 10002
	// 忽略这种错误
	if mixin.IsErrorCodes(err, 10002) {
		return nil
	}

	return err
}

func (s *messageService) Meet(ctx context.Context, userID string) error {
	if userID == s.c.ClientID {
		return nil
	}

	_, err := s.c.CreateContactConversation(ctx, userID)
	return err
}
