package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/referral/types"
)

var _ types.QueryServer = Keeper{}

// Referral handles the gRPC query for a single referral
func (k Keeper) Referral(goCtx context.Context, req *types.QueryGetReferralRequest) (*types.QueryGetReferralResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	referral, found := k.GetReferral(ctx, req.Id)
	if !found {
		return nil, types.ErrReferralNotFound
	}

	return &types.QueryGetReferralResponse{Referral: &referral}, nil
}

// ReferralAll handles the gRPC query for all referrals
func (k Keeper) ReferralAll(goCtx context.Context, req *types.QueryAllReferralRequest) (*types.QueryAllReferralResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	referral := k.GetAllReferral(ctx)

	return &types.QueryAllReferralResponse{
		Referral: referral,
	}, nil
}

// ReferralByAgent handles the gRPC query for referrals of an agent
func (k Keeper) ReferralByAgent(goCtx context.Context, req *types.QueryByAgentRequest) (*types.QueryAllReferralResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	referral := k.GetReferralByAgent(ctx, req.Agent)

	return &types.QueryAllReferralResponse{
		Referral: referral,
	}, nil
}
