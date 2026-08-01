package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateAsuransi{}, "asuransi/CreateAsuransi", nil)
	cdc.RegisterConcrete(&MsgUpdateAsuransiStatus{}, "asuransi/UpdateAsuransiStatus", nil)
	cdc.RegisterConcrete(&MsgSubmitClaim{}, "asuransi/SubmitClaim", nil)
	cdc.RegisterConcrete(&MsgDecideClaim{}, "asuransi/DecideClaim", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAsuransi{},
		&MsgUpdateAsuransiStatus{},
		&MsgSubmitClaim{},
		&MsgDecideClaim{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
