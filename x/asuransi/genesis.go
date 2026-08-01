package asuransi

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/asuransi/keeper"
	"github.com/bizchain/blockchain/x/asuransi/types"
)

// InitGenesis initializes the asuransi module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetAsuransiCount(ctx, genState.AsuransiCount)
	k.SetClaimCount(ctx, genState.ClaimCount)

	for _, asuransi := range genState.AsuransiList {
		k.SetAsuransi(ctx, *asuransi)
	}
	for _, claim := range genState.ClaimList {
		k.SetClaim(ctx, *claim)
	}
}

// ExportGenesis exports the asuransi module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.AsuransiList = k.GetAllAsuransi(ctx)
	genesis.ClaimList = k.GetAllClaims(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.AsuransiList) > 0 {
		genesis.AsuransiCount = genesis.AsuransiList[len(genesis.AsuransiList)-1].Id
	}
	if len(genesis.ClaimList) > 0 {
		genesis.ClaimCount = genesis.ClaimList[len(genesis.ClaimList)-1].Id
	}

	return genesis
}
