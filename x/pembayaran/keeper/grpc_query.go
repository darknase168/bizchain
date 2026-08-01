package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/pembayaran/types"
)

var _ types.QueryServer = Keeper{}

// Pembayaran handles the gRPC query for a single pembayaran
func (k Keeper) Pembayaran(goCtx context.Context, req *types.QueryGetPembayaranRequest) (*types.QueryGetPembayaranResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pembayaran, found := k.GetPembayaran(ctx, req.Id)
	if !found {
		return nil, types.ErrPembayaranNotFound
	}

	return &types.QueryGetPembayaranResponse{Pembayaran: &pembayaran}, nil
}

// PembayaranAll handles the gRPC query for all pembayaran
func (k Keeper) PembayaranAll(goCtx context.Context, req *types.QueryAllPembayaranRequest) (*types.QueryAllPembayaranResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pembayaran := k.GetAllPembayaran(ctx)

	return &types.QueryAllPembayaranResponse{
		Pembayaran: pembayaran,
	}, nil
}

// PembayaranByJamaah handles the gRPC query for payments of a jamaah
func (k Keeper) PembayaranByJamaah(goCtx context.Context, req *types.QueryByJamaahRequest) (*types.QueryAllPembayaranResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pembayaran := k.GetPembayaranByJamaah(ctx, req.Jamaah)

	return &types.QueryAllPembayaranResponse{
		Pembayaran: pembayaran,
	}, nil
}
