package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
)

// RegisterLegacyAminoCodec registers legacy amino codecs via the module basics.
func RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	std.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterLegacyAminoCodec(amino)
}

// RegisterInterfaces registers the module interfaces.
func RegisterInterfaces(interfaceRegistry types.InterfaceRegistry) {
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
}
