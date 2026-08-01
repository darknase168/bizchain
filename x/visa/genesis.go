package visa

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/visa/keeper"
	"github.com/bizchain/blockchain/x/visa/types"
)

// InitGenesis initializes the visa module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetVisaCount(ctx, genState.VisaCount)

	for _, visa := range genState.VisaList {
		k.SetVisa(ctx, *visa)
	}
}

// ExportGenesis exports the visa module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.VisaList = k.GetAllVisa(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.VisaList) > 0 {
		genesis.VisaCount = genesis.VisaList[len(genesis.VisaList)-1].Id
	}

	return genesis
}
