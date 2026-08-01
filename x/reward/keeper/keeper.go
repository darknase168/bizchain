package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/reward/types"
)

// Keeper of the reward module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new reward Keeper instance
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

// GetNextRewardID returns the next reward ID
func (k Keeper) GetNextRewardID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetRewardCount sets the reward ID counter
func (k Keeper) SetRewardCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.RewardCountKey, bz)
}

// SetReward stores a reward
func (k Keeper) SetReward(ctx sdk.Context, reward types.Reward) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RewardKey)
	bz := k.cdc.MustMarshal(&reward)
	store.Set(sdk.Uint64ToBigEndian(reward.Id), bz)

	currentCount := k.GetNextRewardID(ctx)
	if reward.Id >= currentCount {
		k.SetRewardCount(ctx, reward.Id)
	}
}

// GetReward retrieves a reward by ID
func (k Keeper) GetReward(ctx sdk.Context, id uint64) (reward types.Reward, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RewardKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return reward, false
	}
	k.cdc.MustUnmarshal(bz, &reward)
	return reward, true
}

// GetAllReward returns all rewards
func (k Keeper) GetAllReward(ctx sdk.Context) (list []*types.Reward) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RewardKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var reward types.Reward
		k.cdc.MustUnmarshal(iterator.Value(), &reward)
		list = append(list, &reward)
	}

	return
}

// SetRewardBalance stores a jamaah's loyalty balance
func (k Keeper) SetRewardBalance(ctx sdk.Context, balance types.RewardBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalanceKey)
	bz := k.cdc.MustMarshal(&balance)
	store.Set([]byte(balance.Jamaah), bz)
}

// GetRewardBalance retrieves a jamaah's loyalty balance
func (k Keeper) GetRewardBalance(ctx sdk.Context, jamaah string) (balance types.RewardBalance, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalanceKey)
	bz := store.Get([]byte(jamaah))
	if bz == nil {
		return balance, false
	}
	k.cdc.MustUnmarshal(bz, &balance)
	return balance, true
}

// GetAllRewardBalance returns all loyalty balances
func (k Keeper) GetAllRewardBalance(ctx sdk.Context) (list []*types.RewardBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalanceKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var balance types.RewardBalance
		k.cdc.MustUnmarshal(iterator.Value(), &balance)
		list = append(list, &balance)
	}

	return
}

// RemoveReward removes a reward from the store
func (k Keeper) RemoveReward(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RewardKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
