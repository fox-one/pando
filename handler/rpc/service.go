package rpc

import (
	"context"
	"sort"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/handler/request"
	"github.com/fox-one/pando/handler/rpc/api"
	"github.com/fox-one/pando/handler/rpc/view"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store"
	"github.com/spf13/cast"
	"github.com/twitchtv/twirp"
)

func New(
	assets core.AssetStore,
	vaults core.VaultStore,
	collaterals core.CollateralStore,
	transactions core.TransactionStore,
) api.TwirpServer {
	svc := &service{
		assets:       assets,
		vaults:       vaults,
		collaterals:  collaterals,
		transactions: transactions,
	}

	opts := []interface{}{
		twirp.WithServerJSONSkipDefaults(false),
	}

	return api.NewPandoServer(svc, opts...)
}

type service struct {
	assets       core.AssetStore
	vaults       core.VaultStore
	collaterals  core.CollateralStore
	transactions core.TransactionStore
}

func (s *service) ReadAsset(ctx context.Context, req *api.Req_ReadAsset) (*api.Asset, error) {
	asset, err := s.assets.Find(ctx, req.Id)
	if err != nil {
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

func (s *service) ListAssets(ctx context.Context, _ *api.Req_ListAssets) (*api.Resp_ListAssets, error) {
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

func (s *service) ListCollaterals(ctx context.Context, _ *api.Req_ListCollaterals) (*api.Resp_ListCollaterals, error) {
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

func (s *service) FindCollateral(ctx context.Context, req *api.Req_FindCollateral) (*api.Collateral, error) {
	cat, err := s.collaterals.Find(ctx, req.Id)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: collaterals.Find")
		return nil, err
	}

	return view.Collateral(cat), nil
}

func (s *service) ListVaults(ctx context.Context, _ *api.Req_ListVaults) (*api.Resp_ListVaults, error) {
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

func (s *service) FindVault(ctx context.Context, req *api.Req_FindVault) (*api.Vault, error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		logger.FromContext(ctx).Debugln("rpc: authentication required")
		return nil, twirp.NewError(twirp.Unauthenticated, "authentication required")
	}

	vat, err := s.vaults.Find(ctx, req.Id)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("rpc: vaults.Find")
		return nil, err
	}

	if vat.UserID != user.MixinID {
		return nil, twirp.NewError(twirp.Unauthenticated, "authentication required")
	}

	return view.Vault(vat), nil
}

func (s *service) FindTransaction(ctx context.Context, req *api.Req_FindTransaction) (*api.Transaction, error) {
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

func (s *service) ListTransactions(ctx context.Context, req *api.Req_ListTransactions) (*api.Resp_ListTransactions, error) {
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
