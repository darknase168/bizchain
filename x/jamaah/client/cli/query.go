package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/jamaah/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the jamaah module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListJamaah(),
		GetCmdShowJamaah(),
	)

	return queryCmd
}

// GetCmdListJamaah queries all pilgrims
func GetCmdListJamaah() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-jamaah",
		Short: "List all pilgrims",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.JamaahAll(cmd.Context(), &types.QueryAllJamaahRequest{})
			if err != nil {
				return err
			}

			for _, j := range resp.Jamaah {
				fmt.Printf("ID: %d | Name: %s | Passport: %s | Status: %s | DID: %s\n",
					j.Id, j.Name, j.PassportNumber, j.Status, j.Did)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowJamaah queries a specific pilgrim
func GetCmdShowJamaah() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-jamaah [id]",
		Short: "Show a pilgrim by ID",
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
			resp, err := queryClient.Jamaah(cmd.Context(), &types.QueryGetJamaahRequest{Id: id})
			if err != nil {
				return err
			}

			j := resp.Jamaah
			fmt.Printf("ID: %d\n", j.Id)
			fmt.Printf("Name: %s\n", j.Name)
			fmt.Printf("Phone: %s\n", j.Phone)
			fmt.Printf("Email: %s\n", j.Email)
			fmt.Printf("Address: %s\n", j.Address)
			fmt.Printf("Passport: %s\n", j.PassportNumber)
			fmt.Printf("Status: %s\n", j.Status)
			fmt.Printf("DID: %s\n", j.Did)
			fmt.Printf("Creator: %s\n", j.Creator)
			if len(j.Documents) > 0 {
				fmt.Println("Documents:")
				for _, d := range j.Documents {
					fmt.Printf("  - %s | hash: %s | %s\n", d.DocType, d.Hash, d.UploadedAt)
				}
			}
			if len(j.Vaccinations) > 0 {
				fmt.Println("Vaccinations:")
				for _, v := range j.Vaccinations {
					fmt.Printf("  - %s | %s | issuer: %s\n", v.VaccineName, v.Date, v.Issuer)
				}
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
