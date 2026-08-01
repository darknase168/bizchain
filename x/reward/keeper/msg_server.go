package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/reward/types"
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

// AwardReward awards loyalty tokens to a jamaah
func (k msgServer) AwardReward(goCtx context.Context, msg *types.MsgAwardReward) (*types.MsgAwardRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	points, ok := sdkmath.NewIntFromString(msg.Points)
	if !ok || !points.IsPositive() {
		return nil, types.ErrInvalidPoints
	}

	rewardID := k.GetNextRewardID(ctx)
	reward := types.Reward{
		Id:         rewardID,
		Jamaah:     msg.Jamaah,
		Points:     points.String(),
		RewardType: msg.RewardType,
		Reason:     msg.Reason,
		Status:     "awarded",
		Creator:    msg.Creator,
		CreatedAt:  nowStr(ctx),
	}
	k.SetReward(ctx, reward)

	// Update balance
	balance, found := k.GetRewardBalance(ctx, msg.Jamaah)
	if !found {
		balance = types.RewardBalance{Jamaah: msg.Jamaah, Balance: "0", Earned: "0", Redeemed: "0"}
	}
	bal, _ := sdkmath.NewIntFromString(balance.Balance)
	earned, _ := sdkmath.NewIntFromString(balance.Earned)
	balance.Balance = bal.Add(points).String()
	balance.Earned = earned.Add(points).String()
	k.SetRewardBalance(ctx, balance)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAwardReward,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyRewardID, fmt.Sprintf("%d", rewardID)),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.Jamaah),
			sdk.NewAttribute(types.AttributeKeyPoints, points.String()),
			sdk.NewAttribute(types.AttributeKeyRewardType, msg.RewardType),
		),
	)

	return &types.MsgAwardRewardResponse{Id: rewardID}, nil
}

// RedeemReward redeems loyalty tokens
func (k msgServer) RedeemReward(goCtx context.Context, msg *types.MsgRedeemReward) (*types.MsgRedeemRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	points, ok := sdkmath.NewIntFromString(msg.Points)
	if !ok || !points.IsPositive() {
		return nil, types.ErrInvalidPoints
	}

	balance, found := k.GetRewardBalance(ctx, msg.Jamaah)
	if !found {
		return nil, types.ErrInsufficientBalance
	}
	bal, _ := sdkmath.NewIntFromString(balance.Balance)
	redeemed, _ := sdkmath.NewIntFromString(balance.Redeemed)
	if bal.LT(points) {
		return nil, types.ErrInsufficientBalance
	}

	balance.Balance = bal.Sub(points).String()
	balance.Redeemed = redeemed.Add(points).String()
	k.SetRewardBalance(ctx, balance)

	rewardID := k.GetNextRewardID(ctx)
	reward := types.Reward{
		Id:         rewardID,
		Jamaah:     msg.Jamaah,
		Points:     points.String(),
		RewardType: "redeem",
		Reason:     msg.Reason,
		Status:     "redeemed",
		Creator:    msg.Creator,
		CreatedAt:  nowStr(ctx),
	}
	k.SetReward(ctx, reward)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRedeemReward,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyRewardID, fmt.Sprintf("%d", rewardID)),
			sdk.NewAttribute(types.AttributeKeyJamaah, msg.Jamaah),
			sdk.NewAttribute(types.AttributeKeyPoints, points.String()),
			sdk.NewAttribute(types.AttributeKeyBalance, balance.Balance),
		),
	)

	return &types.MsgRedeemRewardResponse{}, nil
}
