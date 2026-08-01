package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/jamaah/types"
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

// CreateJamaah registers a new pilgrim (self-sovereign DID on-chain)
func (k msgServer) CreateJamaah(goCtx context.Context, msg *types.MsgCreateJamaah) (*types.MsgCreateJamaahResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	jamaahID := k.GetNextJamaahID(ctx)
	jamaah := types.Jamaah{
		Id:             jamaahID,
		Name:           msg.Name,
		Phone:          msg.Phone,
		Email:          msg.Email,
		Address:        msg.Address,
		PassportNumber: msg.PassportNumber,
		PhotoHash:      msg.PhotoHash,
		Status:         "active",
		Did:            "did:point:" + msg.Creator,
		Creator:        msg.Creator,
		CreatedAt:      nowStr(ctx),
		UpdatedAt:      nowStr(ctx),
	}
	k.SetJamaah(ctx, jamaah)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateJamaah,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyJamaahID, fmt.Sprintf("%d", jamaahID)),
			sdk.NewAttribute(types.AttributeKeyJamaahName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateJamaahResponse{Id: jamaahID}, nil
}

// UpdateJamaah updates pilgrim data
func (k msgServer) UpdateJamaah(goCtx context.Context, msg *types.MsgUpdateJamaah) (*types.MsgUpdateJamaahResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	jamaah, found := k.GetJamaah(ctx, msg.Id)
	if !found {
		return nil, types.ErrJamaahNotFound
	}
	if jamaah.Creator != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	jamaah.Name = msg.Name
	jamaah.Phone = msg.Phone
	jamaah.Email = msg.Email
	jamaah.Address = msg.Address
	jamaah.PassportNumber = msg.PassportNumber
	jamaah.PhotoHash = msg.PhotoHash
	if msg.Status != "" {
		jamaah.Status = msg.Status
	}
	jamaah.UpdatedAt = nowStr(ctx)
	k.SetJamaah(ctx, jamaah)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateJamaah,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyJamaahID, fmt.Sprintf("%d", msg.Id)),
		),
	)

	return &types.MsgUpdateJamaahResponse{}, nil
}

// AddDocument adds a hashed document to a pilgrim (passport, visa, vaccine, health)
func (k msgServer) AddDocument(goCtx context.Context, msg *types.MsgAddDocument) (*types.MsgAddDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	jamaah, found := k.GetJamaah(ctx, msg.JamaahId)
	if !found {
		return nil, types.ErrJamaahNotFound
	}

	jamaah.Documents = append(jamaah.Documents, &types.DocumentHash{
		DocType:    msg.DocType,
		Hash:       msg.Hash,
		StorageRef: msg.StorageRef,
		UploadedAt: nowStr(ctx),
	})
	jamaah.UpdatedAt = nowStr(ctx)
	k.SetJamaah(ctx, jamaah)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddDocument,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyJamaahID, fmt.Sprintf("%d", msg.JamaahId)),
			sdk.NewAttribute(types.AttributeKeyDocType, msg.DocType),
		),
	)

	return &types.MsgAddDocumentResponse{}, nil
}

// AddVaccination adds a vaccination record to a pilgrim
func (k msgServer) AddVaccination(goCtx context.Context, msg *types.MsgAddVaccination) (*types.MsgAddVaccinationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	jamaah, found := k.GetJamaah(ctx, msg.JamaahId)
	if !found {
		return nil, types.ErrJamaahNotFound
	}

	jamaah.Vaccinations = append(jamaah.Vaccinations, &types.VaccinationRecord{
		VaccineName: msg.VaccineName,
		Date:        msg.Date,
		Issuer:      msg.Issuer,
		Batch:       msg.Batch,
	})
	jamaah.UpdatedAt = nowStr(ctx)
	k.SetJamaah(ctx, jamaah)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddVaccination,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyJamaahID, fmt.Sprintf("%d", msg.JamaahId)),
		),
	)

	return &types.MsgAddVaccinationResponse{}, nil
}
