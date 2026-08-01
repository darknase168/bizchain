package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// asuransi module errors
var (
	ErrAsuransiNotFound   = sdkerrors.Register(ModuleName, 1, "asuransi not found")
	ErrInvalidAsuransiID  = sdkerrors.Register(ModuleName, 2, "invalid asuransi id")
	ErrClaimNotFound      = sdkerrors.Register(ModuleName, 3, "claim not found")
	ErrInvalidClaimID     = sdkerrors.Register(ModuleName, 4, "invalid claim id")
	ErrInvalidAmount      = sdkerrors.Register(ModuleName, 5, "invalid amount")
	ErrInvalidStatus      = sdkerrors.Register(ModuleName, 6, "invalid status transition")
	ErrClaimDecided       = sdkerrors.Register(ModuleName, 7, "claim already decided")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 8, "unauthorized access")
)
