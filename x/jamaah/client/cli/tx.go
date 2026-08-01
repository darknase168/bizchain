package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/jamaah/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "jamaah module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateJamaah(),
		GetCmdUpdateJamaah(),
		GetCmdAddDocument(),
		GetCmdAddVaccination(),
	)

	return txCmd
}

// GetCmdCreateJamaah registers a new pilgrim
func GetCmdCreateJamaah() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-jamaah [name] [passport-number]",
		Short: "Register a new pilgrim (self-sovereign DID)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			phone, _ := cmd.Flags().GetString("phone")
			email, _ := cmd.Flags().GetString("email")
			address, _ := cmd.Flags().GetString("address")
			photoHash, _ := cmd.Flags().GetString("photo-hash")

			msg := &types.MsgCreateJamaah{
				Creator:        clientCtx.GetFromAddress().String(),
				Name:           args[0],
				PassportNumber: args[1],
				Phone:          phone,
				Email:          email,
				Address:        address,
				PhotoHash:      photoHash,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("phone", "", "Phone number")
	cmd.Flags().String("email", "", "Email address")
	cmd.Flags().String("address", "", "Home address")
	cmd.Flags().String("photo-hash", "", "Photo hash (IPFS/object storage)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateJamaah updates pilgrim data
func GetCmdUpdateJamaah() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-jamaah [id] [name] [passport-number]",
		Short: "Update pilgrim data",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			phone, _ := cmd.Flags().GetString("phone")
			email, _ := cmd.Flags().GetString("email")
			address, _ := cmd.Flags().GetString("address")
			photoHash, _ := cmd.Flags().GetString("photo-hash")
			status, _ := cmd.Flags().GetString("status")

			msg := &types.MsgUpdateJamaah{
				Creator:        clientCtx.GetFromAddress().String(),
				Id:             id,
				Name:           args[1],
				PassportNumber: args[2],
				Phone:          phone,
				Email:          email,
				Address:        address,
				PhotoHash:      photoHash,
				Status:         status,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("phone", "", "Phone number")
	cmd.Flags().String("email", "", "Email address")
	cmd.Flags().String("address", "", "Home address")
	cmd.Flags().String("photo-hash", "", "Photo hash")
	cmd.Flags().String("status", "", "Status (active, inactive, blocked)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAddDocument adds a hashed document to a pilgrim
func GetCmdAddDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-document [jamaah-id] [doc-type] [hash]",
		Short: "Add a hashed document (passport, visa, vaccine, health)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			jamaahID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			storageRef, _ := cmd.Flags().GetString("storage-ref")

			msg := &types.MsgAddDocument{
				Creator:    clientCtx.GetFromAddress().String(),
				JamaahId:   jamaahID,
				DocType:    args[1],
				Hash:       args[2],
				StorageRef: storageRef,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("storage-ref", "", "Off-chain storage reference (IPFS)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAddVaccination adds a vaccination record
func GetCmdAddVaccination() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-vaccination [jamaah-id] [vaccine-name] [date]",
		Short: "Add a vaccination record",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			jamaahID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			issuer, _ := cmd.Flags().GetString("issuer")
			batch, _ := cmd.Flags().GetString("batch")

			msg := &types.MsgAddVaccination{
				Creator:     clientCtx.GetFromAddress().String(),
				JamaahId:    jamaahID,
				VaccineName: args[1],
				Date:        args[2],
				Issuer:      issuer,
				Batch:       batch,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("issuer", "", "Vaccine issuer")
	cmd.Flags().String("batch", "", "Vaccine batch number")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
