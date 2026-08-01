package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/asuransi/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "asuransi module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateAsuransi(),
		GetCmdUpdateAsuransiStatus(),
		GetCmdSubmitClaim(),
		GetCmdDecideClaim(),
	)

	return txCmd
}

// GetCmdCreateAsuransi issues a digital insurance policy
func GetCmdCreateAsuransi() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [jamaah] [policy-type] [premium]",
		Short: "Issue a digital insurance policy",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coverage, _ := cmd.Flags().GetString("coverage")
			startDate, _ := cmd.Flags().GetString("start-date")
			endDate, _ := cmd.Flags().GetString("end-date")
			documentHash, _ := cmd.Flags().GetString("document-hash")
			provider, _ := cmd.Flags().GetString("provider")

			msg := &types.MsgCreateAsuransi{
				Creator:      clientCtx.GetFromAddress().String(),
				Jamaah:       args[0],
				PolicyType:   args[1],
				Premium:      args[2],
				Coverage:     coverage,
				StartDate:    startDate,
				EndDate:      endDate,
				DocumentHash: documentHash,
				Provider:     provider,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("coverage", "", "Max coverage amount")
	cmd.Flags().String("start-date", "", "Policy start date")
	cmd.Flags().String("end-date", "", "Policy end date")
	cmd.Flags().String("document-hash", "", "Policy document hash")
	cmd.Flags().String("provider", "", "Insurance provider")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateAsuransiStatus updates policy status
func GetCmdUpdateAsuransiStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-status [id] [status]",
		Short: "Update policy status (active, expired, cancelled, claimed)",
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

			msg := &types.MsgUpdateAsuransiStatus{
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

// GetCmdSubmitClaim files an insurance claim
func GetCmdSubmitClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-claim [asuransi-id] [reason] [amount]",
		Short: "Submit an insurance claim",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asuransiID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			evidenceHash, _ := cmd.Flags().GetString("evidence-hash")

			msg := &types.MsgSubmitClaim{
				Creator:      clientCtx.GetFromAddress().String(),
				AsuransiId:   asuransiID,
				Reason:       args[1],
				Amount:       args[2],
				EvidenceHash: evidenceHash,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("evidence-hash", "", "Evidence document hash")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDecideClaim approves/rejects a claim
func GetCmdDecideClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decide-claim [claim-id] [status]",
		Short: "Approve or reject a claim (approved, rejected)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			claimID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			note, _ := cmd.Flags().GetString("decision-note")

			msg := &types.MsgDecideClaim{
				Creator: clientCtx.GetFromAddress().String(),
				ClaimId: claimID,
				Status:  args[1],
				Note:    note,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("decision-note", "", "Decision note")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
