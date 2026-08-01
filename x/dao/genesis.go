package dao

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/dao/keeper"
	"github.com/bizchain/blockchain/x/dao/types"
)

// InitGenesis initializes the dao module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetProposalCount(ctx, genState.ProposalCount)

	for _, proposal := range genState.ProposalList {
		k.SetProposal(ctx, *proposal)
	}
}

// ExportGenesis exports the dao module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.ProposalList = k.GetAllProposals(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.ProposalList) > 0 {
		genesis.ProposalCount = genesis.ProposalList[len(genesis.ProposalList)-1].Id
	}

	return genesis
}
