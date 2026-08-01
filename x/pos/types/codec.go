package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateProduct{}, "pos/CreateProduct", nil)
	cdc.RegisterConcrete(&MsgUpdateProduct{}, "pos/UpdateProduct", nil)
	cdc.RegisterConcrete(&MsgDeleteProduct{}, "pos/DeleteProduct", nil)
	cdc.RegisterConcrete(&MsgCreateTransaction{}, "pos/CreateTransaction", nil)
	cdc.RegisterConcrete(&MsgCancelTransaction{}, "pos/CancelTransaction", nil)
	cdc.RegisterConcrete(&MsgAddStock{}, "pos/AddStock", nil)
	cdc.RegisterConcrete(&MsgAdjustStock{}, "pos/AdjustStock", nil)
	cdc.RegisterConcrete(&MsgCreateUnit{}, "pos/CreateUnit", nil)
	cdc.RegisterConcrete(&MsgUpdateUnit{}, "pos/UpdateUnit", nil)
	cdc.RegisterConcrete(&MsgDeleteUnit{}, "pos/DeleteUnit", nil)
	cdc.RegisterConcrete(&MsgCreateAccount{}, "pos/CreateAccount", nil)
	cdc.RegisterConcrete(&MsgCreateJournalEntry{}, "pos/CreateJournalEntry", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProduct{},
		&MsgUpdateProduct{},
		&MsgDeleteProduct{},
		&MsgCreateTransaction{},
		&MsgCancelTransaction{},
		&MsgAddStock{},
		&MsgAdjustStock{},
		&MsgCreateUnit{},
		&MsgUpdateUnit{},
		&MsgDeleteUnit{},
		&MsgCreateAccount{},
		&MsgCreateJournalEntry{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
