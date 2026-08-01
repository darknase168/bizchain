package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/keberangkatan/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "keberangkatan module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateKeberangkatan(),
		GetCmdAdvanceStage(),
		GetCmdUpdateDeparture(),
		GetCmdAddBaggage(),
		GetCmdUpdateBaggageStatus(),
	)

	return txCmd
}

// GetCmdCreateKeberangkatan starts tracking a journey
func GetCmdCreateKeberangkatan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-keberangkatan [jamaah] [paket-id] [pembayaran-id]",
		Short: "Start tracking a departure journey (stage 1: Daftar)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			paketID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			pembayaranID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			departure, _ := cmd.Flags().GetString("departure")
			returnDate, _ := cmd.Flags().GetString("return")
			manasik, _ := cmd.Flags().GetString("manasik")

			msg := &types.MsgCreateKeberangkatan{
				Creator:       clientCtx.GetFromAddress().String(),
				Jamaah:        args[0],
				PaketId:       paketID,
				PembayaranId:  pembayaranID,
				DepartureDate: departure,
				ReturnDate:    returnDate,
				ManasikDate:   manasik,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("departure", "", "Departure date")
	cmd.Flags().String("return", "", "Return date")
	cmd.Flags().String("manasik", "", "Manasik date")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAdvanceStage advances a journey to the next stage
func GetCmdAdvanceStage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "advance-stage [id]",
		Short: "Advance journey to the next stage (1-9)",
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

			msg := &types.MsgAdvanceStage{
				Creator: clientCtx.GetFromAddress().String(),
				Id:      id,
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

// GetCmdUpdateDeparture updates departure/return/manasik dates
func GetCmdUpdateDeparture() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-departure [id]",
		Short: "Update departure/return/manasik dates",
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

			departure, _ := cmd.Flags().GetString("departure")
			returnDate, _ := cmd.Flags().GetString("return")
			manasik, _ := cmd.Flags().GetString("manasik")
			hotelConfirm, _ := cmd.Flags().GetString("hotel-confirm")
			airlineConfirm, _ := cmd.Flags().GetString("airline-confirm")

			msg := &types.MsgUpdateDeparture{
				Creator:        clientCtx.GetFromAddress().String(),
				Id:             id,
				DepartureDate:  departure,
				ReturnDate:     returnDate,
				ManasikDate:    manasik,
				HotelConfirm:   hotelConfirm,
				AirlineConfirm: airlineConfirm,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("departure", "", "Departure date")
	cmd.Flags().String("return", "", "Return date")
	cmd.Flags().String("manasik", "", "Manasik date")
	cmd.Flags().String("hotel-confirm", "", "Hotel confirmation")
	cmd.Flags().String("airline-confirm", "", "Airline confirmation")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdAddBaggage adds a checked baggage item
func GetCmdAddBaggage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-baggage [id] [tag] [weight]",
		Short: "Add a checked baggage item (QR/NFC tracking)",
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

			msg := &types.MsgAddBaggage{
				Creator: clientCtx.GetFromAddress().String(),
				Id:      id,
				Tag:     args[1],
				Weight:  args[2],
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

// GetCmdUpdateBaggageStatus updates baggage tracking status
func GetCmdUpdateBaggageStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-baggage-status [id] [baggage-id] [status]",
		Short: "Update baggage status (checked_in, in_transit, arrived, delivered)",
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
			baggageID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateBaggageStatus{
				Creator:   clientCtx.GetFromAddress().String(),
				Id:        id,
				BaggageId: baggageID,
				Status:    args[2],
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
