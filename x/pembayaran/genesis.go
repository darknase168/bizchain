package pembayaran

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/pembayaran/keeper"
	"github.com/bizchain/blockchain/x/pembayaran/types"
)

// InitGenesis initializes the pembayaran module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetPembayaranCount(ctx, genState.PembayaranCount)
	k.SetInstallmentCount(ctx, genState.InstallmentCount)

	for _, pembayaran := range genState.PembayaranList {
		k.SetPembayaran(ctx, *pembayaran)
	}
}

// ExportGenesis exports the pembayaran module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.PembayaranList = k.GetAllPembayaran(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.PembayaranList) > 0 {
		genesis.PembayaranCount = genesis.PembayaranList[len(genesis.PembayaranList)-1].Id
	}
	for _, pembayaran := range genesis.PembayaranList {
		for _, inst := range pembayaran.Installments {
			if inst.Id > genesis.InstallmentCount {
				genesis.InstallmentCount = inst.Id
			}
		}
	}

	return genesis
}
