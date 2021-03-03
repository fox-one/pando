package api

import (
	"net/http"
	"time"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/api/actions"
	"github.com/fox-one/pando/handler/api/system"
	"github.com/fox-one/pando/handler/auth"
	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pando/handler/rpc"
	"github.com/fox-one/pando/pkg/reversetwirp"
	"github.com/fox-one/pkg/property"
	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

func New(
	sessions core.Session,
	userz core.UserService,
	assets core.AssetStore,
	vaults core.VaultStore,
	flips core.FlipStore,
	properties property.Store,
	collaterals core.CollateralStore,
	transactions core.TransactionStore,
	walletz core.WalletService,
	notifier core.Notifier,
	system *core.System,
) *Server {
	return &Server{
		sessions:     sessions,
		userz:        userz,
		assets:       assets,
		vaults:       vaults,
		flips:        flips,
		properties:   properties,
		collaterals:  collaterals,
		transactions: transactions,
		walletz:      walletz,
		notifier:     notifier,
		system:       system,
	}
}

type Server struct {
	sessions     core.Session
	userz        core.UserService
	assets       core.AssetStore
	vaults       core.VaultStore
	flips        core.FlipStore
	properties   property.Store
	collaterals  core.CollateralStore
	transactions core.TransactionStore
	walletz      core.WalletService
	notifier     core.Notifier
	system       *core.System
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.HandleAuthentication(s.sessions))
	r.Use(render.WrapResponse(true))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, twirp.NotFoundError("not found"))
	})

	r.Get("/time", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		render.JSON(w, render.H{
			"iso":   t.Format(time.RFC3339),
			"epoch": t.Unix(),
		})
	})

	r.Get("/info", system.HandleInfo(s.system))
	r.Post("/login", auth.HandleOauth(s.userz, s.sessions, s.notifier))

	svr := rpc.New(s.assets, s.vaults, s.flips, s.properties, s.collaterals, s.transactions).TwirpServer()
	rt := reversetwirp.NewSingleTwirpServerProxy(svr)

	r.Route("/assets", func(r chi.Router) {
		r.Get("/", rt.Handle("ListAssets", nil))
		r.Get("/{id}", rt.Handle("ReadAsset", nil))
	})

	r.Route("/cats", func(r chi.Router) {
		r.Get("/", rt.Handle("ListCollaterals", nil))
		r.Get("/{id}", rt.Handle("FindCollateral", nil))
	})

	r.Route("/vats", func(r chi.Router) {
		r.Get("/", rt.Handle("ListVaults", nil))
		r.Get("/{id}", rt.Handle("FindVault", nil))
	})

	r.Route("/flips", func(r chi.Router) {
		r.Get("/", rt.Handle("ListFlips", nil))
		r.Get("/options", rt.Handle("ReadFlipOption", nil))
		r.Get("/{id}", rt.Handle("FindFlip", nil))
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/{id}", rt.Handle("FindTransaction", nil))
		r.Get("/", rt.Handle("ListTransactions", nil))
		r.Get("/cats/{collateral_id}", rt.Handle("ListTransactions", nil))
		r.Get("/vats/{vault_id}", rt.Handle("ListTransactions", nil))
		r.Get("/flips/{flip_id}", rt.Handle("ListTransactions", nil))

	})

	r.Route("/actions", func(r chi.Router) {
		r.Post("/", actions.HandleCreate(s.walletz, s.system))
	})

	return r
}
