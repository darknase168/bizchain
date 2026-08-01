package reward

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/reward/keeper"
	"github.com/bizchain/blockchain/x/reward/types"
)

// InitGenesis initializes the reward module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetRewardCount(ctx, genState.RewardCount)

	for _, reward := range genState.RewardList {
		k.SetReward(ctx, *reward)
	}
	for _, balance := range genState.BalanceList {
		k.SetRewardBalance(ctx, *balance)
	}
}

// ExportGenesis exports the reward module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.RewardList = k.GetAllReward(ctx)
	genesis.BalanceList = k.GetAllRewardBalance(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.RewardList) > 0 {
		genesis.RewardCount = genesis.RewardList[len(genesis.RewardList)-1].Id
	}

	return genesis
}
