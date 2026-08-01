package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// paket module errors
var (
	ErrPaketNotFound  = sdkerrors.Register(ModuleName, 1, "paket not found")
	ErrInvalidPaketID = sdkerrors.Register(ModuleName, 2, "invalid paket ID")
	ErrQuotaFull      = sdkerrors.Register(ModuleName, 3, "paket quota is full")
	ErrInvalidQuota   = sdkerrors.Register(ModuleName, 4, "invalid quota or booked count")
	ErrInvalidRating  = sdkerrors.Register(ModuleName, 5, "rating must be between 1 and 5")
	ErrUnauthorized   = sdkerrors.Register(ModuleName, 6, "unauthorized access")
)
