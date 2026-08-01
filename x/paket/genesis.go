package paket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/paket/keeper"
	"github.com/bizchain/blockchain/x/paket/types"
)

// InitGenesis initializes the paket module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetPaketCount(ctx, genState.PaketCount)

	for _, paket := range genState.PaketList {
		k.SetPaket(ctx, *paket)
	}
}

// ExportGenesis exports the paket module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.PaketList = k.GetAllPaket(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.PaketList) > 0 {
		genesis.PaketCount = genesis.PaketList[len(genesis.PaketList)-1].Id
	}

	return genesis
}
