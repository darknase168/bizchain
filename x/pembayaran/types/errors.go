package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// pembayaran module errors
var (
	ErrPembayaranNotFound = sdkerrors.Register(ModuleName, 1, "pembayaran not found")
	ErrInvalidAmount      = sdkerrors.Register(ModuleName, 2, "invalid amount")
	ErrInstallmentNotFound = sdkerrors.Register(ModuleName, 3, "installment not found")
	ErrInstallmentPaid    = sdkerrors.Register(ModuleName, 4, "installment already paid")
	ErrEscrowStageNotFound = sdkerrors.Register(ModuleName, 5, "escrow stage not found")
	ErrEscrowReleased     = sdkerrors.Register(ModuleName, 6, "escrow stage already released")
	ErrPaymentCompleted   = sdkerrors.Register(ModuleName, 7, "payment already completed")
	ErrPaymentCancelled   = sdkerrors.Register(ModuleName, 8, "payment is cancelled")
	ErrInvalidStatus      = sdkerrors.Register(ModuleName, 9, "invalid status transition")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 10, "unauthorized access")
	ErrEscrowStageOrder   = sdkerrors.Register(ModuleName, 11, "escrow stage must be released in order")
)
