package referral

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/referral/keeper"
	"github.com/bizchain/blockchain/x/referral/types"
)

// InitGenesis initializes the referral module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetReferralCount(ctx, genState.ReferralCount)

	for _, referral := range genState.ReferralList {
		k.SetReferral(ctx, *referral)
	}
}

// ExportGenesis exports the referral module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.ReferralList = k.GetAllReferral(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.ReferralList) > 0 {
		genesis.ReferralCount = genesis.ReferralList[len(genesis.ReferralList)-1].Id
	}

	return genesis
}
