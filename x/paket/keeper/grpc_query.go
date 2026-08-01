package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/paket/types"
)

var _ types.QueryServer = Keeper{}

// Paket handles the gRPC query for a single paket
func (k Keeper) Paket(goCtx context.Context, req *types.QueryGetPaketRequest) (*types.QueryGetPaketResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	paket, found := k.GetPaket(ctx, req.Id)
	if !found {
		return nil, types.ErrPaketNotFound
	}

	return &types.QueryGetPaketResponse{Paket: &paket}, nil
}

// PaketAll handles the gRPC query for all paket (marketplace)
func (k Keeper) PaketAll(goCtx context.Context, req *types.QueryAllPaketRequest) (*types.QueryAllPaketResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	paket := k.GetAllPaket(ctx)

	return &types.QueryAllPaketResponse{
		Paket: paket,
	}, nil
}
