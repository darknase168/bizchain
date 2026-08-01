package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// reward module errors
var (
	ErrRewardNotFound  = sdkerrors.Register(ModuleName, 1, "reward not found")
	ErrInvalidRewardID = sdkerrors.Register(ModuleName, 2, "invalid reward ID")
	ErrInvalidPoints   = sdkerrors.Register(ModuleName, 3, "invalid points amount")
	ErrInsufficientBalance = sdkerrors.Register(ModuleName, 4, "insufficient reward balance")
	ErrUnauthorized   = sdkerrors.Register(ModuleName, 5, "unauthorized access")
)
