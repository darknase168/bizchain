package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "agen"

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
	EventTypeCreateAgen        = "create_agen"
	EventTypeUpdateAgen        = "update_agen"
	EventTypeAddComplaint      = "add_complaint"
	EventTypeResolveComplaint  = "resolve_complaint"
	EventTypeRecordPerformance = "record_performance"

	AttributeKeyAgenID         = "agen_id"
	AttributeKeyAgenName       = "agen_name"
	AttributeKeyAgenLevel      = "agen_level"
	AttributeKeyParentID       = "parent_id"
	AttributeKeyComplaintID    = "complaint_id"
	AttributeKeyCreator        = "creator"
)

// KVStore keys
var (
	// AgenKey prefix for agen store
	AgenKey = []byte{0x01}
	// AgenCountKey for the agen ID counter
	AgenCountKey = []byte{0x02}
	// ComplaintCountKey for the complaint ID counter
	ComplaintCountKey = []byte{0x03}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x04}
)

// GetAgenKey returns the store key for an agen by ID
func GetAgenKey(id uint64) []byte {
	return append(AgenKey, sdk.Uint64ToBigEndian(id)...)
}
