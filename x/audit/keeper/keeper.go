package keeper

import (
	"encoding/binary"

	cosmossdklog "cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/audit/types"
)

// Keeper of the audit module
type Keeper struct {
	storeKey  storetypes.StoreKey
	memKey    storetypes.StoreKey
	cdc       codec.BinaryCodec
	authority string
}

// NewKeeper creates a new audit Keeper instance
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

// GetNextAuditLogID returns the next audit log ID
func (k Keeper) GetNextAuditLogID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AuditCountKey)
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz) + 1
}

// SetAuditCount sets the audit log ID counter
func (k Keeper) SetAuditCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.AuditCountKey, bz)
}

// SetAuditLog stores an audit log
func (k Keeper) SetAuditLog(ctx sdk.Context, auditLog types.AuditLog) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditLogKey)
	bz := k.cdc.MustMarshal(&auditLog)
	store.Set(sdk.Uint64ToBigEndian(auditLog.Id), bz)

	currentCount := k.GetNextAuditLogID(ctx)
	if auditLog.Id >= currentCount {
		k.SetAuditCount(ctx, auditLog.Id)
	}
}

// GetAuditLog retrieves an audit log by ID
func (k Keeper) GetAuditLog(ctx sdk.Context, id uint64) (auditLog types.AuditLog, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditLogKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return auditLog, false
	}
	k.cdc.MustUnmarshal(bz, &auditLog)
	return auditLog, true
}

// GetAllAuditLog returns all audit logs
func (k Keeper) GetAllAuditLog(ctx sdk.Context) (list []*types.AuditLog) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditLogKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var auditLog types.AuditLog
		k.cdc.MustUnmarshal(iterator.Value(), &auditLog)
		list = append(list, &auditLog)
	}

	return
}

// GetAuditLogByModule returns audit logs filtered by module
func (k Keeper) GetAuditLogByModule(ctx sdk.Context, module string) (list []*types.AuditLog) {
	for _, auditLog := range k.GetAllAuditLog(ctx) {
		if auditLog.Module == module {
			list = append(list, auditLog)
		}
	}
	return
}

// RemoveAuditLog removes an audit log from the store
func (k Keeper) RemoveAuditLog(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditLogKey)
	store.Delete(sdk.Uint64ToBigEndian(id))
}
