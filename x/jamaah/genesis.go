package jamaah

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/jamaah/keeper"
	"github.com/bizchain/blockchain/x/jamaah/types"
)

// InitGenesis initializes the jamaah module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetJamaahCount(ctx, genState.JamaahCount)

	for _, jamaah := range genState.JamaahList {
		k.SetJamaah(ctx, *jamaah)
	}
}

// ExportGenesis exports the jamaah module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.JamaahList = k.GetAllJamaah(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.JamaahList) > 0 {
		genesis.JamaahCount = genesis.JamaahList[len(genesis.JamaahList)-1].Id
	}

	return genesis
}
