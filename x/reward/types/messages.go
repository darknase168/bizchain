package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgAwardReward{}
	_ sdk.Msg = &MsgRedeemReward{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func validatePoints(name, value string) error {
	points, ok := sdkmath.NewIntFromString(value)
	if !ok || !points.IsPositive() {
		return sdkerrors.Wrapf(ErrInvalidPoints, "%s must be a valid positive integer", name)
	}
	return nil
}

// ---------- MsgAwardReward ----------

func (msg *MsgAwardReward) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Jamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "jamaah cannot be empty")
	}
	if err := validatePoints("points", msg.Points); err != nil {
		return err
	}
	return nil
}

func (msg *MsgAwardReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgRedeemReward ----------

func (msg *MsgRedeemReward) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Jamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "jamaah cannot be empty")
	}
	if err := validatePoints("points", msg.Points); err != nil {
		return err
	}
	return nil
}

func (msg *MsgRedeemReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
