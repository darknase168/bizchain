package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// ticket module errors
var (
	ErrTicketNotFound  = sdkerrors.Register(ModuleName, 1, "ticket not found")
	ErrInvalidTicketID = sdkerrors.Register(ModuleName, 2, "invalid ticket ID")
	ErrInvalidStatus   = sdkerrors.Register(ModuleName, 3, "invalid ticket status")
	ErrUnauthorized    = sdkerrors.Register(ModuleName, 4, "unauthorized access")
)
