package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateKeberangkatan{}
	_ sdk.Msg = &MsgAdvanceStage{}
	_ sdk.Msg = &MsgUpdateDeparture{}
	_ sdk.Msg = &MsgAddBaggage{}
	_ sdk.Msg = &MsgUpdateBaggageStatus{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateKeberangkatan ----------

func (msg *MsgCreateKeberangkatan) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Jamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "jamaah cannot be empty")
	}
	return nil
}

func (msg *MsgCreateKeberangkatan) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAdvanceStage ----------

func (msg *MsgAdvanceStage) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidKeberangkatanID, "keberangkatan ID cannot be zero")
	}
	return nil
}

func (msg *MsgAdvanceStage) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateDeparture ----------

func (msg *MsgUpdateDeparture) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidKeberangkatanID, "keberangkatan ID cannot be zero")
	}
	return nil
}

func (msg *MsgUpdateDeparture) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAddBaggage ----------

func (msg *MsgAddBaggage) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidKeberangkatanID, "keberangkatan ID cannot be zero")
	}
	if len(msg.Tag) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "baggage tag cannot be empty")
	}
	return nil
}

func (msg *MsgAddBaggage) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateBaggageStatus ----------

func (msg *MsgUpdateBaggageStatus) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidKeberangkatanID, "keberangkatan ID cannot be zero")
	}
	switch msg.Status {
	case "checked_in", "in_transit", "arrived", "delivered":
	default:
		return sdkerrors.Wrap(ErrInvalidBaggageStatus, "status must be checked_in, in_transit, arrived or delivered")
	}
	return nil
}

func (msg *MsgUpdateBaggageStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
