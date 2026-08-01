package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/ticket/types"
)

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func nowStr(ctx sdk.Context) string {
	return ctx.BlockTime().UTC().Format("2006-01-02T15:04:05Z")
}

// IssueTicket issues an NFT ticket for a jamaah
func (k msgServer) IssueTicket(goCtx context.Context, msg *types.MsgIssueTicket) (*types.MsgIssueTicketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ticketID := k.GetNextTicketID(ctx)
	nftID := k.BuildNFTID(ctx, ticketID)
	ticket := types.Ticket{
		Id:           ticketID,
		Jamaah:       msg.Jamaah,
		PaketId:      msg.PaketId,
		Airline:      msg.Airline,
		FlightNumber: msg.FlightNumber,
		Seat:         msg.Seat,
		Schedule:     msg.Schedule,
		QrCode:       msg.QrCode,
		NftId:        nftID,
		Status:       "issued",
		DocumentHash: msg.DocumentHash,
		Creator:      msg.Creator,
		CreatedAt:    nowStr(ctx),
	}
	k.SetTicket(ctx, ticket)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIssueTicket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTicketID, fmt.Sprintf("%d", ticketID)),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.Jamaah),
			sdk.NewAttribute(types.AttributeKeyNFTID, nftID),
			sdk.NewAttribute(types.AttributeKeyStatus, "issued"),
		),
	)

	return &types.MsgIssueTicketResponse{Id: ticketID, NftId: nftID}, nil
}

// UpdateTicketStatus updates ticket status (checked_in, boarded, etc)
func (k msgServer) UpdateTicketStatus(goCtx context.Context, msg *types.MsgUpdateTicketStatus) (*types.MsgUpdateTicketStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ticket, found := k.GetTicket(ctx, msg.Id)
	if !found {
		return nil, types.ErrTicketNotFound
	}

	switch msg.Status {
	case "issued", "checked_in", "boarded", "used", "void":
	default:
		return nil, types.ErrInvalidStatus
	}

	ticket.Status = msg.Status
	k.SetTicket(ctx, ticket)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateTicketStatus,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTicketID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.Status),
		),
	)

	return &types.MsgUpdateTicketStatusResponse{}, nil
}
