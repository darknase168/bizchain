package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateAgen{}
	_ sdk.Msg = &MsgUpdateAgen{}
	_ sdk.Msg = &MsgAddComplaint{}
	_ sdk.Msg = &MsgResolveComplaint{}
	_ sdk.Msg = &MsgRecordPerformance{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateAgen ----------

func (msg *MsgCreateAgen) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "name cannot be empty")
	}
	if len(msg.Address) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "address cannot be empty")
	}
	switch msg.Level {
	case "pusat", "cabang", "subagen", "":
	default:
		return ErrInvalidLevel
	}
	return nil
}

func (msg *MsgCreateAgen) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateAgen ----------

func (msg *MsgUpdateAgen) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidAgenID, "agen ID cannot be zero")
	}
	switch msg.Level {
	case "pusat", "cabang", "subagen", "":
	default:
		return ErrInvalidLevel
	}
	return nil
}

func (msg *MsgUpdateAgen) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAddComplaint ----------

func (msg *MsgAddComplaint) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.AgenId == 0 {
		return sdkerrors.Wrap(ErrInvalidAgenID, "agen ID cannot be zero")
	}
	if len(msg.Reason) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "reason cannot be empty")
	}
	return nil
}

func (msg *MsgAddComplaint) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgResolveComplaint ----------

func (msg *MsgResolveComplaint) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.AgenId == 0 {
		return sdkerrors.Wrap(ErrInvalidAgenID, "agen ID cannot be zero")
	}
	if msg.ComplaintId == 0 {
		return sdkerrors.Wrap(ErrInvalidComplaintID, "complaint ID cannot be zero")
	}
	if len(msg.Resolution) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "resolution cannot be empty")
	}
	return nil
}

func (msg *MsgResolveComplaint) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgRecordPerformance ----------

func (msg *MsgRecordPerformance) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.AgenId == 0 {
		return sdkerrors.Wrap(ErrInvalidAgenID, "agen ID cannot be zero")
	}
	if len(msg.Period) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "period cannot be empty")
	}
	return nil
}

func (msg *MsgRecordPerformance) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
