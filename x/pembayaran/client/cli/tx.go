package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/pembayaran/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "pembayaran module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreatePembayaran(),
		GetCmdPayInstallment(),
		GetCmdReleaseEscrow(),
		GetCmdRefundPembayaran(),
		GetCmdCancelPembayaran(),
	)

	return txCmd
}

// GetCmdCreatePembayaran creates a new escrow payment
func GetCmdCreatePembayaran() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pembayaran [jamaah] [paket-id] [total] [down-payment]",
		Short: "Create a new escrow payment (DP)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			paketID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			installmentCount, _ := cmd.Flags().GetUint32("installments")
			paymentMethod, _ := cmd.Flags().GetString("payment-method")
			dueDates, _ := cmd.Flags().GetStringSlice("due-dates")
			stagesStr, _ := cmd.Flags().GetString("stages")

			escrowStages, err := parseEscrowStages(stagesStr)
			if err != nil {
				return err
			}

			msg := &types.MsgCreatePembayaran{
				Creator:             clientCtx.GetFromAddress().String(),
				Jamaah:              args[0],
				PaketId:             paketID,
				Total:               args[2],
				DownPayment:         args[3],
				PaymentMethod:       paymentMethod,
				InstallmentCount:    installmentCount,
				InstallmentDueDates: dueDates,
				EscrowStages:        escrowStages,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint32("installments", 0, "Number of installments (cicilan)")
	cmd.Flags().String("payment-method", "transfer", "Payment method (transfer, qris, cash)")
	cmd.Flags().StringSlice("due-dates", nil, "Installment due dates, e.g. --due-dates 2026-02-01 --due-dates 2026-03-01")
	cmd.Flags().String("stages", "", "Escrow stages as 'name:amount,name:amount', e.g. 'visa:2000000,tiket:3000000,hotel:1000000'")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// parseEscrowStages parses "name:amount,name:amount" into EscrowStages
func parseEscrowStages(s string) ([]*types.EscrowStage, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	stages := make([]*types.EscrowStage, 0, len(parts))
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid escrow stage %q, expected name:amount", part)
		}
		stages = append(stages, &types.EscrowStage{
			Name:   kv[0],
			Amount: kv[1],
		})
	}
	return stages, nil
}

// GetCmdPayInstallment records an installment payment
func GetCmdPayInstallment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay-installment [pembayaran-id] [installment-id] [amount]",
		Short: "Pay an installment (cicilan)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pembayaranID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			installmentID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgPayInstallment{
				Creator:       clientCtx.GetFromAddress().String(),
				PembayaranId:  pembayaranID,
				InstallmentId: installmentID,
				Amount:        args[2],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdReleaseEscrow releases escrow funds for a stage
func GetCmdReleaseEscrow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release-escrow [pembayaran-id] [stage-name]",
		Short: "Release escrow funds for a stage (visa, tiket, hotel)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pembayaranID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgReleaseEscrow{
				Creator:      clientCtx.GetFromAddress().String(),
				PembayaranId: pembayaranID,
				StageName:    args[1],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdRefundPembayaran refunds an escrow payment
func GetCmdRefundPembayaran() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-pembayaran [pembayaran-id]",
		Short: "Refund an escrow payment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pembayaranID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			reason, _ := cmd.Flags().GetString("reason")

			msg := &types.MsgRefundPembayaran{
				Creator:      clientCtx.GetFromAddress().String(),
				PembayaranId: pembayaranID,
				Reason:       reason,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("reason", "", "Refund reason")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCancelPembayaran cancels a pending payment
func GetCmdCancelPembayaran() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-pembayaran [pembayaran-id]",
		Short: "Cancel a pending payment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pembayaranID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCancelPembayaran{
				Creator:      clientCtx.GetFromAddress().String(),
				PembayaranId: pembayaranID,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
