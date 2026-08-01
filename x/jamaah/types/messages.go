package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateJamaah{}
	_ sdk.Msg = &MsgUpdateJamaah{}
	_ sdk.Msg = &MsgAddDocument{}
	_ sdk.Msg = &MsgAddVaccination{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateJamaah ----------

func (msg *MsgCreateJamaah) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "name cannot be empty")
	}
	return nil
}

func (msg *MsgCreateJamaah) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateJamaah ----------

func (msg *MsgUpdateJamaah) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidJamaahID, "jamaah ID cannot be zero")
	}
	return nil
}

func (msg *MsgUpdateJamaah) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAddDocument ----------

func (msg *MsgAddDocument) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.JamaahId == 0 {
		return sdkerrors.Wrap(ErrInvalidJamaahID, "jamaah ID cannot be zero")
	}
	if len(msg.DocType) == 0 {
		return sdkerrors.Wrap(ErrInvalidDocType, "doc type cannot be empty")
	}
	if len(msg.Hash) == 0 {
		return sdkerrors.Wrap(ErrInvalidDocType, "document hash cannot be empty")
	}
	return nil
}

func (msg *MsgAddDocument) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAddVaccination ----------

func (msg *MsgAddVaccination) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.JamaahId == 0 {
		return sdkerrors.Wrap(ErrInvalidJamaahID, "jamaah ID cannot be zero")
	}
	if len(msg.VaccineName) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "vaccine name cannot be empty")
	}
	return nil
}

func (msg *MsgAddVaccination) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
