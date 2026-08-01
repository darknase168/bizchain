package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/paket/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the paket module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListPaket(),
		GetCmdShowPaket(),
	)

	return queryCmd
}

// GetCmdListPaket queries all packages (marketplace)
func GetCmdListPaket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-paket",
		Short: "List all packages (marketplace)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.PaketAll(cmd.Context(), &types.QueryAllPaketRequest{})
			if err != nil {
				return err
			}

			for _, p := range resp.Paket {
				fmt.Printf("ID: %d | %s | %s | Quota: %d/%d | Status: %s | Category: %s\n",
					p.Id, p.Name, p.Price, p.Booked, p.Quota, p.Status, p.Category)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowPaket queries a specific package
func GetCmdShowPaket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-paket [id]",
		Short: "Show a package by ID",
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
			resp, err := queryClient.Paket(cmd.Context(), &types.QueryGetPaketRequest{Id: id})
			if err != nil {
				return err
			}

			p := resp.Paket
			fmt.Printf("ID: %d\n", p.Id)
			fmt.Printf("Name: %s\n", p.Name)
			fmt.Printf("Price: %s\n", p.Price)
			fmt.Printf("Schedule: %s\n", p.Schedule)
			fmt.Printf("Quota: %d (booked %d)\n", p.Quota, p.Booked)
			fmt.Printf("Hotel: %s\n", p.Hotel)
			fmt.Printf("Airline: %s\n", p.Airline)
			fmt.Printf("Muthawif: %s\n", p.Muthawif)
			fmt.Printf("Status: %s\n", p.Status)
			fmt.Printf("Departure: %s | Return: %s\n", p.DepartureDate, p.ReturnDate)
			fmt.Printf("Category: %s\n", p.Category)
			fmt.Printf("Creator: %s\n", p.Creator)
			if len(p.Reviews) > 0 {
				fmt.Println("Reviews:")
				for _, r := range p.Reviews {
					fmt.Printf("  - %d stars by %s: %s\n", r.Rating, r.Reviewer, r.Comment)
				}
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
