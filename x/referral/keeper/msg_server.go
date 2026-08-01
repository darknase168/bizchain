package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/referral/types"
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

// CreateReferral records a referral relationship
func (k msgServer) CreateReferral(goCtx context.Context, msg *types.MsgCreateReferral) (*types.MsgCreateReferralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	rate := msg.CommissionRate
	if rate == "" {
		rate = types.DefaultParams().DefaultCommissionRate
	}
	if _, ok := sdkmath.NewIntFromString(rate); !ok {
		return nil, types.ErrInvalidRate
	}

	referralID := k.GetNextReferralID(ctx)
	referral := types.Referral{
		Id:             referralID,
		Agent:          msg.Agent,
		ReferredJamaah: msg.ReferredJamaah,
		PaketId:        msg.PaketId,
		CommissionRate: rate,
		Commission:     "0",
		Status:         "pending",
		Creator:        msg.Creator,
		CreatedAt:      nowStr(ctx),
	}
	k.SetReferral(ctx, referral)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateReferral,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyReferralID, fmt.Sprintf("%d", referralID)),
			sdk.NewAttribute(types.AttributeKeyAgent, msg.Agent),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.ReferredJamaah),
		),
	)

	return &types.MsgCreateReferralResponse{Id: referralID}, nil
}

// SettleReferral pays out commission to the agent
func (k msgServer) SettleReferral(goCtx context.Context, msg *types.MsgSettleReferral) (*types.MsgSettleReferralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	referral, found := k.GetReferral(ctx, msg.Id)
	if !found {
		return nil, types.ErrReferralNotFound
	}
	if referral.Status == "paid" {
		return nil, types.ErrAlreadySettled
	}

	// commission = rate basis poin dari nilai paket (untuk demo: 1 paket = 10_000_000)
	rate, ok := sdkmath.NewIntFromString(referral.CommissionRate)
	if !ok {
		return nil, types.ErrInvalidRate
	}
	paketValue := sdkmath.NewInt(10_000_000)
	commission := paketValue.Mul(rate).QuoRaw(1000) // rate basis poin / 1000 -> persen

	referral.Commission = commission.String()
	referral.Status = "paid"
	referral.PaidAt = nowStr(ctx)
	k.SetReferral(ctx, referral)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSettleReferral,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyReferralID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyCommission, commission.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, "paid"),
		),
	)

	return &types.MsgSettleReferralResponse{Commission: commission.String()}, nil
}
