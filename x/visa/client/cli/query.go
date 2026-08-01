package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/visa/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the visa module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListVisa(),
		GetCmdShowVisa(),
	)

	return queryCmd
}

// GetCmdListVisa queries all visas
func GetCmdListVisa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-visa",
		Short: "List all visas",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.VisaAll(cmd.Context(), &types.QueryAllVisaRequest{})
			if err != nil {
				return err
			}

			for _, v := range resp.Visa {
				fmt.Printf("ID: %d | Jamaah: %s | Paket: %d | Status: %s | No: %s\n",
					v.Id, v.Jamaah, v.PaketId, v.Status, v.VisaNumber)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowVisa queries a specific visa
func GetCmdShowVisa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-visa [id]",
		Short: "Show a visa by ID",
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
			resp, err := queryClient.Visa(cmd.Context(), &types.QueryGetVisaRequest{Id: id})
			if err != nil {
				return err
			}

			v := resp.Visa
			fmt.Printf("ID: %d\n", v.Id)
			fmt.Printf("Jamaah: %s\n", v.Jamaah)
			fmt.Printf("Paket ID: %d\n", v.PaketId)
			fmt.Printf("Status: %s\n", v.Status)
			fmt.Printf("Visa Number: %s\n", v.VisaNumber)
			fmt.Printf("Issue Date: %s\n", v.IssueDate)
			fmt.Printf("Expiry Date: %s\n", v.ExpiryDate)
			fmt.Printf("Document Hash: %s\n", v.DocumentHash)
			fmt.Printf("Notes: %s\n", v.Notes)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
