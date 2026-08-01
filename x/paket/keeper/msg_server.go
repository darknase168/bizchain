package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/paket/types"
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

// CreatePaket creates a new package (smart contract paket)
func (k msgServer) CreatePaket(goCtx context.Context, msg *types.MsgCreatePaket) (*types.MsgCreatePaketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Quota == 0 {
		return nil, types.ErrInvalidQuota
	}

	paketID := k.GetNextPaketID(ctx)
	paket := types.Paket{
		Id:            paketID,
		Name:          msg.Name,
		Price:         msg.Price,
		Schedule:      msg.Schedule,
		Quota:         msg.Quota,
		Booked:        0,
		Hotel:         msg.Hotel,
		Airline:       msg.Airline,
		Muthawif:      msg.Muthawif,
		Status:        "open",
		DepartureDate: msg.DepartureDate,
		ReturnDate:    msg.ReturnDate,
		Category:      msg.Category,
		Description:   msg.Description,
		ImageUrl:      msg.ImageUrl,
		Creator:       msg.Creator,
		CreatedAt:     nowStr(ctx),
		UpdatedAt:     nowStr(ctx),
	}
	k.SetPaket(ctx, paket)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreatePaket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPaketID, fmt.Sprintf("%d", paketID)),
			sdk.NewAttribute(types.AttributeKeyPaketName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreatePaketResponse{Id: paketID}, nil
}

// UpdatePaket updates package data
func (k msgServer) UpdatePaket(goCtx context.Context, msg *types.MsgUpdatePaket) (*types.MsgUpdatePaketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	paket, found := k.GetPaket(ctx, msg.Id)
	if !found {
		return nil, types.ErrPaketNotFound
	}
	if paket.Creator != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	paket.Name = msg.Name
	paket.Price = msg.Price
	paket.Schedule = msg.Schedule
	paket.Quota = msg.Quota
	paket.Hotel = msg.Hotel
	paket.Airline = msg.Airline
	paket.Muthawif = msg.Muthawif
	paket.DepartureDate = msg.DepartureDate
	paket.ReturnDate = msg.ReturnDate
	paket.Category = msg.Category
	paket.Description = msg.Description
	paket.ImageUrl = msg.ImageUrl
	if msg.Status != "" {
		paket.Status = msg.Status
	}
	paket.UpdatedAt = nowStr(ctx)
	k.SetPaket(ctx, paket)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdatePaket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPaketID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgUpdatePaketResponse{}, nil
}

// BookPaket books a package and auto-closes quota when full
func (k msgServer) BookPaket(goCtx context.Context, msg *types.MsgBookPaket) (*types.MsgBookPaketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	paket, found := k.GetPaket(ctx, msg.PaketId)
	if !found {
		return nil, types.ErrPaketNotFound
	}
	if paket.Status == "full" || paket.Status == "closed" || paket.Status == "departed" || paket.Status == "completed" {
		return nil, types.ErrQuotaFull
	}
	if msg.Quantity == 0 {
		return nil, types.ErrInvalidQuota
	}

	newBooked := paket.Booked + msg.Quantity
	quotaClosed := false
	if newBooked >= paket.Quota {
		paket.Status = "full"
		quotaClosed = true
	}
	paket.Booked = newBooked
	paket.UpdatedAt = nowStr(ctx)
	k.SetPaket(ctx, paket)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBookPaket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPaketID, fmt.Sprintf("%d", msg.PaketId)),
			sdk.NewAttribute(types.AttributeKeyQuotaClosed, fmt.Sprintf("%t", quotaClosed)),
		),
	)

	return &types.MsgBookPaketResponse{QuotaClosed: quotaClosed}, nil
}

// AddReview adds a rating/review to a package (marketplace)
func (k msgServer) AddReview(goCtx context.Context, msg *types.MsgAddReview) (*types.MsgAddReviewResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	paket, found := k.GetPaket(ctx, msg.PaketId)
	if !found {
		return nil, types.ErrPaketNotFound
	}
	if msg.Rating < 1 || msg.Rating > 5 {
		return nil, types.ErrInvalidRating
	}

	paket.Reviews = append(paket.Reviews, &types.Review{
		Reviewer:  msg.Creator,
		Rating:    msg.Rating,
		Comment:   msg.Comment,
		CreatedAt: nowStr(ctx),
	})
	paket.UpdatedAt = nowStr(ctx)
	k.SetPaket(ctx, paket)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddReview,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPaketID, fmt.Sprintf("%d", msg.PaketId)),
			sdk.NewAttribute(types.AttributeKeyReviewer, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyRating, fmt.Sprintf("%d", msg.Rating)),
		),
	)

	return &types.MsgAddReviewResponse{}, nil
}
