package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// visa module errors
var (
	ErrVisaNotFound  = sdkerrors.Register(ModuleName, 1, "visa not found")
	ErrInvalidVisaID = sdkerrors.Register(ModuleName, 2, "invalid visa ID")
	ErrInvalidStatus = sdkerrors.Register(ModuleName, 3, "invalid visa status")
	ErrUnauthorized  = sdkerrors.Register(ModuleName, 4, "unauthorized access")
)
