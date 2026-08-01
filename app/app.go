package app

import (
	"encoding/json"
	"fmt"
	"io"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/cosmos/ibc-go/v8/modules/core"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/spf13/cast"

	"github.com/bizchain/blockchain/x/agen"
	agenkeeper "github.com/bizchain/blockchain/x/agen/keeper"
	agentypes "github.com/bizchain/blockchain/x/agen/types"
	"github.com/bizchain/blockchain/x/asuransi"
	asuransikeeper "github.com/bizchain/blockchain/x/asuransi/keeper"
	asuransitypes "github.com/bizchain/blockchain/x/asuransi/types"
	"github.com/bizchain/blockchain/x/audit"
	auditkeeper "github.com/bizchain/blockchain/x/audit/keeper"
	audittypes "github.com/bizchain/blockchain/x/audit/types"
	"github.com/bizchain/blockchain/x/dao"
	daokeeper "github.com/bizchain/blockchain/x/dao/keeper"
	daotypes "github.com/bizchain/blockchain/x/dao/types"
	"github.com/bizchain/blockchain/x/hotel"
	hotelkeeper "github.com/bizchain/blockchain/x/hotel/keeper"
	hoteltypes "github.com/bizchain/blockchain/x/hotel/types"
	"github.com/bizchain/blockchain/x/jamaah"
	jamaahkeeper "github.com/bizchain/blockchain/x/jamaah/keeper"
	jamaahtypes "github.com/bizchain/blockchain/x/jamaah/types"
	"github.com/bizchain/blockchain/x/keberangkatan"
	keberangkatankeeper "github.com/bizchain/blockchain/x/keberangkatan/keeper"
	keberangkatantypes "github.com/bizchain/blockchain/x/keberangkatan/types"
	"github.com/bizchain/blockchain/x/oleholeh"
	oleholehkeeper "github.com/bizchain/blockchain/x/oleholeh/keeper"
	oleholehtypes "github.com/bizchain/blockchain/x/oleholeh/types"
	"github.com/bizchain/blockchain/x/paket"
	paketkeeper "github.com/bizchain/blockchain/x/paket/keeper"
	pakettypes "github.com/bizchain/blockchain/x/paket/types"
	"github.com/bizchain/blockchain/x/pembayaran"
	pembayarankeeper "github.com/bizchain/blockchain/x/pembayaran/keeper"
	pembayarantypes "github.com/bizchain/blockchain/x/pembayaran/types"
	"github.com/bizchain/blockchain/x/pos"
	poskeeper "github.com/bizchain/blockchain/x/pos/keeper"
	postypes "github.com/bizchain/blockchain/x/pos/types"
	"github.com/bizchain/blockchain/x/referral"
	referralkeeper "github.com/bizchain/blockchain/x/referral/keeper"
	referraltypes "github.com/bizchain/blockchain/x/referral/types"
	"github.com/bizchain/blockchain/x/reward"
	rewardkeeper "github.com/bizchain/blockchain/x/reward/keeper"
	rewardtypes "github.com/bizchain/blockchain/x/reward/types"
	"github.com/bizchain/blockchain/x/ticket"
	ticketkeeper "github.com/bizchain/blockchain/x/ticket/keeper"
	tickettypes "github.com/bizchain/blockchain/x/ticket/types"
	"github.com/bizchain/blockchain/x/visa"
	visakeeper "github.com/bizchain/blockchain/x/visa/keeper"
	visatypes "github.com/bizchain/blockchain/x/visa/types"
)

const (
	Name         = "BizChain"
	BondDenom    = "upoint"
	Bech32Prefix = "point"
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	// IBC transfer module account (required by ibctransferkeeper.NewKeeper)
	ibctransfertypes.ModuleName: nil,
}

func GetBech32Prefix() string {
	return Bech32Prefix
}

// EncodingConfig specifies the concrete encoding types for this app.
type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// App extends the Cosmos SDK App
type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry
	txConfig          client.TxConfig

	// Keepers
	AccountKeeper       authkeeper.AccountKeeper
	BankKeeper          bankkeeper.BaseKeeper
	StakingKeeper       *stakingkeeper.Keeper
	ConsensusKeeper     consensuskeeper.Keeper
	ParamsKeeper        paramskeeper.Keeper
	CrisisKeeper        *crisiskeeper.Keeper
	CapabilityKeeper    *capabilitykeeper.Keeper
	UpgradeKeeper       *upgradekeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	ScopedIBCKeeper     capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	POSKeeper           poskeeper.Keeper
	JamaahKeeper        jamaahkeeper.Keeper
	PaketKeeper         paketkeeper.Keeper
	PembayaranKeeper    pembayarankeeper.Keeper
	VisaKeeper          visakeeper.Keeper
	HotelKeeper         hotelkeeper.Keeper
	TicketKeeper        ticketkeeper.Keeper
	ReferralKeeper      referralkeeper.Keeper
	RewardKeeper        rewardkeeper.Keeper
	AuditKeeper         auditkeeper.Keeper
	KeberangkatanKeeper keberangkatankeeper.Keeper
	AgenKeeper          agenkeeper.Keeper
	OlehOlehKeeper      oleholehkeeper.Keeper
	AsuransiKeeper      asuransikeeper.Keeper
	DaoKeeper           daokeeper.Keeper

	// Module Manager
	mm           *module.Manager
	configurator module.Configurator

	// BasicModuleManager holds codec-carrying AppModuleBasic values extracted
	// from the module manager. It is used for CLI command wiring so that
	// module GetTxCmd/GetQueryCmd (which dereference the codec) work correctly.
	BasicModuleManager module.BasicManager
}

// NewApp creates a new BizChain App
func NewApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	encodingConfig := MakeEncodingConfig()
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(encodingConfig.TxConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		consensusparamtypes.StoreKey,
		paramstypes.StoreKey,
		crisistypes.StoreKey,
		capabilitytypes.StoreKey,
		upgradetypes.StoreKey,
		ibcexported.StoreKey,
		ibctransfertypes.StoreKey,
		postypes.StoreKey,
		jamaahtypes.StoreKey,
		pakettypes.StoreKey,
		pembayarantypes.StoreKey,
		visatypes.StoreKey,
		hoteltypes.StoreKey,
		tickettypes.StoreKey,
		referraltypes.StoreKey,
		rewardtypes.StoreKey,
		audittypes.StoreKey,
		keberangkatantypes.StoreKey,
		agentypes.StoreKey,
		oleholehtypes.StoreKey,
		asuransitypes.StoreKey,
		daotypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(
		capabilitytypes.MemStoreKey,
		postypes.MemStoreKey,
		jamaahtypes.MemStoreKey,
		pakettypes.MemStoreKey,
		pembayarantypes.MemStoreKey,
		visatypes.MemStoreKey,
		hoteltypes.MemStoreKey,
		tickettypes.MemStoreKey,
		referraltypes.MemStoreKey,
		rewardtypes.MemStoreKey,
		audittypes.MemStoreKey,
		keberangkatantypes.MemStoreKey,
		agentypes.MemStoreKey,
		oleholehtypes.MemStoreKey,
		asuransitypes.MemStoreKey,
		daotypes.MemStoreKey,
	)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		txConfig:          encodingConfig.TxConfig,
	}

	// Initialize keepers
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// Set the BaseApp's parameter store
	app.ConsensusKeeper = consensuskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(app.ConsensusKeeper.ParamsStore)

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(GetBech32Prefix()),
		GetBech32Prefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedAddresses(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(GetBech32Prefix()+"valoper"),
		authcodec.NewBech32Codec(GetBech32Prefix()+"valcons"),
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.AccountKeeper.AddressCodec(),
	)

	// IBC: capability keeper
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.CapabilityKeeper.Seal()

	// IBC: upgrade keeper (required by IBC core)
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// IBC: core keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// IBC: transfer keeper
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	// Initialize POS module keeper
	app.POSKeeper = *poskeeper.NewKeeper(
		appCodec,
		keys[postypes.StoreKey],
		memKeys[postypes.MemStoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Initialize Haji/Umroh module keepers
	app.JamaahKeeper = *jamaahkeeper.NewKeeper(appCodec, keys[jamaahtypes.StoreKey], memKeys[jamaahtypes.MemStoreKey], authority)
	app.PaketKeeper = *paketkeeper.NewKeeper(appCodec, keys[pakettypes.StoreKey], memKeys[pakettypes.MemStoreKey], authority)
	app.PembayaranKeeper = *pembayarankeeper.NewKeeper(appCodec, keys[pembayarantypes.StoreKey], memKeys[pembayarantypes.MemStoreKey], authority)
	app.VisaKeeper = *visakeeper.NewKeeper(appCodec, keys[visatypes.StoreKey], memKeys[visatypes.MemStoreKey], authority)
	app.HotelKeeper = *hotelkeeper.NewKeeper(appCodec, keys[hoteltypes.StoreKey], memKeys[hoteltypes.MemStoreKey], authority)
	app.TicketKeeper = *ticketkeeper.NewKeeper(appCodec, keys[tickettypes.StoreKey], memKeys[tickettypes.MemStoreKey], authority)
	app.ReferralKeeper = *referralkeeper.NewKeeper(appCodec, keys[referraltypes.StoreKey], memKeys[referraltypes.MemStoreKey], authority)
	app.RewardKeeper = *rewardkeeper.NewKeeper(appCodec, keys[rewardtypes.StoreKey], memKeys[rewardtypes.MemStoreKey], authority)
	app.AuditKeeper = *auditkeeper.NewKeeper(appCodec, keys[audittypes.StoreKey], memKeys[audittypes.MemStoreKey], authority)
	app.KeberangkatanKeeper = *keberangkatankeeper.NewKeeper(appCodec, keys[keberangkatantypes.StoreKey], memKeys[keberangkatantypes.MemStoreKey], authority)
	app.AgenKeeper = *agenkeeper.NewKeeper(appCodec, keys[agentypes.StoreKey], memKeys[agentypes.MemStoreKey], authority)
	app.OlehOlehKeeper = *oleholehkeeper.NewKeeper(appCodec, keys[oleholehtypes.StoreKey], memKeys[oleholehtypes.MemStoreKey], authority)
	app.AsuransiKeeper = *asuransikeeper.NewKeeper(appCodec, keys[asuransitypes.StoreKey], memKeys[asuransitypes.MemStoreKey], authority)
	app.DaoKeeper = *daokeeper.NewKeeper(appCodec, keys[daotypes.StoreKey], memKeys[daotypes.MemStoreKey], authority)

	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// Module Manager
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		consensus.NewAppModule(appCodec, app.ConsensusKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
		ibctm.NewAppModule(),
		pos.NewAppModule(appCodec, app.POSKeeper),
		jamaah.NewAppModule(appCodec, app.JamaahKeeper),
		paket.NewAppModule(appCodec, app.PaketKeeper),
		pembayaran.NewAppModule(appCodec, app.PembayaranKeeper),
		visa.NewAppModule(appCodec, app.VisaKeeper),
		hotel.NewAppModule(appCodec, app.HotelKeeper),
		ticket.NewAppModule(appCodec, app.TicketKeeper),
		referral.NewAppModule(appCodec, app.ReferralKeeper),
		reward.NewAppModule(appCodec, app.RewardKeeper),
		audit.NewAppModule(appCodec, app.AuditKeeper),
		keberangkatan.NewAppModule(appCodec, app.KeberangkatanKeeper),
		agen.NewAppModule(appCodec, app.AgenKeeper),
		oleholeh.NewAppModule(appCodec, app.OlehOlehKeeper),
		asuransi.NewAppModule(appCodec, app.AsuransiKeeper),
		dao.NewAppModule(appCodec, app.DaoKeeper),
	)

	// Build a codec-carrying BasicManager from the module manager so that
	// module CLI commands (GetTxCmd/GetQueryCmd) can be wired correctly.
	app.BasicModuleManager = make(module.BasicManager)
	for name, m := range app.mm.Modules {
		if amb, ok := m.(module.AppModuleBasic); ok {
			app.BasicModuleManager[name] = amb
		}
	}

	// Genesis order
	genesisModuleOrder := []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		genutiltypes.ModuleName,
		crisistypes.ModuleName,
		capabilitytypes.ModuleName,
		upgradetypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		ibctm.ModuleName,
		consensusparamtypes.ModuleName,
		postypes.ModuleName,
		jamaahtypes.ModuleName,
		pakettypes.ModuleName,
		pembayarantypes.ModuleName,
		visatypes.ModuleName,
		hoteltypes.ModuleName,
		tickettypes.ModuleName,
		referraltypes.ModuleName,
		rewardtypes.ModuleName,
		audittypes.ModuleName,
		keberangkatantypes.ModuleName,
		agentypes.ModuleName,
		oleholehtypes.ModuleName,
		asuransitypes.ModuleName,
		daotypes.ModuleName,
	}
	app.mm.SetOrderInitGenesis(genesisModuleOrder...)
	app.mm.SetOrderExportGenesis(genesisModuleOrder...)

	app.mm.SetOrderBeginBlockers(genesisModuleOrder...)
	app.mm.SetOrderEndBlockers(genesisModuleOrder...)

	app.mm.RegisterInvariants(app.CrisisKeeper)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	if err := app.mm.RegisterServices(app.configurator); err != nil {
		panic(err)
	}

	// Initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// Initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// Set AnteHandler with zero-fee support
	app.setAnteHandler(encodingConfig)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(fmt.Errorf("failed to load latest version: %w", err))
		}
	}

	return app
}

func (app *App) setAnteHandler(encodingConfig EncodingConfig) {
	handlerOptions := ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		FeegrantKeeper:  nil,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
	}

	handler, err := ante.NewAnteHandler(handlerOptions)
	if err != nil {
		panic(fmt.Sprintf("failed to create ante handler: %s", err))
	}
	app.SetAnteHandler(handler)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *App) ModuleManager() *module.Manager {
	return app.mm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig serverconfig.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new CometBFT queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	cmtApp := server.NewCometABCIWrapper(app)
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		cmtApp.Query,
	)
}

func (app *App) RegisterNodeService(clientCtx client.Context, cfg serverconfig.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *App) ExportAppStateAndValidators(forZeroHeight bool, jailAllowedAddrs, modulesToExport []string) (servertypes.ExportedApp, error) {
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
	}
	genState, err := app.mm.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	appState, err := json.MarshalIndent(genState, "", " ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	validators, err := staking.WriteValidators(ctx, app.StakingKeeper)
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}

func init() {
	setConfig()
}

func setConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32Prefix, Bech32Prefix+"pub")
	config.SetBech32PrefixForValidator(Bech32Prefix+"valoper", Bech32Prefix+"valoperpub")
	config.SetBech32PrefixForConsensusNode(Bech32Prefix+"valcons", Bech32Prefix+"valconspub")
	config.SetPurpose(44)
	config.SetCoinType(118)
	config.Seal()
}

func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func MakeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          authcodec.NewBech32Codec(GetBech32Prefix()),
			ValidatorAddressCodec: authcodec.NewBech32Codec(GetBech32Prefix() + "valoper"),
		},
	})
	if err != nil {
		panic(err)
	}
	appCodec := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(appCodec, authtx.DefaultSignModes)

	std.RegisterLegacyAminoCodec(amino)
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             appCodec,
		TxConfig:          txConfig,
		Amino:             amino,
	}
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	return paramsKeeper
}

// GetMaccPerms returns a copy of the module account permissions.
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// ModuleBasics defines the module basics
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	genutil.AppModuleBasic{},
	consensus.AppModuleBasic{},
	crisis.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	capability.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctransfer.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	pos.AppModuleBasic{},
	jamaah.AppModuleBasic{},
	paket.AppModuleBasic{},
	pembayaran.AppModuleBasic{},
	visa.AppModuleBasic{},
	hotel.AppModuleBasic{},
	ticket.AppModuleBasic{},
	referral.AppModuleBasic{},
	reward.AppModuleBasic{},
	audit.AppModuleBasic{},
	keberangkatan.AppModuleBasic{},
	agen.AppModuleBasic{},
	oleholeh.AppModuleBasic{},
	asuransi.AppModuleBasic{},
	dao.AppModuleBasic{},
)

var (
	_ servertypes.Application = &App{}
)
