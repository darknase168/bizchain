package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/asuransi/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Asuransi queries a policy by ID
func (k Keeper) Asuransi(c context.Context, req *types.QueryGetAsuransiRequest) (*types.QueryGetAsuransiResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	asuransi, found := k.GetAsuransi(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "asuransi not found")
	}

	return &types.QueryGetAsuransiResponse{Asuransi: &asuransi}, nil
}

// AsuransiAll queries all policies
func (k Keeper) AsuransiAll(c context.Context, req *types.QueryAllAsuransiRequest) (*types.QueryAllAsuransiResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	list := k.GetAllAsuransi(ctx)

	return &types.QueryAllAsuransiResponse{
		Asuransi: list,
	}, nil
}

// Claim queries a claim by ID
func (k Keeper) Claim(c context.Context, req *types.QueryGetClaimRequest) (*types.QueryGetClaimResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	claim, found := k.GetClaim(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "claim not found")
	}

	return &types.QueryGetClaimResponse{Claim: &claim}, nil
}

// ClaimAll queries all claims
func (k Keeper) ClaimAll(c context.Context, req *types.QueryAllClaimRequest) (*types.QueryAllClaimResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	list := k.GetAllClaims(ctx)

	return &types.QueryAllClaimResponse{
		Claim: list,
	}, nil
}
