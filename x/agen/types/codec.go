package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateAgen{}, "agen/CreateAgen", nil)
	cdc.RegisterConcrete(&MsgUpdateAgen{}, "agen/UpdateAgen", nil)
	cdc.RegisterConcrete(&MsgAddComplaint{}, "agen/AddComplaint", nil)
	cdc.RegisterConcrete(&MsgResolveComplaint{}, "agen/ResolveComplaint", nil)
	cdc.RegisterConcrete(&MsgRecordPerformance{}, "agen/RecordPerformance", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAgen{},
		&MsgUpdateAgen{},
		&MsgAddComplaint{},
		&MsgResolveComplaint{},
		&MsgRecordPerformance{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
