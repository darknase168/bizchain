package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// referral module errors
var (
	ErrReferralNotFound  = sdkerrors.Register(ModuleName, 1, "referral not found")
	ErrInvalidReferralID = sdkerrors.Register(ModuleName, 2, "invalid referral ID")
	ErrInvalidRate       = sdkerrors.Register(ModuleName, 3, "invalid commission rate")
	ErrAlreadySettled    = sdkerrors.Register(ModuleName, 4, "referral already settled")
	ErrUnauthorized      = sdkerrors.Register(ModuleName, 5, "unauthorized access")
)
