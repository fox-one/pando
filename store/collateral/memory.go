package collateral

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/pando/core"
)

func Memory() core.CollateralStore {
	return &memoryStore{
		m: map[string]core.Collateral{},
	}
}

type memoryStore struct {
	m  map[string]core.Collateral
	id int64
}

func (s *memoryStore) Create(ctx context.Context, cat *core.Collateral) error {
	if _, ok := s.m[cat.TraceID]; ok {
		return fmt.Errorf("cat with trace id %q already existed", cat.TraceID)
	}

	s.id += 1
	cat.ID = s.id
	cat.CreatedAt = time.Now()
	s.m[cat.TraceID] = *cat
	return nil
}

func (s *memoryStore) Update(ctx context.Context, cat *core.Collateral, version int64) error {
	if cat.Version >= version {
		return nil
	}

	if _, ok := s.m[cat.TraceID]; !ok {
		return fmt.Errorf("cat with trace id %q not existed", cat.TraceID)
	}

	cat.Version = version
	s.m[cat.TraceID] = *cat
	return nil
}

func (s *memoryStore) Find(ctx context.Context, traceID string) (*core.Collateral, error) {
	cat := s.m[traceID]
	return &cat, nil
}

func (s *memoryStore) List(ctx context.Context) ([]*core.Collateral, error) {
	var cats []*core.Collateral
	for key := range s.m {
		cat := s.m[key]
		cats = append(cats, &cat)
	}

	return cats, nil
}
