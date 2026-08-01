package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "asuransi"

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
	EventTypeCreateAsuransi      = "create_asuransi"
	EventTypeUpdateAsuransi      = "update_asuransi"
	EventTypeSubmitClaim         = "submit_claim"
	EventTypeDecideClaim         = "decide_claim"

	AttributeKeyAsuransiID  = "asuransi_id"
	AttributeKeyClaimID     = "claim_id"
	AttributeKeyJamaah      = "jamaah"
	AttributeKeyPolicyType  = "policy_type"
	AttributeKeyStatus      = "status"
	AttributeKeyCreator     = "creator"
)

// KVStore keys
var (
	// AsuransiKey prefix for asuransi store
	AsuransiKey = []byte{0x01}
	// AsuransiCountKey for the asuransi ID counter
	AsuransiCountKey = []byte{0x02}
	// ClaimKey prefix for claim store
	ClaimKey = []byte{0x03}
	// ClaimCountKey for the claim ID counter
	ClaimCountKey = []byte{0x04}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x05}
)

// GetAsuransiKey returns the store key for an asuransi by ID
func GetAsuransiKey(id uint64) []byte {
	return append(AsuransiKey, sdk.Uint64ToBigEndian(id)...)
}

// GetClaimKey returns the store key for a claim by ID
func GetClaimKey(id uint64) []byte {
	return append(ClaimKey, sdk.Uint64ToBigEndian(id)...)
}
