package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/jamaah/types"
)

var _ types.QueryServer = Keeper{}

// Jamaah handles the gRPC query for a single jamaah
func (k Keeper) Jamaah(goCtx context.Context, req *types.QueryGetJamaahRequest) (*types.QueryGetJamaahResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	jamaah, found := k.GetJamaah(ctx, req.Id)
	if !found {
		return nil, types.ErrJamaahNotFound
	}

	return &types.QueryGetJamaahResponse{Jamaah: &jamaah}, nil
}

// JamaahAll handles the gRPC query for all jamaah
func (k Keeper) JamaahAll(goCtx context.Context, req *types.QueryAllJamaahRequest) (*types.QueryAllJamaahResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	jamaah := k.GetAllJamaah(ctx)

	return &types.QueryAllJamaahResponse{
		Jamaah: jamaah,
	}, nil
}
