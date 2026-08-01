package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/visa/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "visa module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateVisa(),
		GetCmdUpdateVisaStatus(),
	)

	return txCmd
}

// GetCmdCreateVisa starts a visa application
func GetCmdCreateVisa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-visa [jamaah] [paket-id]",
		Short: "Start a visa application",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			paketID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			docHash, _ := cmd.Flags().GetString("document-hash")
			notes, _ := cmd.Flags().GetString("notes")

			msg := &types.MsgCreateVisa{
				Creator:      clientCtx.GetFromAddress().String(),
				Jamaah:       args[0],
				PaketId:      paketID,
				DocumentHash: docHash,
				Notes:        notes,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("document-hash", "", "Hash of the passport document used")
	cmd.Flags().String("notes", "", "Application notes")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateVisaStatus updates visa processing status
func GetCmdUpdateVisaStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-visa-status [id] [status]",
		Short: "Update visa status (processing, issued, rejected, expired)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			visaNumber, _ := cmd.Flags().GetString("visa-number")
			issueDate, _ := cmd.Flags().GetString("issue-date")
			expiryDate, _ := cmd.Flags().GetString("expiry-date")

			msg := &types.MsgUpdateVisaStatus{
				Creator:    clientCtx.GetFromAddress().String(),
				Id:         id,
				Status:     args[1],
				VisaNumber: visaNumber,
				IssueDate:  issueDate,
				ExpiryDate: expiryDate,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("visa-number", "", "Visa number")
	cmd.Flags().String("issue-date", "", "Issue date")
	cmd.Flags().String("expiry-date", "", "Expiry date")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
