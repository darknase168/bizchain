package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/audit/types"
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

// LogAction records an immutable audit trail entry
func (k msgServer) LogAction(goCtx context.Context, msg *types.MsgLogAction) (*types.MsgLogActionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	logID := k.GetNextAuditLogID(ctx)
	auditLog := types.AuditLog{
		Id:        logID,
		Module:    msg.Module,
		Action:    msg.Action,
		Actor:     msg.Actor,
		TargetId:  msg.TargetId,
		DataHash:  msg.DataHash,
		Metadata:  msg.Metadata,
		CreatedAt: nowStr(ctx),
	}
	k.SetAuditLog(ctx, auditLog)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLogAction,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyLogID, fmt.Sprintf("%d", logID)),
			sdk.NewAttribute(types.AttributeKeyAction, msg.Action),
			sdk.NewAttribute(types.AttributeKeyActor, msg.Actor),
			sdk.NewAttribute(types.AttributeKeyTargetID, msg.TargetId),
		),
	)

	return &types.MsgLogActionResponse{Id: logID}, nil
}
