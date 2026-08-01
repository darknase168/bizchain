package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// agen module errors
var (
	ErrAgenNotFound       = sdkerrors.Register(ModuleName, 1, "agen not found")
	ErrInvalidAgenID      = sdkerrors.Register(ModuleName, 2, "invalid agen id")
	ErrInvalidLevel       = sdkerrors.Register(ModuleName, 3, "invalid agent level (pusat, cabang, subagen)")
	ErrComplaintNotFound  = sdkerrors.Register(ModuleName, 4, "complaint not found")
	ErrComplaintResolved  = sdkerrors.Register(ModuleName, 5, "complaint already resolved")
	ErrInvalidComplaintID = sdkerrors.Register(ModuleName, 6, "invalid complaint id")
	ErrInvalidParent      = sdkerrors.Register(ModuleName, 7, "invalid parent agent id")
	ErrInvalidRate        = sdkerrors.Register(ModuleName, 8, "invalid commission rate")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 9, "unauthorized access")
)
