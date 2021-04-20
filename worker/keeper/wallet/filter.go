package wallet

import (
	"context"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/patrickmn/go-cache"
)

func FilterTrace(walletz core.WalletService, exp time.Duration) core.WalletService {
	return &filterTrace{
		WalletService: walletz,
		traces:        cache.New(exp, time.Minute),
		nonce:         uuid.New(),
	}
}

type filterTrace struct {
	core.WalletService
	traces *cache.Cache
	nonce  string
}

func remixTrace(transfer *core.Transfer, nonce string) *core.Transfer {
	t := new(core.Transfer)
	*t = *transfer
	t.TraceID = uuid.Modify(transfer.TraceID, nonce)
	return t
}

func (s *filterTrace) HandleTransfer(ctx context.Context, transfer *core.Transfer) error {
	transfer = remixTrace(transfer, s.nonce)

	if _, ok := s.traces.Get(transfer.TraceID); ok {
		return nil
	}

	if err := s.WalletService.HandleTransfer(ctx, transfer); err != nil {
		return err
	}

	s.traces.SetDefault(transfer.TraceID, nil)
	return nil
}
