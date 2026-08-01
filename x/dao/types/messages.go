package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
)

// ensure Msg types implement sdk.Msg
var (
	_ sdk.Msg = &MsgCreateProposal{}
	_ sdk.Msg = &MsgCastVote{}
	_ sdk.Msg = &MsgCloseProposal{}
)

func validCreator(creator string) error {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ---------- MsgCreateProposal ----------

func (msg *MsgCreateProposal) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if len(msg.Title) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "title cannot be empty")
	}
	if len(msg.Options) < 2 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "at least 2 options required")
	}
	return nil
}

func (msg *MsgCreateProposal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCastVote ----------

func (msg *MsgCastVote) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.ProposalId == 0 {
		return sdkerrors.Wrap(ErrInvalidProposalID, "proposal ID cannot be zero")
	}
	if len(msg.Option) == 0 {
		return sdkerrors.Wrap(sdkerrorstypes.ErrInvalidRequest, "option cannot be empty")
	}
	return nil
}

func (msg *MsgCastVote) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ---------- MsgCloseProposal ----------

func (msg *MsgCloseProposal) ValidateBasic() error {
	if err := validCreator(msg.Creator); err != nil {
		return err
	}
	if msg.ProposalId == 0 {
		return sdkerrors.Wrap(ErrInvalidProposalID, "proposal ID cannot be zero")
	}
	return nil
}

func (msg *MsgCloseProposal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
