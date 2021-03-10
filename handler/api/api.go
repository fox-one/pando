package api

import (
	"net/http"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/api/actions"
	"github.com/fox-one/pando/handler/api/system"
	"github.com/fox-one/pando/handler/api/user"
	"github.com/fox-one/pando/handler/auth"
	"github.com/fox-one/pando/handler/render"
	"github.com/fox-one/pando/handler/rpc"
	"github.com/fox-one/pando/pkg/reversetwirp"
	"github.com/go-chi/chi"
	"github.com/twitchtv/twirp"
)

func New(
	sessions core.Session,
	userz core.UserService,
	assets core.AssetStore,
	vaults core.VaultStore,
	flips core.FlipStore,
	collaterals core.CollateralStore,
	transactions core.TransactionStore,
	walletz core.WalletService,
	notifier core.Notifier,
	oracles core.OracleStore,
	system *core.System,
) *Server {
	return &Server{
		sessions:     sessions,
		userz:        userz,
		assets:       assets,
		vaults:       vaults,
		flips:        flips,
		collaterals:  collaterals,
		transactions: transactions,
		walletz:      walletz,
		notifier:     notifier,
		oracles:      oracles,
		system:       system,
	}
}

type Server struct {
	sessions     core.Session
	userz        core.UserService
	assets       core.AssetStore
	vaults       core.VaultStore
	flips        core.FlipStore
	collaterals  core.CollateralStore
	transactions core.TransactionStore
	walletz      core.WalletService
	notifier     core.Notifier
	oracles      core.OracleStore
	system       *core.System
}

// @title Pando API
// @version 1.0
// @description Pando Api Doc
// @termsOfService http://pando.im/terms/

// @contact.name API Support
// @contact.url http://www.pando.im/support
// @contact.email support@pando.im

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes https
// @host pando-test-api.fox.one
// @BasePath /api
// @query.collection.format multi
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.HandleAuthentication(s.sessions))
	r.Use(render.WrapResponse(true))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, twirp.NotFoundError("not found"))
	})

	r.Get("/time", system.HandleTime())
	r.Get("/info", system.HandleInfo(s.system))

	r.Post("/login", user.HandleOauth(s.userz, s.sessions, s.notifier))

	svr := rpc.New(s.assets, s.vaults, s.flips, s.oracles, s.collaterals, s.transactions).TwirpServer()
	rt := reversetwirp.NewSingleTwirpServerProxy(svr)

	r.Route("/assets", func(r chi.Router) {
		r.Get("/", rt.Handle("ListAssets", nil))
		r.Get("/{id}", rt.Handle("FindAsset", nil))
	})

	r.Route("/oracles", func(r chi.Router) {
		r.Get("/", rt.Handle("ListOracles", nil))
		r.Get("/{id}", rt.Handle("FindOracle", nil))
	})

	r.Route("/cats", func(r chi.Router) {
		r.Get("/", rt.Handle("ListCollaterals", nil))
		r.Get("/{id}", rt.Handle("FindCollateral", nil))
	})

	r.Route("/vats", func(r chi.Router) {
		r.Get("/", rt.Handle("ListVaults", nil))
		r.Get("/{id}", rt.Handle("FindVault", nil))
	})

	r.Route("/me", func(r chi.Router) {
		r.Get("/vats", rt.Handle("ListMyVaults", nil))
	})

	r.Route("/flips", func(r chi.Router) {
		r.Get("/", rt.Handle("ListFlips", nil))
		r.Get("/{id}", rt.Handle("FindFlip", nil))
		r.Get("/{id}/events", rt.Handle("ListFlipEvents", nil))
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/{id}", rt.Handle("FindTransaction", nil))
		r.Get("/", rt.Handle("ListTransactions", nil))
	})

	r.Route("/actions", func(r chi.Router) {
		r.Post("/", actions.HandleCreate(s.walletz, s.system))
	})

	return r
}
