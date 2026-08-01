package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateAsuransi{}
	_ sdk.Msg = &MsgUpdateAsuransiStatus{}
	_ sdk.Msg = &MsgSubmitClaim{}
	_ sdk.Msg = &MsgDecideClaim{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateAsuransi ----------

func (msg *MsgCreateAsuransi) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Jamaah) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "jamaah cannot be empty")
	}
	if len(msg.PolicyType) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "policy type cannot be empty")
	}
	if len(msg.Premium) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "premium cannot be empty")
	}
	return nil
}

func (msg *MsgCreateAsuransi) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgUpdateAsuransiStatus ----------

func (msg *MsgUpdateAsuransiStatus) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidAsuransiID, "asuransi ID cannot be zero")
	}
	switch msg.Status {
	case "active", "expired", "cancelled", "claimed":
	default:
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "invalid status")
	}
	return nil
}

func (msg *MsgUpdateAsuransiStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgSubmitClaim ----------

func (msg *MsgSubmitClaim) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.AsuransiId == 0 {
		return sdkerrors.Wrap(ErrInvalidAsuransiID, "asuransi ID cannot be zero")
	}
	if len(msg.Reason) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "reason cannot be empty")
	}
	if len(msg.Amount) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "amount cannot be empty")
	}
	return nil
}

func (msg *MsgSubmitClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgDecideClaim ----------

func (msg *MsgDecideClaim) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.ClaimId == 0 {
		return sdkerrors.Wrap(ErrInvalidClaimID, "claim ID cannot be zero")
	}
	switch msg.Status {
	case "approved", "rejected":
	default:
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "invalid decision status")
	}
	return nil
}

func (msg *MsgDecideClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
