package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/hotel/types"
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

// CreateHotel registers a new hotel
func (k msgServer) CreateHotel(goCtx context.Context, msg *types.MsgCreateHotel) (*types.MsgCreateHotelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hotelID := k.GetNextHotelID(ctx)
	hotel := types.Hotel{
		Id:            hotelID,
		Name:          msg.Name,
		City:          msg.City,
		Address:       msg.Address,
		StarRating:    msg.StarRating,
		PricePerNight: msg.PricePerNight,
		RoomType:      msg.RoomType,
		AvailableRooms: msg.AvailableRooms,
		DistanceHaram: msg.DistanceHaram,
		Status:        "active",
		Creator:       msg.Creator,
		CreatedAt:     nowStr(ctx),
		UpdatedAt:     nowStr(ctx),
	}
	k.SetHotel(ctx, hotel)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateHotel,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyHotelID, fmt.Sprintf("%d", hotelID)),
			sdk.NewAttribute(types.AttributeKeyHotelName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateHotelResponse{Id: hotelID}, nil
}

// UpdateHotel updates hotel data
func (k msgServer) UpdateHotel(goCtx context.Context, msg *types.MsgUpdateHotel) (*types.MsgUpdateHotelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hotel, found := k.GetHotel(ctx, msg.Id)
	if !found {
		return nil, types.ErrHotelNotFound
	}

	hotel.Name = msg.Name
	hotel.City = msg.City
	hotel.Address = msg.Address
	hotel.StarRating = msg.StarRating
	hotel.PricePerNight = msg.PricePerNight
	hotel.RoomType = msg.RoomType
	hotel.AvailableRooms = msg.AvailableRooms
	hotel.DistanceHaram = msg.DistanceHaram
	if msg.Status != "" {
		hotel.Status = msg.Status
	}
	hotel.UpdatedAt = nowStr(ctx)
	k.SetHotel(ctx, hotel)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateHotel,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyHotelID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyStatus, hotel.Status),
		),
	)

	return &types.MsgUpdateHotelResponse{}, nil
}
