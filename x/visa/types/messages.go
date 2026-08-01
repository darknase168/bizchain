package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateVisa{}
	_ sdk.Msg = &MsgUpdateVisaStatus{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateVisa ----------

func (msg *MsgCreateVisa) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Jamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "jamaah cannot be empty")
	}
	return nil
}

func (msg *MsgCreateVisa) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateVisaStatus ----------

func (msg *MsgUpdateVisaStatus) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidVisaID, "visa ID cannot be zero")
	}
	switch msg.Status {
	case "processing", "issued", "rejected", "expired":
	default:
		return sdkerrors.Wrap(ErrInvalidStatus, "status must be processing, issued, rejected or expired")
	}
	return nil
}

func (msg *MsgUpdateVisaStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
