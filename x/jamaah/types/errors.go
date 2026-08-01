package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// jamaah module errors
var (
	ErrJamaahNotFound  = sdkerrors.Register(ModuleName, 1, "jamaah not found")
	ErrInvalidJamaahID = sdkerrors.Register(ModuleName, 2, "invalid jamaah ID")
	ErrInvalidDocType  = sdkerrors.Register(ModuleName, 3, "invalid document type")
	ErrUnauthorized    = sdkerrors.Register(ModuleName, 4, "unauthorized access")
)
