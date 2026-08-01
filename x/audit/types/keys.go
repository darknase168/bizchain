package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "audit"

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
	EventTypeLogAction = "log_action"

	AttributeKeyLogID    = "log_id"
	AttributeKeyModule   = "audit_module"
	AttributeKeyAction   = "action"
	AttributeKeyActor    = "actor"
	AttributeKeyTargetID = "target_id"
	AttributeKeyCreator  = "creator"
)

// KVStore keys
var (
	// AuditLogKey prefix for audit log store
	AuditLogKey = []byte{0x01}
	// AuditCountKey for the audit log ID counter
	AuditCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetAuditLogKey returns the store key for an audit log by ID
func GetAuditLogKey(id uint64) []byte {
	return append(AuditLogKey, sdk.Uint64ToBigEndian(id)...)
}
