package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "dao"

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
	EventTypeCreateProposal = "create_proposal"
	EventTypeCastVote       = "cast_vote"
	EventTypeCloseProposal  = "close_proposal"

	AttributeKeyProposalID = "proposal_id"
	AttributeKeyTitle      = "title"
	AttributeKeyOption     = "option"
	AttributeKeyVoter      = "voter"
	AttributeKeyResult     = "result"
	AttributeKeyCreator    = "creator"
)

// KVStore keys
var (
	// ProposalKey prefix for proposal store
	ProposalKey = []byte{0x01}
	// ProposalCountKey for the proposal ID counter
	ProposalCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetProposalKey returns the store key for a proposal by ID
func GetProposalKey(id uint64) []byte {
	return append(ProposalKey, sdk.Uint64ToBigEndian(id)...)
}
