package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/bizchain/blockchain/x/audit/types"
)

var _ types.QueryServer = Keeper{}

// AuditLog handles the gRPC query for a single audit log
func (k Keeper) AuditLog(goCtx context.Context, req *types.QueryGetAuditLogRequest) (*types.QueryGetAuditLogResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	auditLog, found := k.GetAuditLog(ctx, req.Id)
	if !found {
		return nil, types.ErrAuditLogNotFound
	}

	return &types.QueryGetAuditLogResponse{AuditLog: &auditLog}, nil
}

// AuditLogAll handles the gRPC query for all audit logs
func (k Keeper) AuditLogAll(goCtx context.Context, req *types.QueryAllAuditLogRequest) (*types.QueryAllAuditLogResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	auditLog := k.GetAllAuditLog(ctx)

	return &types.QueryAllAuditLogResponse{
		AuditLog: auditLog,
	}, nil
}

// AuditLogByModule handles the gRPC query for audit logs filtered by module
func (k Keeper) AuditLogByModule(goCtx context.Context, req *types.QueryAuditLogByModuleRequest) (*types.QueryAuditLogByModuleResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	auditLog := k.GetAuditLogByModule(ctx, req.Module)

	return &types.QueryAuditLogByModuleResponse{
		AuditLog: auditLog,
	}, nil
}
