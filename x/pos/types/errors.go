package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// POS module errors
var (
	ErrProductNotFound      = sdkerrors.Register(ModuleName, 1, "product not found")
	ErrProductAlreadyExist  = sdkerrors.Register(ModuleName, 2, "product already exists")
	ErrInvalidProductID     = sdkerrors.Register(ModuleName, 3, "invalid product ID")
	ErrTransactionNotFound  = sdkerrors.Register(ModuleName, 4, "transaction not found")
	ErrInsufficientStock    = sdkerrors.Register(ModuleName, 5, "insufficient stock")
	ErrInvalidQuantity      = sdkerrors.Register(ModuleName, 6, "invalid quantity")
	ErrInvalidPrice         = sdkerrors.Register(ModuleName, 7, "invalid price")
	ErrProductInactive      = sdkerrors.Register(ModuleName, 8, "product is inactive")
	ErrUnauthorized         = sdkerrors.Register(ModuleName, 9, "unauthorized access")
	ErrUnitNotFound         = sdkerrors.Register(ModuleName, 10, "unit not found")
	ErrInvalidUnitID        = sdkerrors.Register(ModuleName, 11, "invalid unit ID")
	ErrAccountNotFound      = sdkerrors.Register(ModuleName, 12, "account not found")
	ErrInvalidAccountType   = sdkerrors.Register(ModuleName, 13, "invalid account type")
	ErrJournalNotFound      = sdkerrors.Register(ModuleName, 14, "journal entry not found")
	ErrInvalidJournalEntry  = sdkerrors.Register(ModuleName, 15, "journal entry must be balanced (debit = credit)")
	ErrInvalidBundle        = sdkerrors.Register(ModuleName, 16, "invalid bundle components")
	ErrInvalidPriceLevel    = sdkerrors.Register(ModuleName, 17, "invalid price level")
	ErrTransactionCancelled = sdkerrors.Register(ModuleName, 18, "transaction is already cancelled")
)
