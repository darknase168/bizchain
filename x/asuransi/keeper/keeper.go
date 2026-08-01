package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/asuransi/types"
)

// Keeper of the asuransi module
type Keeper struct {
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	cdc      codec.BinaryCodec
	authority string
}

// NewKeeper creates a new asuransi Keeper instance
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

// GetNextAsuransiID returns the next asuransi ID
func (k Keeper) GetNextAsuransiID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AsuransiCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetAsuransiCount sets the asuransi ID counter
func (k Keeper) SetAsuransiCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.AsuransiCountKey, bz)
}

// SetAsuransi stores an asuransi policy
func (k Keeper) SetAsuransi(ctx sdk.Context, asuransi types.Asuransi) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AsuransiKey)
	bz := k.cdc.MustMarshal(&asuransi)
	store.Set(sdk.Uint64ToBigEndian(asuransi.Id), bz)

	currentCount := k.GetNextAsuransiID(ctx)
	if asuransi.Id >= currentCount {
		k.SetAsuransiCount(ctx, asuransi.Id)
	}
}

// GetAsuransi retrieves an asuransi policy by ID
func (k Keeper) GetAsuransi(ctx sdk.Context, id uint64) (asuransi types.Asuransi, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AsuransiKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return asuransi, false
	}
	k.cdc.MustUnmarshal(bz, &asuransi)
	return asuransi, true
}

// GetAllAsuransi returns all asuransi policies
func (k Keeper) GetAllAsuransi(ctx sdk.Context) (list []*types.Asuransi) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AsuransiKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var asuransi types.Asuransi
		k.cdc.MustUnmarshal(iterator.Value(), &asuransi)
		list = append(list, &asuransi)
	}

	return
}

// GetNextClaimID returns the next claim ID
func (k Keeper) GetNextClaimID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ClaimCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetClaimCount sets the claim ID counter
func (k Keeper) SetClaimCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.ClaimCountKey, bz)
}

// SetClaim stores a claim
func (k Keeper) SetClaim(ctx sdk.Context, claim types.AsuransiClaim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClaimKey)
	bz := k.cdc.MustMarshal(&claim)
	store.Set(sdk.Uint64ToBigEndian(claim.Id), bz)

	currentCount := k.GetNextClaimID(ctx)
	if claim.Id >= currentCount {
		k.SetClaimCount(ctx, claim.Id)
	}
}

// GetClaim retrieves a claim by ID
func (k Keeper) GetClaim(ctx sdk.Context, id uint64) (claim types.AsuransiClaim, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClaimKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return claim, false
	}
	k.cdc.MustUnmarshal(bz, &claim)
	return claim, true
}

// GetAllClaims returns all claims
func (k Keeper) GetAllClaims(ctx sdk.Context) (list []*types.AsuransiClaim) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClaimKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var claim types.AsuransiClaim
		k.cdc.MustUnmarshal(iterator.Value(), &claim)
		list = append(list, &claim)
	}

	return
}
