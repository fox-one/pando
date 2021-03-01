package rpc

import (
	"context"
	"net/http"
	"sort"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/auth"
	"github.com/fox-one/pando/handler/request"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/handler/rpc/view"
	"github.com/fox-one/pando/pkg/maker/flip"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/fox-one/pkg/store"
	"github.com/spf13/cast"
	"github.com/twitchtv/twirp"
)

func New(
	assets core.AssetStore,
	vaults core.VaultStore,
	flips core.FlipStore,
	properties property.Store,
	collaterals core.CollateralStore,
	transactions core.TransactionStore,
) *Server {
	return &Server{
		assets:       assets,
		vaults:       vaults,
		flips:        flips,
		properties:   properties,
		collaterals:  collaterals,
		transactions: transactions,
	}
}

type Server struct {
	assets       core.AssetStore
	vaults       core.VaultStore
	flips        core.FlipStore
	properties   property.Store
	collaterals  core.CollateralStore
	transactions core.TransactionStore
}

func (s *Server) TwirpServer() api.TwirpServer {
	opts := []interface{}{
		twirp.WithServerJSONSkipDefaults(false),
	}

	return api.NewPandoServer(s, opts...)
}

func (s *Server) Handle(sessions core.Session) http.Handler {
	return auth.HandleAuthentication(sessions)(s.TwirpServer())
}

func (s *Server) ReadAsset(ctx context.Context, req *api.Req_ReadAsset) (*api.Asset, error) {
	asset, err := s.assets.Find(ctx, req.Id)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, twirp.NotFoundError("asset not found")
		}

		logger.FromContext(ctx).WithError(err).Errorf("rpc: assets.Find(%s)", req.Id)
		return nil, err
	}

	chain := asset
	if asset.ChainID != chain.ID {
		chain, err = s.assets.Find(ctx, asset.ChainID)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorf("rpc: assets.Find(%s)", asset.ChainID)
			return nil, err
		}
	}

	return view.Asset(asset, chain), nil
}

func (s *Server) ListAssets(ctx context.Context, _ *api.Req_ListAssets) (*api.Resp_ListAssets, error) {
	assets, err := s.assets.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: assets.ListAll")
		return nil, err
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Price.GreaterThan(assets[j].Price)
	})

	chains := make(map[string]*core.Asset, 32)
	for _, asset := range assets {
		if asset.ID == asset.ChainID {
			chains[asset.ID] = asset
		}
	}

	resp := &api.Resp_ListAssets{}
	for _, asset := range assets {
		resp.Assets = append(resp.Assets, view.Asset(asset, chains[asset.ChainID]))
	}

	return resp, nil
}

func (s *Server) ListCollaterals(ctx context.Context, _ *api.Req_ListCollaterals) (*api.Resp_ListCollaterals, error) {
	cats, err := s.collaterals.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: collaterals.List")
		return nil, err
	}

	resp := &api.Resp_ListCollaterals{}
	for _, cat := range cats {
		resp.Collaterals = append(resp.Collaterals, view.Collateral(cat))
	}

	return resp, nil
}

func (s *Server) FindCollateral(ctx context.Context, req *api.Req_FindCollateral) (*api.Collateral, error) {
	cat, err := s.collaterals.Find(ctx, req.Id)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: collaterals.Find")
		return nil, err
	}

	if cat.ID == 0 {
		return nil, twirp.NotFoundError("cat not init")
	}

	return view.Collateral(cat), nil
}

func (s *Server) ListVaults(ctx context.Context, _ *api.Req_ListVaults) (*api.Resp_ListVaults, error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		logger.FromContext(ctx).Debugln("rpc: authentication required")
		return nil, twirp.NewError(twirp.Unauthenticated, "authentication required")
	}

	vats, err := s.vaults.ListUser(ctx, user.MixinID)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: vaults.ListUser")
		return nil, err
	}

	resp := &api.Resp_ListVaults{}
	for _, vat := range vats {
		resp.Vaults = append(resp.Vaults, view.Vault(vat))
	}

	return resp, nil
}

func (s *Server) FindVault(ctx context.Context, req *api.Req_FindVault) (*api.Vault, error) {
	vat, err := s.vaults.Find(ctx, req.Id)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: vaults.Find")
		return nil, err
	}

	if vat.ID == 0 {
		return nil, twirp.NotFoundError("vat not init")
	}

	return view.Vault(vat), nil
}

func (s *Server) FindFlip(ctx context.Context, req *api.Req_FindFlip) (*api.Flip, error) {
	flip, err := s.flips.Find(ctx, req.Id)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: flips.Find")
		return nil, err
	}

	if flip.ID == 0 {
		return nil, twirp.NotFoundError("flip not init")
	}

	if user, ok := request.UserFrom(ctx); !ok || user.MixinID != flip.Guy {
		flip.Guy = ""
	}

	return view.Flip(flip), nil
}

func (s *Server) ListFlips(ctx context.Context, req *api.Req_ListFlips) (*api.Resp_ListFlips, error) {
	fromID := cast.ToInt64(req.Cursor)
	limit := 50
	if l := int(req.Limit); l > 0 && l < limit {
		limit = l
	}

	flips, err := s.flips.List(ctx, fromID, limit+1)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: flips.List")
		return nil, err
	}

	resp := &api.Resp_ListFlips{
		Pagination: &api.Pagination{},
	}

	for idx, f := range flips {
		f.Guy = ""
		resp.Flips = append(resp.Flips, view.Flip(f))

		if idx == limit-1 {
			resp.Pagination.NextCursor = cast.ToString(f.ID)
			resp.Pagination.HasNext = true
			break
		}
	}

	return resp, nil
}

func (s *Server) ReadFlipOption(ctx context.Context, _ *api.Req_ReadFlipOption) (*api.FlipOption, error) {
	opt, err := flip.ReadOptions(ctx, s.properties)
	if err != nil {
		return nil, err
	}

	return view.FlipOption(opt), nil
}

func (s *Server) FindTransaction(ctx context.Context, req *api.Req_FindTransaction) (*api.Transaction, error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		logger.FromContext(ctx).Debugln("rpc: authentication required")
		return nil, twirp.NewError(twirp.Unauthenticated, "authentication required")
	}

	tx, err := s.transactions.FindFollow(ctx, user.MixinID, req.Follow)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, twirp.NotFoundError("transaction not found")
		}

		return nil, err
	}

	return view.Transaction(tx), nil
}

func (s *Server) ListTransactions(ctx context.Context, req *api.Req_ListTransactions) (*api.Resp_ListTransactions, error) {
	fromID := cast.ToInt64(req.Cursor)
	limit := 50
	if l := int(req.Limit); l > 0 && l < limit {
		limit = l
	}

	transactions, err := s.transactions.ListTarget(ctx, req.Target, fromID, limit+1)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: transactions.ListTarget")
		return nil, err
	}

	resp := &api.Resp_ListTransactions{
		Pagination: &api.Pagination{},
	}

	for idx, t := range transactions {
		resp.Transactions = append(resp.Transactions, view.Transaction(t))

		if idx == limit-1 {
			resp.Pagination.NextCursor = cast.ToString(t.ID)
			resp.Pagination.HasNext = true
			break
		}
	}

	return resp, nil
}
