package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/hotel/types"
)

var _ types.QueryServer = Keeper{}

// Hotel handles the gRPC query for a single hotel
func (k Keeper) Hotel(goCtx context.Context, req *types.QueryGetHotelRequest) (*types.QueryGetHotelResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hotel, found := k.GetHotel(ctx, req.Id)
	if !found {
		return nil, types.ErrHotelNotFound
	}

	return &types.QueryGetHotelResponse{Hotel: &hotel}, nil
}

// HotelAll handles the gRPC query for all hotels
func (k Keeper) HotelAll(goCtx context.Context, req *types.QueryAllHotelRequest) (*types.QueryAllHotelResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hotel := k.GetAllHotel(ctx)

	return &types.QueryAllHotelResponse{
		Hotel: hotel,
	}, nil
}
