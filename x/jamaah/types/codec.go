package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateJamaah{}, "jamaah/CreateJamaah", nil)
	cdc.RegisterConcrete(&MsgUpdateJamaah{}, "jamaah/UpdateJamaah", nil)
	cdc.RegisterConcrete(&MsgAddDocument{}, "jamaah/AddDocument", nil)
	cdc.RegisterConcrete(&MsgAddVaccination{}, "jamaah/AddVaccination", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateJamaah{},
		&MsgUpdateJamaah{},
		&MsgAddDocument{},
		&MsgAddVaccination{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
