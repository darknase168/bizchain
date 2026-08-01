package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateKeberangkatan{}, "keberangkatan/CreateKeberangkatan", nil)
	cdc.RegisterConcrete(&MsgAdvanceStage{}, "keberangkatan/AdvanceStage", nil)
	cdc.RegisterConcrete(&MsgUpdateDeparture{}, "keberangkatan/UpdateDeparture", nil)
	cdc.RegisterConcrete(&MsgAddBaggage{}, "keberangkatan/AddBaggage", nil)
	cdc.RegisterConcrete(&MsgUpdateBaggageStatus{}, "keberangkatan/UpdateBaggageStatus", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateKeberangkatan{},
		&MsgAdvanceStage{},
		&MsgUpdateDeparture{},
		&MsgAddBaggage{},
		&MsgUpdateBaggageStatus{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
