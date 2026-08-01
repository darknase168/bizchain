package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/asuransi/types"
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

// CreateAsuransi issues a digital insurance policy
func (k msgServer) CreateAsuransi(goCtx context.Context, msg *types.MsgCreateAsuransi) (*types.MsgCreateAsuransiResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	premium, ok := sdkmath.NewIntFromString(msg.Premium)
	if !ok || !premium.IsPositive() {
		return nil, types.ErrInvalidAmount
	}

	asuransiID := k.GetNextAsuransiID(ctx)
	asuransi := types.Asuransi{
		Id:           asuransiID,
		Jamaah:       msg.Jamaah,
		PolicyType:   msg.PolicyType,
		Premium:      msg.Premium,
		Coverage:     msg.Coverage,
		StartDate:    msg.StartDate,
		EndDate:      msg.EndDate,
		Status:       "active",
		DocumentHash: msg.DocumentHash,
		Provider:     msg.Provider,
		Creator:      msg.Creator,
		CreatedAt:    nowStr(ctx),
		UpdatedAt:    nowStr(ctx),
	}
	k.SetAsuransi(ctx, asuransi)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateAsuransi,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAsuransiID, fmt.Sprintf("%d", asuransiID)),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.Jamaah),
			sdk.NewAttribute(types.AttributeKeyPolicyType, msg.PolicyType),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgCreateAsuransiResponse{Id: asuransiID}, nil
}

// UpdateAsuransiStatus updates policy status
func (k msgServer) UpdateAsuransiStatus(goCtx context.Context, msg *types.MsgUpdateAsuransiStatus) (*types.MsgUpdateAsuransiStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	asuransi, found := k.GetAsuransi(ctx, msg.Id)
	if !found {
		return nil, types.ErrAsuransiNotFound
	}
	if asuransi.Creator != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	asuransi.Status = msg.Status
	asuransi.UpdatedAt = nowStr(ctx)
	k.SetAsuransi(ctx, asuransi)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateAsuransi,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAsuransiID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.Status),
		),
	)

	return &types.MsgUpdateAsuransiStatusResponse{}, nil
}

// SubmitClaim files an insurance claim
func (k msgServer) SubmitClaim(goCtx context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	asuransi, found := k.GetAsuransi(ctx, msg.AsuransiId)
	if !found {
		return nil, types.ErrAsuransiNotFound
	}

	amount, ok := sdkmath.NewIntFromString(msg.Amount)
	if !ok || !amount.IsPositive() {
		return nil, types.ErrInvalidAmount
	}
	// cap the claim against the policy coverage
	if asuransi.Coverage != "" {
		coverage, cok := sdkmath.NewIntFromString(asuransi.Coverage)
		if cok && amount.GT(coverage) {
			return nil, types.ErrInvalidAmount
		}
	}

	claimID := k.GetNextClaimID(ctx)
	claim := types.AsuransiClaim{
		Id:           claimID,
		AsuransiId:   asuransi.Id,
		Jamaah:       asuransi.Jamaah,
		Reason:       msg.Reason,
		Amount:       msg.Amount,
		Status:       "submitted",
		EvidenceHash: msg.EvidenceHash,
		DecisionBy:   "",
		DecisionNote: "",
		SubmittedAt:  nowStr(ctx),
		DecidedAt:    "",
	}
	k.SetClaim(ctx, claim)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitClaim,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyClaimID, fmt.Sprintf("%d", claimID)),
			sdk.NewAttribute(types.AttributeKeyAsuransiID, fmt.Sprintf("%d", asuransi.Id)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgSubmitClaimResponse{ClaimId: claimID}, nil
}

// DecideClaim approves or rejects a claim
func (k msgServer) DecideClaim(goCtx context.Context, msg *types.MsgDecideClaim) (*types.MsgDecideClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	claim, found := k.GetClaim(ctx, msg.ClaimId)
	if !found {
		return nil, types.ErrClaimNotFound
	}
	if claim.Status == "approved" || claim.Status == "rejected" || claim.Status == "paid" {
		return nil, types.ErrClaimDecided
	}

	claim.Status = msg.Status
	claim.DecisionBy = msg.Creator
	claim.DecisionNote = msg.Note
	claim.DecidedAt = nowStr(ctx)
	k.SetClaim(ctx, claim)

	// mark the policy as claimed when approved
	if msg.Status == "approved" {
		if asuransi, f := k.GetAsuransi(ctx, claim.AsuransiId); f {
			asuransi.Status = "claimed"
			asuransi.UpdatedAt = nowStr(ctx)
			k.SetAsuransi(ctx, asuransi)
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDecideClaim,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyClaimID, fmt.Sprintf("%d", msg.ClaimId)),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.Status),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgDecideClaimResponse{}, nil
}
