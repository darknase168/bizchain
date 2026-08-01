package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/agen/types"
)

// Keeper of the agen module
type Keeper struct {
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	cdc      codec.BinaryCodec
	authority string
}

// NewKeeper creates a new agen Keeper instance
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

// GetNextAgenID returns the next agen ID
func (k Keeper) GetNextAgenID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AgenCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetAgenCount sets the agen ID counter
func (k Keeper) SetAgenCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.AgenCountKey, bz)
}

// GetNextComplaintID returns the next complaint ID
func (k Keeper) GetNextComplaintID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ComplaintCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetComplaintCount sets the complaint ID counter
func (k Keeper) SetComplaintCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.ComplaintCountKey, bz)
}

// SetAgen stores an agen
func (k Keeper) SetAgen(ctx sdk.Context, agen types.Agen) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AgenKey)
	bz := k.cdc.MustMarshal(&agen)
	store.Set(sdk.Uint64ToBigEndian(agen.Id), bz)

	currentCount := k.GetNextAgenID(ctx)
	if agen.Id >= currentCount {
		k.SetAgenCount(ctx, agen.Id)
	}
}

// GetAgen retrieves an agen by ID
func (k Keeper) GetAgen(ctx sdk.Context, id uint64) (agen types.Agen, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AgenKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return agen, false
	}
	k.cdc.MustUnmarshal(bz, &agen)
	return agen, true
}

// GetAgenByAddress retrieves an agen by wallet address
func (k Keeper) GetAgenByAddress(ctx sdk.Context, address string) (agen types.Agen, found bool) {
	for _, a := range k.GetAllAgen(ctx) {
		if a.Address == address {
			return *a, true
		}
	}
	return agen, false
}

// GetAllAgen returns all agen
func (k Keeper) GetAllAgen(ctx sdk.Context) (list []*types.Agen) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AgenKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var agen types.Agen
		k.cdc.MustUnmarshal(iterator.Value(), &agen)
		list = append(list, &agen)
	}

	return
}

// GetDownline returns all direct downline of an agen (children)
func (k Keeper) GetDownline(ctx sdk.Context, parentID uint64) (list []*types.Agen) {
	for _, a := range k.GetAllAgen(ctx) {
		pid, err := ParseUint64(a.ParentId)
		if err == nil && pid == parentID {
			list = append(list, a)
		}
	}
	return
}

// ParseUint64 parses a string into uint64
func ParseUint64(s string) (uint64, error) {
	var v uint64
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, types.ErrInvalidAgenID
		}
		v = v*10 + uint64(s[i]-'0')
	}
	return v, nil
}
