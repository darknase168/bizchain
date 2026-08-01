package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "visa"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Event attribute keys
const (
	EventTypeCreateVisa       = "create_visa"
	EventTypeUpdateVisaStatus = "update_visa_status"

	AttributeKeyVisaID   = "visa_id"
	AttributeKeyJamaah   = "jamaah"
	AttributeKeyStatus   = "status"
	AttributeKeyVisaNo   = "visa_number"
	AttributeKeyCreator  = "creator"
)

// KVStore keys
var (
	// VisaKey prefix for visa store
	VisaKey = []byte{0x01}
	// VisaCountKey for the visa ID counter
	VisaCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetVisaKey returns the store key for a visa by ID
func GetVisaKey(id uint64) []byte {
	return append(VisaKey, sdk.Uint64ToBigEndian(id)...)
}
