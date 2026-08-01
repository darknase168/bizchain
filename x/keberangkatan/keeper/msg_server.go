package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/keberangkatan/types"
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

// CreateKeberangkatan starts tracking a pilgrim's departure journey (stage 1)
func (k msgServer) CreateKeberangkatan(goCtx context.Context, msg *types.MsgCreateKeberangkatan) (*types.MsgCreateKeberangkatanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keberangkatanID := k.GetNextKeberangkatanID(ctx)
	keberangkatan := types.Keberangkatan{
		Id:            keberangkatanID,
		Jamaah:        msg.Jamaah,
		PaketId:       msg.PaketId,
		PembayaranId:  msg.PembayaranId,
		Stage:         1,
		StatusLabel:   types.StageLabels[1],
		DepartureDate: msg.DepartureDate,
		ReturnDate:    msg.ReturnDate,
		ManasikDate:   msg.ManasikDate,
		Creator:       msg.Creator,
		CreatedAt:     nowStr(ctx),
		UpdatedAt:     nowStr(ctx),
	}
	k.SetKeberangkatan(ctx, keberangkatan)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateKeberangkatan,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyKeberangkatanID, fmt.Sprintf("%d", keberangkatanID)),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.Jamaah),
			sdk.NewAttribute(types.AttributeKeyStage, "1"),
			sdk.NewAttribute(types.AttributeKeyStatusLabel, types.StageLabels[1]),
		),
	)

	return &types.MsgCreateKeberangkatanResponse{Id: keberangkatanID}, nil
}

// AdvanceStage advances a journey to the next stage (1-9)
func (k msgServer) AdvanceStage(goCtx context.Context, msg *types.MsgAdvanceStage) (*types.MsgAdvanceStageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keberangkatan, found := k.GetKeberangkatan(ctx, msg.Id)
	if !found {
		return nil, types.ErrKeberangkatanNotFound
	}
	if keberangkatan.Stage >= 9 {
		return nil, types.ErrJourneyCompleted
	}

	keberangkatan.Stage++
	keberangkatan.StatusLabel = types.StageLabels[keberangkatan.Stage]
	keberangkatan.UpdatedAt = nowStr(ctx)
	k.SetKeberangkatan(ctx, keberangkatan)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAdvanceStage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyKeberangkatanID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyStage, fmt.Sprintf("%d", keberangkatan.Stage)),
			sdk.NewAttribute(types.AttributeKeyStatusLabel, keberangkatan.StatusLabel),
		),
	)

	return &types.MsgAdvanceStageResponse{Stage: keberangkatan.Stage, StatusLabel: keberangkatan.StatusLabel}, nil
}

// UpdateDeparture updates departure/return/manasik dates
func (k msgServer) UpdateDeparture(goCtx context.Context, msg *types.MsgUpdateDeparture) (*types.MsgUpdateDepartureResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keberangkatan, found := k.GetKeberangkatan(ctx, msg.Id)
	if !found {
		return nil, types.ErrKeberangkatanNotFound
	}

	if msg.DepartureDate != "" {
		keberangkatan.DepartureDate = msg.DepartureDate
	}
	if msg.ReturnDate != "" {
		keberangkatan.ReturnDate = msg.ReturnDate
	}
	if msg.ManasikDate != "" {
		keberangkatan.ManasikDate = msg.ManasikDate
	}
	if msg.HotelConfirm != "" {
		keberangkatan.HotelConfirm = msg.HotelConfirm
	}
	if msg.AirlineConfirm != "" {
		keberangkatan.AirlineConfirm = msg.AirlineConfirm
	}
	keberangkatan.UpdatedAt = nowStr(ctx)
	k.SetKeberangkatan(ctx, keberangkatan)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateDeparture,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyKeberangkatanID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgUpdateDepartureResponse{}, nil
}

// AddBaggage adds a checked baggage item (QR/NFC tracking)
func (k msgServer) AddBaggage(goCtx context.Context, msg *types.MsgAddBaggage) (*types.MsgAddBaggageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keberangkatan, found := k.GetKeberangkatan(ctx, msg.Id)
	if !found {
		return nil, types.ErrKeberangkatanNotFound
	}

	baggageID := k.GetNextBaggageID(ctx)
	keberangkatan.Baggage = append(keberangkatan.Baggage, &types.BaggageItem{
		Id:        baggageID,
		Tag:       msg.Tag,
		Weight:    msg.Weight,
		Status:    "checked_in",
		UpdatedAt: nowStr(ctx),
	})
	keberangkatan.UpdatedAt = nowStr(ctx)
	k.SetKeberangkatan(ctx, keberangkatan)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddBaggage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyKeberangkatanID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyBaggageID, fmt.Sprintf("%d", baggageID)),
		),
	)

	return &types.MsgAddBaggageResponse{BaggageId: baggageID}, nil
}

// UpdateBaggageStatus updates baggage tracking status
func (k msgServer) UpdateBaggageStatus(goCtx context.Context, msg *types.MsgUpdateBaggageStatus) (*types.MsgUpdateBaggageStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keberangkatan, found := k.GetKeberangkatan(ctx, msg.Id)
	if !found {
		return nil, types.ErrKeberangkatanNotFound
	}

	switch msg.Status {
	case "checked_in", "in_transit", "arrived", "delivered":
	default:
		return nil, types.ErrInvalidBaggageStatus
	}

	for _, baggage := range keberangkatan.Baggage {
		if baggage.Id == msg.BaggageId {
			baggage.Status = msg.Status
			baggage.UpdatedAt = nowStr(ctx)
			keberangkatan.UpdatedAt = nowStr(ctx)
			k.SetKeberangkatan(ctx, keberangkatan)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeUpdateBaggageStatus,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
					sdk.NewAttribute(types.AttributeKeyKeberangkatanID, fmt.Sprintf("%d", msg.Id)),
					sdk.NewAttribute(types.AttributeKeyBaggageID, fmt.Sprintf("%d", msg.BaggageId)),
				),
			)

			return &types.MsgUpdateBaggageStatusResponse{}, nil
		}
	}

	return nil, types.ErrBaggageNotFound
}
