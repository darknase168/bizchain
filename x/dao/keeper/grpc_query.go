package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/dao/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Proposal queries a proposal by ID
func (k Keeper) Proposal(c context.Context, req *types.QueryGetProposalRequest) (*types.QueryGetProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	proposal, found := k.GetProposal(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "proposal not found")
	}

	return &types.QueryGetProposalResponse{Proposal: &proposal}, nil
}

// ProposalAll queries all proposals
func (k Keeper) ProposalAll(c context.Context, req *types.QueryAllProposalRequest) (*types.QueryAllProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	list := k.GetAllProposals(ctx)

	return &types.QueryAllProposalResponse{
		Proposal: list,
	}, nil
}
