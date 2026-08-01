package audit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/audit/keeper"
	"github.com/bizchain/blockchain/x/audit/types"
)

// InitGenesis initializes the audit module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetAuditCount(ctx, genState.AuditCount)

	for _, auditLog := range genState.AuditLogList {
		k.SetAuditLog(ctx, *auditLog)
	}
}

// ExportGenesis exports the audit module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.AuditLogList = k.GetAllAuditLog(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.AuditLogList) > 0 {
		genesis.AuditCount = genesis.AuditLogList[len(genesis.AuditLogList)-1].Id
	}

	return genesis
}
