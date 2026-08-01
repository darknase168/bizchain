package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/ticket/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "ticket module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdIssueTicket(),
		GetCmdUpdateTicketStatus(),
	)

	return txCmd
}

// GetCmdIssueTicket issues an NFT ticket
func GetCmdIssueTicket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-ticket [jamaah] [paket-id] [airline] [flight-number]",
		Short: "Issue an NFT flight ticket",
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

			seat, _ := cmd.Flags().GetString("seat")
			schedule, _ := cmd.Flags().GetString("schedule")
			qrCode, _ := cmd.Flags().GetString("qr-code")
			docHash, _ := cmd.Flags().GetString("document-hash")

			msg := &types.MsgIssueTicket{
				Creator:      clientCtx.GetFromAddress().String(),
				Jamaah:       args[0],
				PaketId:      paketID,
				Airline:      args[2],
				FlightNumber: args[3],
				Seat:         seat,
				Schedule:     schedule,
				QrCode:       qrCode,
				DocumentHash: docHash,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("seat", "", "Seat number")
	cmd.Flags().String("schedule", "", "Flight schedule")
	cmd.Flags().String("qr-code", "", "QR code payload")
	cmd.Flags().String("document-hash", "", "Document hash")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateTicketStatus updates ticket status
func GetCmdUpdateTicketStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-ticket-status [id] [status]",
		Short: "Update ticket status (issued, checked_in, boarded, used, void)",
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

			msg := &types.MsgUpdateTicketStatus{
				Creator: clientCtx.GetFromAddress().String(),
				Id:      id,
				Status:  args[1],
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
