package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "jamaah"

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
	EventTypeCreateJamaah   = "create_jamaah"
	EventTypeUpdateJamaah   = "update_jamaah"
	EventTypeAddDocument    = "add_document"
	EventTypeAddVaccination = "add_vaccination"

	AttributeKeyJamaahID   = "jamaah_id"
	AttributeKeyJamaahName = "jamaah_name"
	AttributeKeyDocType    = "doc_type"
	AttributeKeyCreator    = "creator"
)

// KVStore keys
var (
	// JamaahKey prefix for jamaah store
	JamaahKey = []byte{0x01}
	// JamaahCountKey for the jamaah ID counter
	JamaahCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetJamaahKey returns the store key for a jamaah by ID
func GetJamaahKey(id uint64) []byte {
	return append(JamaahKey, sdk.Uint64ToBigEndian(id)...)
}
