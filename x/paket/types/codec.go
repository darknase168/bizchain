package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePaket{}, "paket/CreatePaket", nil)
	cdc.RegisterConcrete(&MsgUpdatePaket{}, "paket/UpdatePaket", nil)
	cdc.RegisterConcrete(&MsgBookPaket{}, "paket/BookPaket", nil)
	cdc.RegisterConcrete(&MsgAddReview{}, "paket/AddReview", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePaket{},
		&MsgUpdatePaket{},
		&MsgBookPaket{},
		&MsgAddReview{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
