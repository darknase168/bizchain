package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "referral"

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
	EventTypeCreateReferral = "create_referral"
	EventTypeSettleReferral = "settle_referral"

	AttributeKeyReferralID = "referral_id"
	AttributeKeyAgent      = "agent"
	AttributeKeyJamaah     = "referred_jamaah"
	AttributeKeyCommission = "commission"
	AttributeKeyStatus     = "status"
	AttributeKeyCreator    = "creator"
)

// KVStore keys
var (
	// ReferralKey prefix for referral store
	ReferralKey = []byte{0x01}
	// ReferralCountKey for the referral ID counter
	ReferralCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetReferralKey returns the store key for a referral by ID
func GetReferralKey(id uint64) []byte {
	return append(ReferralKey, sdk.Uint64ToBigEndian(id)...)
}
