package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/agen/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "agen module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateAgen(),
		GetCmdUpdateAgen(),
		GetCmdAddComplaint(),
		GetCmdResolveComplaint(),
		GetCmdRecordPerformance(),
	)

	return txCmd
}

// GetCmdCreateAgen registers a new agent
func GetCmdCreateAgen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-agen [name] [address]",
		Short: "Register a new agent (pusat/cabang/subagen)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			parentID, _ := cmd.Flags().GetString("parent-id")
			level, _ := cmd.Flags().GetString("level")
			commissionRate, _ := cmd.Flags().GetString("commission-rate")

			msg := &types.MsgCreateAgen{
				Creator:        clientCtx.GetFromAddress().String(),
				Name:           args[0],
				Address:        args[1],
				ParentId:       parentID,
				Level:          level,
				CommissionRate: commissionRate,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("parent-id", "", "Parent agent ID (0 for pusat)")
	cmd.Flags().String("level", "subagen", "Agent level (pusat, cabang, subagen)")
	cmd.Flags().String("commission-rate", "", "Commission rate in basis points")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateAgen updates agent data
func GetCmdUpdateAgen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-agen [id]",
		Short: "Update agent data (level, parent, commission, status)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			parentID, _ := cmd.Flags().GetString("parent-id")
			level, _ := cmd.Flags().GetString("level")
			commissionRate, _ := cmd.Flags().GetString("commission-rate")
			status, _ := cmd.Flags().GetString("status")

			msg := &types.MsgUpdateAgen{
				Creator:        clientCtx.GetFromAddress().String(),
				Id:             id,
				Name:           name,
				ParentId:       parentID,
				Level:          level,
				CommissionRate: commissionRate,
				Status:         status,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("name", "", "Agent name")
	cmd.Flags().String("parent-id", "", "New parent agent ID")
	cmd.Flags().String("level", "", "Agent level (pusat, cabang, subagen)")
	cmd.Flags().String("commission-rate", "", "Commission rate in basis points")
	cmd.Flags().String("status", "", "Status (active, inactive, suspended)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAddComplaint files a complaint against an agent
func GetCmdAddComplaint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-complaint [agen-id] [reason]",
		Short: "File a complaint against an agent",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			agenID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgAddComplaint{
				Creator: clientCtx.GetFromAddress().String(),
				AgenId:  agenID,
				Reason:  args[1],
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

// GetCmdResolveComplaint resolves a complaint
func GetCmdResolveComplaint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve-complaint [agen-id] [complaint-id] [resolution]",
		Short: "Resolve a complaint against an agent",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			agenID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			complaintID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgResolveComplaint{
				Creator:     clientCtx.GetFromAddress().String(),
				AgenId:      agenID,
				ComplaintId: complaintID,
				Resolution:  args[2],
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

// GetCmdRecordPerformance records agent performance
func GetCmdRecordPerformance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-performance [agen-id] [period] [sales] [volume]",
		Short: "Record agent performance for a period (rekam jejak)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			agenID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			sales, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			ratingAvg, _ := cmd.Flags().GetString("rating-avg")

			msg := &types.MsgRecordPerformance{
				Creator:   clientCtx.GetFromAddress().String(),
				AgenId:    agenID,
				Period:    args[1],
				Sales:     sales,
				Volume:    args[3],
				RatingAvg: ratingAvg,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("rating-avg", "", "Average rating for the period (0-5)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
