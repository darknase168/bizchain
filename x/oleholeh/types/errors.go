package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// oleholeh module errors
var (
	ErrProductNotFound  = sdkerrors.Register(ModuleName, 1, "product not found")
	ErrInvalidProductID = sdkerrors.Register(ModuleName, 2, "invalid product id")
	ErrOrderNotFound    = sdkerrors.Register(ModuleName, 3, "order not found")
	ErrInvalidOrderID   = sdkerrors.Register(ModuleName, 4, "invalid order id")
	ErrInvalidAmount    = sdkerrors.Register(ModuleName, 5, "invalid amount")
	ErrInsufficientStock = sdkerrors.Register(ModuleName, 6, "insufficient stock")
	ErrInvalidStatus    = sdkerrors.Register(ModuleName, 7, "invalid status transition")
	ErrUnauthorized     = sdkerrors.Register(ModuleName, 8, "unauthorized access")
)
