package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreatePembayaran{}
	_ sdk.Msg = &MsgPayInstallment{}
	_ sdk.Msg = &MsgReleaseEscrow{}
	_ sdk.Msg = &MsgRefundPembayaran{}
	_ sdk.Msg = &MsgCancelPembayaran{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func validateAmount(name, value string) error {
	amount, ok := sdkmath.NewIntFromString(value)
	if !ok || !amount.IsPositive() {
		return sdkerrors.Wrapf(ErrInvalidAmount, "%s must be a valid positive integer", name)
	}
	return nil
}

// ---------- MsgCreatePembayaran ----------

func (msg *MsgCreatePembayaran) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Jamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "jamaah cannot be empty")
	}
	if err := validateAmount("total", msg.Total); err != nil {
		return err
	}
	return nil
}

func (msg *MsgCreatePembayaran) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgPayInstallment ----------

func (msg *MsgPayInstallment) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.PembayaranId == 0 {
		return sdkerrors.Wrap(ErrPembayaranNotFound, "pembayaran ID cannot be zero")
	}
	if msg.InstallmentId == 0 {
		return sdkerrors.Wrap(ErrInstallmentNotFound, "installment ID cannot be zero")
	}
	if err := validateAmount("amount", msg.Amount); err != nil {
		return err
	}
	return nil
}

func (msg *MsgPayInstallment) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgReleaseEscrow ----------

func (msg *MsgReleaseEscrow) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.PembayaranId == 0 {
		return sdkerrors.Wrap(ErrPembayaranNotFound, "pembayaran ID cannot be zero")
	}
	if len(msg.StageName) == 0 {
		return sdkerrors.Wrap(ErrEscrowStageNotFound, "stage name cannot be empty")
	}
	return nil
}

func (msg *MsgReleaseEscrow) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgRefundPembayaran ----------

func (msg *MsgRefundPembayaran) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.PembayaranId == 0 {
		return sdkerrors.Wrap(ErrPembayaranNotFound, "pembayaran ID cannot be zero")
	}
	return nil
}

func (msg *MsgRefundPembayaran) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCancelPembayaran ----------

func (msg *MsgCancelPembayaran) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.PembayaranId == 0 {
		return sdkerrors.Wrap(ErrPembayaranNotFound, "pembayaran ID cannot be zero")
	}
	return nil
}

func (msg *MsgCancelPembayaran) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
