package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// audit module errors
var (
	ErrAuditLogNotFound  = sdkerrors.Register(ModuleName, 1, "audit log not found")
	ErrInvalidAuditLogID = sdkerrors.Register(ModuleName, 2, "invalid audit log ID")
	ErrInvalidAction     = sdkerrors.Register(ModuleName, 3, "invalid action")
	ErrUnauthorized      = sdkerrors.Register(ModuleName, 4, "unauthorized access")
)
