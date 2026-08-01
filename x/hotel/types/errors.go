package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// hotel module errors
var (
	ErrHotelNotFound  = sdkerrors.Register(ModuleName, 1, "hotel not found")
	ErrInvalidHotelID = sdkerrors.Register(ModuleName, 2, "invalid hotel ID")
	ErrInvalidPrice   = sdkerrors.Register(ModuleName, 3, "invalid price")
	ErrUnauthorized   = sdkerrors.Register(ModuleName, 4, "unauthorized access")
)
