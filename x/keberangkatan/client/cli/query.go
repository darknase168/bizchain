package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/keberangkatan/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the keberangkatan module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListKeberangkatan(),
		GetCmdShowKeberangkatan(),
		GetCmdKeberangkatanByJamaah(),
	)

	return queryCmd
}

// GetCmdListKeberangkatan queries all journeys
func GetCmdListKeberangkatan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-keberangkatan",
		Short: "List all departure journeys",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.KeberangkatanAll(cmd.Context(), &types.QueryAllKeberangkatanRequest{})
			if err != nil {
				return err
			}

			for _, k := range resp.Keberangkatan {
				fmt.Printf("ID: %d | Jamaah: %s | Paket: %d | Stage: %d | %s\n",
					k.Id, k.Jamaah, k.PaketId, k.Stage, k.StatusLabel)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowKeberangkatan queries a specific journey
func GetCmdShowKeberangkatan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-keberangkatan [id]",
		Short: "Show a journey by ID",
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
			resp, err := queryClient.Keberangkatan(cmd.Context(), &types.QueryGetKeberangkatanRequest{Id: id})
			if err != nil {
				return err
			}

			k := resp.Keberangkatan
			fmt.Printf("ID: %d\n", k.Id)
			fmt.Printf("Jamaah: %s\n", k.Jamaah)
			fmt.Printf("Paket ID: %d\n", k.PaketId)
			fmt.Printf("Pembayaran ID: %d\n", k.PembayaranId)
			fmt.Printf("Stage: %d - %s\n", k.Stage, k.StatusLabel)
			fmt.Printf("Departure: %s | Return: %s | Manasik: %s\n", k.DepartureDate, k.ReturnDate, k.ManasikDate)
			fmt.Printf("Hotel Confirm: %s\n", k.HotelConfirm)
			fmt.Printf("Airline Confirm: %s\n", k.AirlineConfirm)
			if len(k.Baggage) > 0 {
				fmt.Println("Baggage:")
				for _, b := range k.Baggage {
					fmt.Printf("  - #%d | tag: %s | weight: %s | status: %s\n", b.Id, b.Tag, b.Weight, b.Status)
				}
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdKeberangkatanByJamaah queries journeys of a pilgrim
func GetCmdKeberangkatanByJamaah() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keberangkatan-by-jamaah [jamaah]",
		Short: "List journeys of a pilgrim",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.KeberangkatanByJamaah(cmd.Context(), &types.QueryKeberangkatanByJamaahRequest{Jamaah: args[0]})
			if err != nil {
				return err
			}

			for _, k := range resp.Keberangkatan {
				fmt.Printf("ID: %d | Paket: %d | Stage: %d | %s\n",
					k.Id, k.PaketId, k.Stage, k.StatusLabel)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
