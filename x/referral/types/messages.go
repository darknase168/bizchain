package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateReferral{}
	_ sdk.Msg = &MsgSettleReferral{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateReferral ----------

func (msg *MsgCreateReferral) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Agent) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "agent cannot be empty")
	}
	if len(msg.ReferredJamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "referred jamaah cannot be empty")
	}
	if len(msg.CommissionRate) > 0 {
		if _, ok := sdkmath.NewIntFromString(msg.CommissionRate); !ok {
			return sdkerrors.Wrap(ErrInvalidRate, "commission rate must be a valid integer")
		}
	}
	return nil
}

func (msg *MsgCreateReferral) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgSettleReferral ----------

func (msg *MsgSettleReferral) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidReferralID, "referral ID cannot be zero")
	}
	return nil
}

func (msg *MsgSettleReferral) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
