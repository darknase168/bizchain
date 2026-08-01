package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateProduct{}
	_ sdk.Msg = &MsgUpdateProduct{}
	_ sdk.Msg = &MsgOrderProduct{}
	_ sdk.Msg = &MsgUpdateOrderStatus{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateProduct ----------

func (msg *MsgCreateProduct) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "name cannot be empty")
	}
	if len(msg.Price) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "price cannot be empty")
	}
	return nil
}

func (msg *MsgCreateProduct) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateProduct ----------

func (msg *MsgUpdateProduct) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidProductID, "product ID cannot be zero")
	}
	return nil
}

func (msg *MsgUpdateProduct) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgOrderProduct ----------

func (msg *MsgOrderProduct) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.ProductId == 0 {
		return sdkerrors.Wrap(ErrInvalidProductID, "product ID cannot be zero")
	}
	if msg.Quantity == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "quantity must be positive")
	}
	return nil
}

func (msg *MsgOrderProduct) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateOrderStatus ----------

func (msg *MsgUpdateOrderStatus) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.OrderId == 0 {
		return sdkerrors.Wrap(ErrInvalidOrderID, "order ID cannot be zero")
	}
	switch msg.Status {
	case "paid", "shipped", "delivered", "cancelled":
	default:
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "invalid order status")
	}
	return nil
}

func (msg *MsgUpdateOrderStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
