package agen

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/agen/keeper"
	"github.com/bizchain/blockchain/x/agen/types"
)

// InitGenesis initializes the agen module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetAgenCount(ctx, genState.AgenCount)
	k.SetComplaintCount(ctx, genState.ComplaintCount)

	for _, agen := range genState.AgenList {
		k.SetAgen(ctx, *agen)
	}
}

// ExportGenesis exports the agen module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.AgenList = k.GetAllAgen(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.AgenList) > 0 {
		genesis.AgenCount = genesis.AgenList[len(genesis.AgenList)-1].Id
	}

	return genesis
}
