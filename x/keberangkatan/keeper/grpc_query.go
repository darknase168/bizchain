package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/keberangkatan/types"
)

var _ types.QueryServer = Keeper{}

// Keberangkatan handles the gRPC query for a single journey
func (k Keeper) Keberangkatan(goCtx context.Context, req *types.QueryGetKeberangkatanRequest) (*types.QueryGetKeberangkatanResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	keberangkatan, found := k.GetKeberangkatan(ctx, req.Id)
	if !found {
		return nil, types.ErrKeberangkatanNotFound
	}

	return &types.QueryGetKeberangkatanResponse{Keberangkatan: &keberangkatan}, nil
}

// KeberangkatanAll handles the gRPC query for all journeys
func (k Keeper) KeberangkatanAll(goCtx context.Context, req *types.QueryAllKeberangkatanRequest) (*types.QueryAllKeberangkatanResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	keberangkatan := k.GetAllKeberangkatan(ctx)

	return &types.QueryAllKeberangkatanResponse{
		Keberangkatan: keberangkatan,
	}, nil
}

// KeberangkatanByJamaah handles the gRPC query for journeys of a pilgrim
func (k Keeper) KeberangkatanByJamaah(goCtx context.Context, req *types.QueryKeberangkatanByJamaahRequest) (*types.QueryKeberangkatanByJamaahResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	keberangkatan := k.GetKeberangkatanByJamaah(ctx, req.Jamaah)

	return &types.QueryKeberangkatanByJamaahResponse{
		Keberangkatan: keberangkatan,
	}, nil
}
