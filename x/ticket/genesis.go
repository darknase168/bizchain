package ticket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bizchain/blockchain/x/ticket/keeper"
	"github.com/bizchain/blockchain/x/ticket/types"
)

// InitGenesis initializes the ticket module genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetTicketCount(ctx, genState.TicketCount)

	for _, ticket := range genState.TicketList {
		k.SetTicket(ctx, *ticket)
	}
}

// ExportGenesis exports the ticket module genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.TicketList = k.GetAllTicket(ctx)
	genesis.Params = types.DefaultParams()

	if len(genesis.TicketList) > 0 {
		genesis.TicketCount = genesis.TicketList[len(genesis.TicketList)-1].Id
	}

	return genesis
}
