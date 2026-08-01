package hotel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/hotel/keeper"
	"github.com/bizchain/blockchain/x/hotel/types"
)

// InitGenesis initializes the hotel module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetHotelCount(ctx, genState.HotelCount)

	for _, hotel := range genState.HotelList {
		k.SetHotel(ctx, *hotel)
	}
}

// ExportGenesis exports the hotel module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.HotelList = k.GetAllHotel(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.HotelList) > 0 {
		genesis.HotelCount = genesis.HotelList[len(genesis.HotelList)-1].Id
	}

	return genesis
}
