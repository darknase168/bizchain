package oleholeh

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/oleholeh/keeper"
	"github.com/bizchain/blockchain/x/oleholeh/types"
)

// InitGenesis initializes the oleholeh module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetProductCount(ctx, genState.ProductCount)
	k.SetOrderCount(ctx, genState.OrderCount)

	for _, product := range genState.ProductList {
		k.SetProduct(ctx, *product)
	}
	for _, order := range genState.OrderList {
		k.SetOrder(ctx, *order)
	}
}

// ExportGenesis exports the oleholeh module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.ProductList = k.GetAllProducts(ctx)
	genesis.OrderList = k.GetAllOrders(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.ProductList) > 0 {
		genesis.ProductCount = genesis.ProductList[len(genesis.ProductList)-1].Id
	}
	if len(genesis.OrderList) > 0 {
		genesis.OrderCount = genesis.OrderList[len(genesis.OrderList)-1].Id
	}

	return genesis
}
