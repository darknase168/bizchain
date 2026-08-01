package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/ticket/types"
)

var _ types.QueryServer = Keeper{}

// Ticket handles the gRPC query for a single ticket
func (k Keeper) Ticket(goCtx context.Context, req *types.QueryGetTicketRequest) (*types.QueryGetTicketResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	ticket, found := k.GetTicket(ctx, req.Id)
	if !found {
		return nil, types.ErrTicketNotFound
	}

	return &types.QueryGetTicketResponse{Ticket: &ticket}, nil
}

// TicketAll handles the gRPC query for all tickets
func (k Keeper) TicketAll(goCtx context.Context, req *types.QueryAllTicketRequest) (*types.QueryAllTicketResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	ticket := k.GetAllTicket(ctx)

	return &types.QueryAllTicketResponse{
		Ticket: ticket,
	}, nil
}
