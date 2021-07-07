package message

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/patrickmn/go-cache"
)

func Duplicate(messages core.MessageStore, exp time.Duration) core.MessageStore {
	return &duplicate{
		MessageStore: messages,
		c:            cache.New(exp, time.Minute),
	}
}

type duplicate struct {
	core.MessageStore
	c *cache.Cache
}

func (s *duplicate) Create(ctx context.Context, messages []*core.Message) error {
	var idx int

	for _, msg := range messages {
		if _, ok := s.c.Get(msg.MessageID); ok {
			continue
		}

		messages[idx] = msg
		idx++
	}

	if messages = messages[:idx]; len(messages) == 0 {
		return nil
	}

	if err := s.MessageStore.Create(ctx, messages); err != nil {
		return err
	}

	for _, msg := range messages {
		s.c.SetDefault(msg.MessageID, nil)
	}

	return nil
}
