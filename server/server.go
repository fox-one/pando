package server

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	s1 := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return s1.Shutdown(ctx)
		}
	})
	g.Go(func() error {
		return s1.ListenAndServe()
	})

	return g.Wait()
}
