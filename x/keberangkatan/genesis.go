package keberangkatan

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/keberangkatan/keeper"
	"github.com/bizchain/blockchain/x/keberangkatan/types"
)

// InitGenesis initializes the keberangkatan module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetKeberangkatanCount(ctx, genState.KeberangkatanCount)
	k.SetBaggageCount(ctx, genState.BaggageCount)

	for _, keberangkatan := range genState.KeberangkatanList {
		k.SetKeberangkatan(ctx, *keberangkatan)
	}
}

// ExportGenesis exports the keberangkatan module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.KeberangkatanList = k.GetAllKeberangkatan(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.KeberangkatanList) > 0 {
		genesis.KeberangkatanCount = genesis.KeberangkatanList[len(genesis.KeberangkatanList)-1].Id
	}
	for _, keberangkatan := range genesis.KeberangkatanList {
		for _, baggage := range keberangkatan.Baggage {
			if baggage.Id > genesis.BaggageCount {
				genesis.BaggageCount = baggage.Id
			}
		}
	}

	return genesis
}
