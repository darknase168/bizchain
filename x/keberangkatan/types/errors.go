package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// keberangkatan module errors
var (
	ErrKeberangkatanNotFound = sdkerrors.Register(ModuleName, 1, "keberangkatan not found")
	ErrInvalidKeberangkatanID = sdkerrors.Register(ModuleName, 2, "invalid keberangkatan ID")
	ErrInvalidStage          = sdkerrors.Register(ModuleName, 3, "invalid stage")
	ErrJourneyCompleted      = sdkerrors.Register(ModuleName, 4, "journey already completed")
	ErrBaggageNotFound       = sdkerrors.Register(ModuleName, 5, "baggage not found")
	ErrInvalidBaggageStatus  = sdkerrors.Register(ModuleName, 6, "invalid baggage status")
	ErrUnauthorized          = sdkerrors.Register(ModuleName, 7, "unauthorized access")
)
