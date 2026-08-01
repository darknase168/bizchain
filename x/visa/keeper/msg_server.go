package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/visa/types"
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

// CreateVisa starts a visa application
func (k msgServer) CreateVisa(goCtx context.Context, msg *types.MsgCreateVisa) (*types.MsgCreateVisaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	visaID := k.GetNextVisaID(ctx)
	visa := types.Visa{
		Id:           visaID,
		Jamaah:       msg.Jamaah,
		PaketId:      msg.PaketId,
		Status:       "processing",
		DocumentHash: msg.DocumentHash,
		Notes:        msg.Notes,
		Creator:      msg.Creator,
		CreatedAt:    nowStr(ctx),
		UpdatedAt:    nowStr(ctx),
	}
	k.SetVisa(ctx, visa)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateVisa,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyVisaID, fmt.Sprintf("%d", visaID)),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.Jamaah),
			sdk.NewAttribute(types.AttributeKeyStatus, "processing"),
		),
	)

	return &types.MsgCreateVisaResponse{Id: visaID}, nil
}

// UpdateVisaStatus updates visa processing status
func (k msgServer) UpdateVisaStatus(goCtx context.Context, msg *types.MsgUpdateVisaStatus) (*types.MsgUpdateVisaStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	visa, found := k.GetVisa(ctx, msg.Id)
	if !found {
		return nil, types.ErrVisaNotFound
	}

	switch msg.Status {
	case "processing", "issued", "rejected", "expired":
	default:
		return nil, types.ErrInvalidStatus
	}

	visa.Status = msg.Status
	if msg.VisaNumber != "" {
		visa.VisaNumber = msg.VisaNumber
	}
	if msg.IssueDate != "" {
		visa.IssueDate = msg.IssueDate
	}
	if msg.ExpiryDate != "" {
		visa.ExpiryDate = msg.ExpiryDate
	}
	visa.UpdatedAt = nowStr(ctx)
	k.SetVisa(ctx, visa)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateVisaStatus,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyVisaID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.Status),
			sdk.NewAttribute(types.AttributeKeyVisaNo, visa.VisaNumber),
		),
	)

	return &types.MsgUpdateVisaStatusResponse{}, nil
}
