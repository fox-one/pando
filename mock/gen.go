package mock

//go:generate go run github.com/golang/mock/mockgen -package=mock -destination=mock_gen.go github.com/fox-one/pando/core AssetStore,AssetService,CollateralStore,FlipStore,MessageStore,MessageService,Notifier,OracleStore,OracleService,ProposalStore,Parliament,Session,TransactionStore,UserStore,UserService,VaultStore,WalletStore,WalletService
