package keeper

import (
	"encoding/binary"
	"fmt"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/ticket/types"
)

// Keeper of the ticket module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new ticket Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	authority string,
) *Keeper {
	return &Keeper{
		storeKey:  storeKey,
		memKey:    memKey,
		cdc:       cdc,
		authority: authority,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) cosmossdklog.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetNextTicketID returns the next ticket ID
func (k Keeper) GetNextTicketID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TicketCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetTicketCount sets the ticket ID counter
func (k Keeper) SetTicketCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.TicketCountKey, bz)
}

// SetTicket stores a ticket
func (k Keeper) SetTicket(ctx sdk.Context, ticket types.Ticket) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TicketKey)
	bz := k.cdc.MustMarshal(&ticket)
	store.Set(sdk.Uint64ToBigEndian(ticket.Id), bz)

	currentCount := k.GetNextTicketID(ctx)
	if ticket.Id >= currentCount {
		k.SetTicketCount(ctx, ticket.Id)
	}
}

// GetTicket retrieves a ticket by ID
func (k Keeper) GetTicket(ctx sdk.Context, id uint64) (ticket types.Ticket, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TicketKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return ticket, false
	}
	k.cdc.MustUnmarshal(bz, &ticket)
	return ticket, true
}

// GetAllTicket returns all tickets
func (k Keeper) GetAllTicket(ctx sdk.Context) (list []*types.Ticket) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TicketKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var ticket types.Ticket
		k.cdc.MustUnmarshal(iterator.Value(), &ticket)
		list = append(list, &ticket)
	}

	return
}

// GetTicketByJamaah returns tickets of a jamaah
func (k Keeper) GetTicketByJamaah(ctx sdk.Context, jamaah string) (list []*types.Ticket) {
	for _, ticket := range k.GetAllTicket(ctx) {
		if ticket.Jamaah == jamaah {
			list = append(list, ticket)
		}
	}
	return
}

// RemoveTicket removes a ticket from the store
func (k Keeper) RemoveTicket(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TicketKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}

// BuildNFTID constructs a pseudo-NFT token identifier for a ticket
func (k Keeper) BuildNFTID(ctx sdk.Context, ticketID uint64) string {
	return fmt.Sprintf("nft:%s:%d", types.ModuleName, ticketID)
}
