package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pembayaran/types"
)

// Keeper of the pembayaran module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new pembayaran Keeper instance
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

// GetNextPembayaranID returns the next pembayaran ID
func (k Keeper) GetNextPembayaranID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PembayaranCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetPembayaranCount sets the pembayaran ID counter
func (k Keeper) SetPembayaranCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.PembayaranCountKey, bz)
}

// GetNextInstallmentID returns the next installment ID
func (k Keeper) GetNextInstallmentID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InstallmentCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetInstallmentCount sets the installment ID counter
func (k Keeper) SetInstallmentCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.InstallmentCountKey, bz)
}

// SetPembayaran stores a pembayaran
func (k Keeper) SetPembayaran(ctx sdk.Context, pembayaran types.Pembayaran) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PembayaranKey)
	bz := k.cdc.MustMarshal(&pembayaran)
	store.Set(sdk.Uint64ToBigEndian(pembayaran.Id), bz)

	currentCount := k.GetNextPembayaranID(ctx)
	if pembayaran.Id >= currentCount {
		k.SetPembayaranCount(ctx, pembayaran.Id)
	}
}

// GetPembayaran retrieves a pembayaran by ID
func (k Keeper) GetPembayaran(ctx sdk.Context, id uint64) (pembayaran types.Pembayaran, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PembayaranKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return pembayaran, false
	}
	k.cdc.MustUnmarshal(bz, &pembayaran)
	return pembayaran, true
}

// GetAllPembayaran returns all pembayaran
func (k Keeper) GetAllPembayaran(ctx sdk.Context) (list []*types.Pembayaran) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PembayaranKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pembayaran types.Pembayaran
		k.cdc.MustUnmarshal(iterator.Value(), &pembayaran)
		list = append(list, &pembayaran)
	}

	return
}

// GetPembayaranByJamaah returns all pembayaran of a jamaah
func (k Keeper) GetPembayaranByJamaah(ctx sdk.Context, jamaah string) (list []*types.Pembayaran) {
	for _, pembayaran := range k.GetAllPembayaran(ctx) {
		if pembayaran.Jamaah == jamaah {
			list = append(list, pembayaran)
		}
	}
	return
}

// RemovePembayaran removes a pembayaran from the store
func (k Keeper) RemovePembayaran(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PembayaranKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
