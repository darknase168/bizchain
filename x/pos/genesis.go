package pos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pos/keeper"
	"github.com/bizchain/blockchain/x/pos/types"
)

// InitGenesis initializes the POS module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all counters
	k.SetProductCount(ctx, genState.ProductCount)
	k.SetTransactionCount(ctx, genState.TransactionCount)
	k.SetUnitCount(ctx, genState.UnitCount)
	k.SetAccountCount(ctx, genState.AccountCount)
	k.SetJournalCount(ctx, genState.JournalCount)

	// Set all units from genesis
	for _, unit := range genState.UnitList {
		k.SetUnit(ctx, *unit)
	}

	// Set all accounts from genesis
	for _, account := range genState.AccountList {
		k.SetAccount(ctx, *account)
	}

	// Set all products from genesis
	for _, product := range genState.ProductList {
		k.SetProduct(ctx, *product)
	}

	// Set all transactions from genesis
	for _, transaction := range genState.TransactionList {
		k.SetTransaction(ctx, *transaction)
	}

	// Set all journal entries from genesis
	for _, entry := range genState.JournalList {
		k.SetJournalEntry(ctx, *entry)
	}
}

// ExportGenesis exports the POS module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.ProductList = k.GetAllProducts(ctx)
	genesis.TransactionList = k.GetAllTransactions(ctx)
	genesis.UnitList = k.GetAllUnits(ctx)
	genesis.AccountList = k.GetAllAccounts(ctx)
	genesis.JournalList = k.GetAllJournalEntries(ctx)
	genesis.Params = types.DefaultParams()

	// Update counts
	if len(genesis.ProductList) > 0 {
		genesis.ProductCount = genesis.ProductList[len(genesis.ProductList)-1].Id
	}
	if len(genesis.TransactionList) > 0 {
		genesis.TransactionCount = genesis.TransactionList[len(genesis.TransactionList)-1].Id
	}
	if len(genesis.UnitList) > 0 {
		genesis.UnitCount = genesis.UnitList[len(genesis.UnitList)-1].Id
	}
	if len(genesis.AccountList) > 0 {
		genesis.AccountCount = genesis.AccountList[len(genesis.AccountList)-1].Id
	}
	if len(genesis.JournalList) > 0 {
		genesis.JournalCount = genesis.JournalList[len(genesis.JournalList)-1].Id
	}

	return genesis
}
