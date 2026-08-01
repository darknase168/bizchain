package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateProduct{}, "oleholeh/CreateProduct", nil)
	cdc.RegisterConcrete(&MsgUpdateProduct{}, "oleholeh/UpdateProduct", nil)
	cdc.RegisterConcrete(&MsgOrderProduct{}, "oleholeh/OrderProduct", nil)
	cdc.RegisterConcrete(&MsgUpdateOrderStatus{}, "oleholeh/UpdateOrderStatus", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProduct{},
		&MsgUpdateProduct{},
		&MsgOrderProduct{},
		&MsgUpdateOrderStatus{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
