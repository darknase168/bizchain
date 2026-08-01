package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/visa/types"
)

// Keeper of the visa module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new visa Keeper instance
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

// GetNextVisaID returns the next visa ID
func (k Keeper) GetNextVisaID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VisaCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetVisaCount sets the visa ID counter
func (k Keeper) SetVisaCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.VisaCountKey, bz)
}

// SetVisa stores a visa
func (k Keeper) SetVisa(ctx sdk.Context, visa types.Visa) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VisaKey)
	bz := k.cdc.MustMarshal(&visa)
	store.Set(sdk.Uint64ToBigEndian(visa.Id), bz)

	currentCount := k.GetNextVisaID(ctx)
	if visa.Id >= currentCount {
		k.SetVisaCount(ctx, visa.Id)
	}
}

// GetVisa retrieves a visa by ID
func (k Keeper) GetVisa(ctx sdk.Context, id uint64) (visa types.Visa, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VisaKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return visa, false
	}
	k.cdc.MustUnmarshal(bz, &visa)
	return visa, true
}

// GetAllVisa returns all visas
func (k Keeper) GetAllVisa(ctx sdk.Context) (list []*types.Visa) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VisaKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var visa types.Visa
		k.cdc.MustUnmarshal(iterator.Value(), &visa)
		list = append(list, &visa)
	}

	return
}

// GetVisaByJamaah returns visas of a jamaah
func (k Keeper) GetVisaByJamaah(ctx sdk.Context, jamaah string) (list []*types.Visa) {
	for _, visa := range k.GetAllVisa(ctx) {
		if visa.Jamaah == jamaah {
			list = append(list, visa)
		}
	}
	return
}

// RemoveVisa removes a visa from the store
func (k Keeper) RemoveVisa(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VisaKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
