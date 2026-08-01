package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateHotel{}
	_ sdk.Msg = &MsgUpdateHotel{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateHotel ----------

func (msg *MsgCreateHotel) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "name cannot be empty")
	}
	if len(msg.PricePerNight) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "price per night cannot be empty")
	}
	price, ok := sdkmath.NewIntFromString(msg.PricePerNight)
	if !ok || !price.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidPrice, "price per night must be a valid positive integer")
	}
	return nil
}

func (msg *MsgCreateHotel) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateHotel ----------

func (msg *MsgUpdateHotel) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidHotelID, "hotel ID cannot be zero")
	}
	return nil
}

func (msg *MsgUpdateHotel) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
