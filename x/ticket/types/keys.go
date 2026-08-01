package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "ticket"

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
	EventTypeIssueTicket       = "issue_ticket"
	EventTypeUpdateTicketStatus = "update_ticket_status"

	AttributeKeyTicketID   = "ticket_id"
	AttributeKeyJamaah     = "jamaah"
	AttributeKeyStatus     = "status"
	AttributeKeyNFTID      = "nft_id"
	AttributeKeyCreator    = "creator"
)

// KVStore keys
var (
	// TicketKey prefix for ticket store
	TicketKey = []byte{0x01}
	// TicketCountKey for the ticket ID counter
	TicketCountKey = []byte{0x02}
	// ParamsKey for the module parameters
	ParamsKey = []byte{0x03}
)

// GetTicketKey returns the store key for a ticket by ID
func GetTicketKey(id uint64) []byte {
	return append(TicketKey, sdk.Uint64ToBigEndian(id)...)
}
