package datadog

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/metric"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/spf13/cast"
)

type Config struct {
	ConversationID string `valid:"uuid,required"`
	Interval       time.Duration
	Version        string
}

func New(
	wallets core.WalletStore,
	properties property.Store,
	messagez core.MessageService,
	cfg Config,
) *Datadog {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	return &Datadog{
		wallets:        wallets,
		properties:     properties,
		messagez:       messagez,
		interval:       cfg.Interval,
		launchAt:       time.Now(),
		conversationID: cfg.ConversationID,
		version:        cfg.Version,
	}
}

type Datadog struct {
	wallets        core.WalletStore
	properties     property.Store
	messagez       core.MessageService
	interval       time.Duration
	launchAt       time.Time
	conversationID string
	version        string
}

func (w *Datadog) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "datadog")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-time.After(dur):
			if err := w.run(ctx); err != nil {
				dur = time.Second
			} else {
				dur = t.Truncate(w.interval).Add(w.interval).Sub(t)
			}
		}
	}
}

func (w *Datadog) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	var groups []metric.Group

	var report bool

	// wallets
	{
		lastOutputID, err := w.wallets.CountOutputs(ctx)
		if err != nil {
			log.WithError(err).Errorln("wallets.CountOutputs")
			return err
		}

		unhandled, err := w.wallets.CountUnhandledTransfers(ctx)
		if err != nil {
			log.WithError(err).Errorln("wallets.CountUnhandledTransfers")
			return err
		}

		groups = append(groups, metric.Group{
			Name: "wallets",
			Entries: []metric.Entry{
				{
					Name:  "unhandled_transfers",
					Value: cast.ToString(unhandled),
				},
				{
					Name:  "last_output_id",
					Value: cast.ToString(lastOutputID),
				},
			},
		})

		report = unhandled > 0 || report
	}

	// properties
	{
		items, err := w.properties.List(ctx)
		if err != nil {
			log.WithError(err).Errorln("properties.List")
			return err
		}

		group := metric.Group{Name: "properties"}
		for k, v := range items {
			group.Entries = append(group.Entries, metric.Entry{
				Name:  k,
				Value: v.String(),
			})

			if k == "sync_checkpoint" {
				t := v.Time()
				report = (!t.IsZero() && time.Since(t) > 5*time.Minute) || report
			}
		}

		groups = append(groups, group)
	}

	// system
	{
		uptime := time.Since(w.launchAt)
		report = uptime < time.Minute || report

		groups = append(groups, metric.Group{
			Name: "system",
			Entries: []metric.Entry{
				{
					Name:  "uptime",
					Value: uptime.String(),
				},
				{
					Name:  "version",
					Value: w.version,
				},
			},
		})
	}

	var b bytes.Buffer
	metric.Render(&b, groups)

	if !report {
		return nil
	}

	msg := core.BuildMessage(&mixin.MessageRequest{
		ConversationID: w.conversationID,
		MessageID:      uuid.New(),
		Category:       mixin.MessageCategoryPlainPost,
		Data:           base64.StdEncoding.EncodeToString(b.Bytes()),
	})

	if err := w.messagez.Send(ctx, []*core.Message{msg}, false); err != nil {
		log.WithError(err).Errorln("messagez.Send")
		return err
	}

	return nil
}
