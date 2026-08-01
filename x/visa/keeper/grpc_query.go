package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/visa/types"
)

var _ types.QueryServer = Keeper{}

// Visa handles the gRPC query for a single visa
func (k Keeper) Visa(goCtx context.Context, req *types.QueryGetVisaRequest) (*types.QueryGetVisaResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	visa, found := k.GetVisa(ctx, req.Id)
	if !found {
		return nil, types.ErrVisaNotFound
	}

	return &types.QueryGetVisaResponse{Visa: &visa}, nil
}

// VisaAll handles the gRPC query for all visas
func (k Keeper) VisaAll(goCtx context.Context, req *types.QueryAllVisaRequest) (*types.QueryAllVisaResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	visa := k.GetAllVisa(ctx)

	return &types.QueryAllVisaResponse{
		Visa: visa,
	}, nil
}
