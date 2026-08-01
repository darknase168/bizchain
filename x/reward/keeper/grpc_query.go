package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/reward/types"
)

var _ types.QueryServer = Keeper{}

// Balance handles the gRPC query for a jamaah's loyalty balance
func (k Keeper) Balance(goCtx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	balance, found := k.GetRewardBalance(ctx, req.Jamaah)
	if !found {
		balance = types.RewardBalance{Jamaah: req.Jamaah, Balance: "0", Earned: "0", Redeemed: "0"}
	}

	return &types.QueryBalanceResponse{Balance: &balance}, nil
}

// Reward handles the gRPC query for a single reward
func (k Keeper) Reward(goCtx context.Context, req *types.QueryGetRewardRequest) (*types.QueryGetRewardResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	reward, found := k.GetReward(ctx, req.Id)
	if !found {
		return nil, types.ErrRewardNotFound
	}

	return &types.QueryGetRewardResponse{Reward: &reward}, nil
}

// RewardAll handles the gRPC query for all rewards
func (k Keeper) RewardAll(goCtx context.Context, req *types.QueryAllRewardRequest) (*types.QueryAllRewardResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	reward := k.GetAllReward(ctx)

	return &types.QueryAllRewardResponse{
		Reward: reward,
	}, nil
}
