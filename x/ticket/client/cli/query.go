package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/ticket/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the ticket module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListTicket(),
		GetCmdShowTicket(),
	)

	return queryCmd
}

// GetCmdListTicket queries all tickets
func GetCmdListTicket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-ticket",
		Short: "List all tickets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.TicketAll(cmd.Context(), &types.QueryAllTicketRequest{})
			if err != nil {
				return err
			}

			for _, t := range resp.Ticket {
				fmt.Printf("ID: %d | Jamaah: %s | Flight: %s %s | Seat: %s | Status: %s | NFT: %s\n",
					t.Id, t.Jamaah, t.Airline, t.FlightNumber, t.Seat, t.Status, t.NftId)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowTicket queries a specific ticket
func GetCmdShowTicket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-ticket [id]",
		Short: "Show a ticket by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Ticket(cmd.Context(), &types.QueryGetTicketRequest{Id: id})
			if err != nil {
				return err
			}

			t := resp.Ticket
			fmt.Printf("ID: %d\n", t.Id)
			fmt.Printf("Jamaah: %s\n", t.Jamaah)
			fmt.Printf("Paket ID: %d\n", t.PaketId)
			fmt.Printf("Airline: %s\n", t.Airline)
			fmt.Printf("Flight: %s\n", t.FlightNumber)
			fmt.Printf("Seat: %s\n", t.Seat)
			fmt.Printf("Schedule: %s\n", t.Schedule)
			fmt.Printf("Status: %s\n", t.Status)
			fmt.Printf("NFT ID: %s\n", t.NftId)
			fmt.Printf("QR: %s\n", t.QrCode)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
