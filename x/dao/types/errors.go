package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// dao module errors
var (
	ErrProposalNotFound = sdkerrors.Register(ModuleName, 1, "proposal not found")
	ErrInvalidProposalID = sdkerrors.Register(ModuleName, 2, "invalid proposal id")
	ErrInvalidOption    = sdkerrors.Register(ModuleName, 3, "invalid vote option")
	ErrAlreadyVoted     = sdkerrors.Register(ModuleName, 4, "already voted")
	ErrProposalClosed   = sdkerrors.Register(ModuleName, 5, "proposal already closed")
	ErrProposalActive   = sdkerrors.Register(ModuleName, 6, "proposal still active")
	ErrDeadlinePassed   = sdkerrors.Register(ModuleName, 7, "voting deadline has passed")
	ErrUnauthorized     = sdkerrors.Register(ModuleName, 8, "unauthorized access")
)
