package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreatePaket{}
	_ sdk.Msg = &MsgUpdatePaket{}
	_ sdk.Msg = &MsgBookPaket{}
	_ sdk.Msg = &MsgAddReview{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreatePaket ----------

func (msg *MsgCreatePaket) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "name cannot be empty")
	}
	if msg.Quota == 0 {
		return sdkerrors.Wrap(ErrInvalidQuota, "quota must be positive")
	}
	if len(msg.Price) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "price cannot be empty")
	}
	price, ok := sdkmath.NewIntFromString(msg.Price)
	if !ok || !price.IsPositive() {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "price must be a valid positive integer")
	}
	return nil
}

func (msg *MsgCreatePaket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdatePaket ----------

func (msg *MsgUpdatePaket) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidPaketID, "paket ID cannot be zero")
	}
	return nil
}

func (msg *MsgUpdatePaket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgBookPaket ----------

func (msg *MsgBookPaket) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.PaketId == 0 {
		return sdkerrors.Wrap(ErrInvalidPaketID, "paket ID cannot be zero")
	}
	if msg.Quantity == 0 {
		return sdkerrors.Wrap(ErrInvalidQuota, "quantity must be positive")
	}
	return nil
}

func (msg *MsgBookPaket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgAddReview ----------

func (msg *MsgAddReview) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.PaketId == 0 {
		return sdkerrors.Wrap(ErrInvalidPaketID, "paket ID cannot be zero")
	}
	if msg.Rating < 1 || msg.Rating > 5 {
		return sdkerrors.Wrap(ErrInvalidRating, "rating must be between 1 and 5")
	}
	return nil
}

func (msg *MsgAddReview) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
