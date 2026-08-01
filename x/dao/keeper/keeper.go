package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/dao/types"
)

// Keeper of the dao module
type Keeper struct {
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	cdc      codec.BinaryCodec
	authority string
}

// NewKeeper creates a new dao Keeper instance
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

// GetNextProposalID returns the next proposal ID
func (k Keeper) GetNextProposalID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetProposalCount sets the proposal ID counter
func (k Keeper) SetProposalCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.ProposalCountKey, bz)
}

// SetProposal stores a proposal
func (k Keeper) SetProposal(ctx sdk.Context, proposal types.DaoProposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKey)
	bz := k.cdc.MustMarshal(&proposal)
	store.Set(sdk.Uint64ToBigEndian(proposal.Id), bz)

	currentCount := k.GetNextProposalID(ctx)
	if proposal.Id >= currentCount {
		k.SetProposalCount(ctx, proposal.Id)
	}
}

// GetProposal retrieves a proposal by ID
func (k Keeper) GetProposal(ctx sdk.Context, id uint64) (proposal types.DaoProposal, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return proposal, false
	}
	k.cdc.MustUnmarshal(bz, &proposal)
	return proposal, true
}

// GetAllProposals returns all proposals
func (k Keeper) GetAllProposals(ctx sdk.Context) (list []*types.DaoProposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal types.DaoProposal
		k.cdc.MustUnmarshal(iterator.Value(), &proposal)
		list = append(list, &proposal)
	}

	return
}

// HasVoted checks whether a voter has already voted on a proposal
func (k Keeper) HasVoted(ctx sdk.Context, proposalID uint64, voter string) bool {
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		return false
	}
	for _, v := range proposal.Votes {
		if v.Voter == voter {
			return true
		}
	}
	return false
}
