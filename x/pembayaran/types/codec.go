package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePembayaran{}, "pembayaran/CreatePembayaran", nil)
	cdc.RegisterConcrete(&MsgPayInstallment{}, "pembayaran/PayInstallment", nil)
	cdc.RegisterConcrete(&MsgReleaseEscrow{}, "pembayaran/ReleaseEscrow", nil)
	cdc.RegisterConcrete(&MsgRefundPembayaran{}, "pembayaran/RefundPembayaran", nil)
	cdc.RegisterConcrete(&MsgCancelPembayaran{}, "pembayaran/CancelPembayaran", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePembayaran{},
		&MsgPayInstallment{},
		&MsgReleaseEscrow{},
		&MsgRefundPembayaran{},
		&MsgCancelPembayaran{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
