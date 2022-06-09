package call

import (
	"context"
	"net/http"

	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	api "github.com/fox-one/pando/handler/rpc/pando"
	"github.com/twitchtv/twirp"
)

func RPC() api.Pando {
	return api.NewPandoProtobufClient(cfg.GetApiHost(), client.GetClient())
}

func WithToken(ctx context.Context, token string) context.Context {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	ctx, _ = twirp.WithHTTPRequestHeaders(ctx, header)
	return ctx
}
