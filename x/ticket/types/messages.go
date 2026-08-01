package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgIssueTicket{}
	_ sdk.Msg = &MsgUpdateTicketStatus{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgIssueTicket ----------

func (msg *MsgIssueTicket) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Jamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "jamaah cannot be empty")
	}
	if len(msg.Airline) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "airline cannot be empty")
	}
	if len(msg.FlightNumber) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "flight number cannot be empty")
	}
	return nil
}

func (msg *MsgIssueTicket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateTicketStatus ----------

func (msg *MsgUpdateTicketStatus) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidTicketID, "ticket ID cannot be zero")
	}
	switch msg.Status {
	case "issued", "checked_in", "boarded", "used", "void":
	default:
		return sdkerrors.Wrap(ErrInvalidStatus, "status must be issued, checked_in, boarded, used or void")
	}
	return nil
}

func (msg *MsgUpdateTicketStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
