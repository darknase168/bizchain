package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/pembayaran/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the pembayaran module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListPembayaran(),
		GetCmdShowPembayaran(),
		GetCmdPembayaranByJamaah(),
	)

	return queryCmd
}

// GetCmdListPembayaran queries all payments
func GetCmdListPembayaran() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-pembayaran",
		Short: "List all payments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.PembayaranAll(cmd.Context(), &types.QueryAllPembayaranRequest{})
			if err != nil {
				return err
			}

			for _, p := range resp.Pembayaran {
				fmt.Printf("ID: %d | Jamaah: %s | Paket: %d | Total: %s | Paid: %s | Remaining: %s | Status: %s\n",
					p.Id, p.Jamaah, p.PaketId, p.Total, p.Paid, p.Remaining, p.Status)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowPembayaran queries a specific payment
func GetCmdShowPembayaran() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-pembayaran [id]",
		Short: "Show a payment by ID",
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
			resp, err := queryClient.Pembayaran(cmd.Context(), &types.QueryGetPembayaranRequest{Id: id})
			if err != nil {
				return err
			}

			p := resp.Pembayaran
			fmt.Printf("ID: %d\n", p.Id)
			fmt.Printf("Jamaah: %s\n", p.Jamaah)
			fmt.Printf("Paket ID: %d\n", p.PaketId)
			fmt.Printf("Total: %s\n", p.Total)
			fmt.Printf("Down Payment: %s\n", p.DownPayment)
			fmt.Printf("Paid: %s\n", p.Paid)
			fmt.Printf("Remaining: %s\n", p.Remaining)
			fmt.Printf("Status: %s\n", p.Status)
			fmt.Printf("Payment Method: %s\n", p.PaymentMethod)
			if len(p.Installments) > 0 {
				fmt.Println("Installments:")
				for _, inst := range p.Installments {
					fmt.Printf("  - #%d | %s | due %s | paid: %t | %s\n", inst.Id, inst.Amount, inst.DueDate, inst.Paid, inst.PaidAt)
				}
			}
			if len(p.EscrowStages) > 0 {
				fmt.Println("Escrow Stages:")
				for _, stage := range p.EscrowStages {
					fmt.Printf("  - %s | %s | released: %t | %s\n", stage.Name, stage.Amount, stage.Released, stage.ReleasedAt)
				}
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdPembayaranByJamaah queries payments of a jamaah
func GetCmdPembayaranByJamaah() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pembayaran-by-jamaah [jamaah]",
		Short: "List payments of a jamaah",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.PembayaranByJamaah(cmd.Context(), &types.QueryByJamaahRequest{Jamaah: args[0]})
			if err != nil {
				return err
			}

			for _, p := range resp.Pembayaran {
				fmt.Printf("ID: %d | Paket: %d | Total: %s | Paid: %s | Remaining: %s | Status: %s\n",
					p.Id, p.PaketId, p.Total, p.Paid, p.Remaining, p.Status)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
