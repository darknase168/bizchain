package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "reward"

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
	EventTypeAwardReward  = "award_reward"
	EventTypeRedeemReward = "redeem_reward"

	AttributeKeyRewardID   = "reward_id"
	AttributeKeyJamaah     = "jamaah"
	AttributeKeyPoints     = "points"
	AttributeKeyRewardType = "reward_type"
	AttributeKeyBalance    = "balance"
	AttributeKeyStatus     = "status"
	AttributeKeyCreator    = "creator"
)

// KVStore keys
var (
	// RewardKey prefix for reward store
	RewardKey = []byte{0x01}
	// BalanceKey prefix for reward balance store (keyed by jamaah address)
	BalanceKey = []byte{0x02}
	// RewardCountKey for the reward ID counter
	RewardCountKey = []byte{0x03}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x04}
)

// GetRewardKey returns the store key for a reward by ID
func GetRewardKey(id uint64) []byte {
	return append(RewardKey, sdk.Uint64ToBigEndian(id)...)
}

// GetBalanceKey returns the store key for a jamaah balance
func GetBalanceKey(jamaah string) []byte {
	return append(BalanceKey, []byte(jamaah)...)
}
