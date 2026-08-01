package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/agen/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Agen queries an agen by ID
func (k Keeper) Agen(c context.Context, req *types.QueryGetAgenRequest) (*types.QueryGetAgenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	agen, found := k.GetAgen(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "agen not found")
	}

	return &types.QueryGetAgenResponse{Agen: &agen}, nil
}

// AgenAll queries all agents
func (k Keeper) AgenAll(c context.Context, req *types.QueryAllAgenRequest) (*types.QueryAllAgenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	list := k.GetAllAgen(ctx)

	return &types.QueryAllAgenResponse{
		Agen: list,
	}, nil
}

// Hierarki queries the downline tree of an agent
func (k Keeper) Hierarki(c context.Context, req *types.QueryHierarkiRequest) (*types.QueryHierarkiResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	downline := k.GetDownline(ctx, req.Id)

	return &types.QueryHierarkiResponse{Downline: downline}, nil
}
