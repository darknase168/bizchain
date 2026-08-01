package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/paket/types"
)

// Keeper of the paket module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new paket Keeper instance
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

// GetNextPaketID returns the next paket ID
func (k Keeper) GetNextPaketID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PaketCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetPaketCount sets the paket ID counter
func (k Keeper) SetPaketCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.PaketCountKey, bz)
}

// SetPaket stores a paket
func (k Keeper) SetPaket(ctx sdk.Context, paket types.Paket) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaketKey)
	bz := k.cdc.MustMarshal(&paket)
	store.Set(sdk.Uint64ToBigEndian(paket.Id), bz)

	currentCount := k.GetNextPaketID(ctx)
	if paket.Id >= currentCount {
		k.SetPaketCount(ctx, paket.Id)
	}
}

// GetPaket retrieves a paket by ID
func (k Keeper) GetPaket(ctx sdk.Context, id uint64) (paket types.Paket, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaketKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return paket, false
	}
	k.cdc.MustUnmarshal(bz, &paket)
	return paket, true
}

// GetAllPaket returns all paket
func (k Keeper) GetAllPaket(ctx sdk.Context) (list []*types.Paket) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaketKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var paket types.Paket
		k.cdc.MustUnmarshal(iterator.Value(), &paket)
		list = append(list, &paket)
	}

	return
}

// GetOpenPaket returns paket that are still open for booking
func (k Keeper) GetOpenPaket(ctx sdk.Context) (list []*types.Paket) {
	for _, paket := range k.GetAllPaket(ctx) {
		if paket.Status == "open" || paket.Status == "" {
			list = append(list, paket)
		}
	}
	return
}

// RemovePaket removes a paket from the store
func (k Keeper) RemovePaket(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaketKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
