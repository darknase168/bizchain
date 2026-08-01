package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/referral/types"
)

// Keeper of the referral module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new referral Keeper instance
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

// GetNextReferralID returns the next referral ID
func (k Keeper) GetNextReferralID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ReferralCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetReferralCount sets the referral ID counter
func (k Keeper) SetReferralCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.ReferralCountKey, bz)
}

// SetReferral stores a referral
func (k Keeper) SetReferral(ctx sdk.Context, referral types.Referral) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ReferralKey)
	bz := k.cdc.MustMarshal(&referral)
	store.Set(sdk.Uint64ToBigEndian(referral.Id), bz)

	currentCount := k.GetNextReferralID(ctx)
	if referral.Id >= currentCount {
		k.SetReferralCount(ctx, referral.Id)
	}
}

// GetReferral retrieves a referral by ID
func (k Keeper) GetReferral(ctx sdk.Context, id uint64) (referral types.Referral, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ReferralKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return referral, false
	}
	k.cdc.MustUnmarshal(bz, &referral)
	return referral, true
}

// GetAllReferral returns all referrals
func (k Keeper) GetAllReferral(ctx sdk.Context) (list []*types.Referral) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ReferralKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var referral types.Referral
		k.cdc.MustUnmarshal(iterator.Value(), &referral)
		list = append(list, &referral)
	}

	return
}

// GetReferralByAgent returns referrals of an agent
func (k Keeper) GetReferralByAgent(ctx sdk.Context, agent string) (list []*types.Referral) {
	for _, referral := range k.GetAllReferral(ctx) {
		if referral.Agent == agent {
			list = append(list, referral)
		}
	}
	return
}

// RemoveReferral removes a referral from the store
func (k Keeper) RemoveReferral(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ReferralKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
