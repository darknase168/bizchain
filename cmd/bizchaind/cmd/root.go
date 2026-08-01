package cmd

import (
	"io"
	"os"

	"cosmossdk.io/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/app"
)

// NewRootCmd creates the root command for bizchaind
func NewRootCmd() (*cobra.Command, app.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()

	// Pre-instantiate the application to obtain a codec-carrying module
	// BasicManager for wiring module CLI commands (v0.50 simapp pattern).
	tempApp := app.NewApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		simtestutil.NewAppOptionsWithFlagHome(app.GetDefaultNodeHome()),
	)
	basicManager := tempApp.BasicModuleManager

	rootCmd := &cobra.Command{
		Use:   "bizchaind",
		Short: "BizChain - A feeless retail/minimarket & education blockchain",
		Long: `BizChain is a custom Cosmos SDK blockchain designed for 
retail/minimarket and education businesses with ZERO transaction fees.

Commands:
  - tx pos create-product        Create a new product
  - tx pos create-transaction    Create a POS sale transaction
  - tx pos add-stock             Add stock to a product
  - query pos list-product       List all products
  - query pos show-product       Show a product by ID
  - query bank balances          Query account balances
  - keys                         Manage wallet keys
  - tx agen create-agen          Register an agent (pusat/cabang/subagen)
  - tx agen record-performance   Record agent performance (rekam jejak)
  - tx oleholeh create           List a souvenir product
  - tx asuransi create           Issue a digital insurance policy
  - tx dao create-proposal       Create a DAO proposal
  - query agen list              List agents
  - query dao list               List DAO proposals
  - tx ibc-transfer transfer     Transfer tokens across IBC
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx := client.Context{}.
				WithCodec(encodingConfig.Codec).
				WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
				WithTxConfig(encodingConfig.TxConfig).
				WithLegacyAmino(encodingConfig.Amino).
				WithInput(os.Stdin).
				WithAccountRetriever(types.AccountRetriever{}).
				WithHomeDir(app.GetDefaultNodeHome()).
				WithViper("")

			clientCtx, err := config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(clientCtx, cmd); err != nil {
				return err
			}

			// Create the home/config directories and load the server config
			// (app.toml, config.toml, min gas prices). This must run for every
			// server command (init, start, etc.) like the standard simapp.
			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := cmtcfg.DefaultConfig()
			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	// Note: the --home and --trace persistent flags are added by svrcmd.Execute
	// via cometbft/libs/cli.PrepareBaseCmd, so they must NOT be registered here
	// (doing so panics with 'flag redefined: home'/'trace').

	rootCmd.AddCommand(
		server.StatusCommand(),
		rpc.ValidatorCommand(),
		queryCommand(basicManager),
		txCommand(basicManager),
		keys.Commands(),
		debug.Cmd(),
		genutilcli.InitCmd(basicManager, app.GetDefaultNodeHome()),
		genutilcli.CollectGenTxsCmd(
			banktypes.GenesisBalancesIterator{},
			app.GetDefaultNodeHome(),
			genutiltypes.DefaultMessageValidator,
			authcodec.NewBech32Codec(app.GetBech32Prefix()+"valoper"),
		),
		genutilcli.GenTxCmd(
			basicManager,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.GetDefaultNodeHome(),
			authcodec.NewBech32Codec(app.GetBech32Prefix()+"valoper"),
		),
		genutilcli.ValidateGenesisCmd(basicManager),
		genutilcli.AddGenesisAccountCmd(
			app.GetDefaultNodeHome(),
			authcodec.NewBech32Codec(app.GetBech32Prefix()),
		),
	)

	// Note: server.AddCommands already registers start, comet, export, version,
	// rollback and module-hash commands on the root command.
	server.AddCommands(rootCmd, app.GetDefaultNodeHome(), newApp, newAppExport, func(_ *cobra.Command) {})

	return rootCmd, encodingConfig
}

func queryCommand(basicManager module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(rpc.QueryEventForTxCmd())
	cmd.AddCommand(
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)
	// Wire query commands from every registered module (pos, agen, oleholeh,
	// asuransi, dao, and the haji/umroh ecosystem modules).
	basicManager.AddQueryCommands(cmd)
	return cmd
}

func txCommand(basicManager module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Aliases:                    []string{"t"},
		Short:                      "Transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(rpc.WaitTxCmd())
	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)
	// Wire transaction commands from every registered module.
	basicManager.AddTxCommands(cmd)
	return cmd
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	return app.NewApp(
		logger,
		db,
		traceStore,
		true,
		appOpts,
	)
}

func newAppExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	bapp := app.NewApp(logger, db, traceStore, true, appOpts)
	if height != -1 {
		if err := bapp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}
	return bapp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

func initAppConfig() (string, interface{}) {
	customAppConfig := serverconfig.DefaultConfig()
	customAppConfig.MinGasPrices = "0upoint"
	customAppConfig.API.Enable = true
	customAppConfig.API.EnableUnsafeCORS = true
	customAppConfig.GRPC.Enable = true
	customAppConfig.GRPCWeb.Enable = true

	template := serverconfig.DefaultConfigTemplate
	return template, customAppConfig
}
