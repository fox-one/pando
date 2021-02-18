package messenger

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pkg/logger"
)

func New(messages core.MessageStore, messagez core.MessageService) *Messenger {
	return &Messenger{
		messages: messages,
		messagez: messagez,
	}
}

type Messenger struct {
	messages core.MessageStore
	messagez core.MessageService
}

func (w *Messenger) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "messenger")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Messenger) run(ctx context.Context) error {
	log := logger.FromContext(ctx)
	const Limit = 300
	const Batch = 70

	messages, err := w.messages.List(ctx, Limit)
	if err != nil {
		log.WithError(err).Error("messengers.ListPair")
		return err
	}

	if len(messages) == 0 {
		return errors.New("list messages: EOF")
	}

	filter := make(map[string]bool)
	var idx int

	for _, msg := range messages {
		if filter[msg.UserID] {
			continue
		}

		messages[idx] = msg
		filter[msg.UserID] = true
		idx++

		if idx >= Batch {
			break
		}
	}

	messages = messages[:idx]
	if err := w.messagez.Send(ctx, messages); err != nil {
		log.WithError(err).Error("messagez.Send")
		return err
	}

	if err := w.messages.Delete(ctx, messages); err != nil {
		log.WithError(err).Error("messagez.Delete")
		return err
	}

	return nil
}
