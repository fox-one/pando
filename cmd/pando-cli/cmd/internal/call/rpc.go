package call

import (
	"github.com/fox-one/pando/cmd/pando-cli/cmd/internal/cfg"
	"github.com/fox-one/pando/handler/rpc/api"
)

func RPC() api.Pando {
	return api.NewPandoProtobufClient(cfg.GetApiHost(), client.GetClient())
}
