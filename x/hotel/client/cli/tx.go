package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/hotel/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "hotel module transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateHotel(),
		GetCmdUpdateHotel(),
	)

	return txCmd
}

// GetCmdCreateHotel registers a new hotel
func GetCmdCreateHotel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-hotel [name] [city] [price-per-night]",
		Short: "Register a new hotel",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			address, _ := cmd.Flags().GetString("address")
			starRating, _ := cmd.Flags().GetString("stars")
			roomType, _ := cmd.Flags().GetString("room-type")
			availableRooms, _ := cmd.Flags().GetUint64("rooms")
			distanceHaram, _ := cmd.Flags().GetString("distance-haram")

			msg := &types.MsgCreateHotel{
				Creator:        clientCtx.GetFromAddress().String(),
				Name:           args[0],
				City:           args[1],
				Address:        address,
				StarRating:     starRating,
				PricePerNight:  args[2],
				RoomType:       roomType,
				AvailableRooms: availableRooms,
				DistanceHaram:  distanceHaram,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("address", "", "Hotel address")
	cmd.Flags().String("stars", "", "Star rating (e.g. 5)")
	cmd.Flags().String("room-type", "", "Room type")
	cmd.Flags().Uint64("rooms", 0, "Available rooms")
	cmd.Flags().String("distance-haram", "", "Distance to Masjidil Haram (meters)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateHotel updates hotel data
func GetCmdUpdateHotel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-hotel [id] [name] [city] [price-per-night]",
		Short: "Update hotel data",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			address, _ := cmd.Flags().GetString("address")
			starRating, _ := cmd.Flags().GetString("stars")
			roomType, _ := cmd.Flags().GetString("room-type")
			availableRooms, _ := cmd.Flags().GetUint64("rooms")
			distanceHaram, _ := cmd.Flags().GetString("distance-haram")
			status, _ := cmd.Flags().GetString("status")

			msg := &types.MsgUpdateHotel{
				Creator:        clientCtx.GetFromAddress().String(),
				Id:             id,
				Name:           args[1],
				City:           args[2],
				Address:        address,
				StarRating:     starRating,
				PricePerNight:  args[3],
				RoomType:       roomType,
				AvailableRooms: availableRooms,
				DistanceHaram:  distanceHaram,
				Status:         status,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("address", "", "Hotel address")
	cmd.Flags().String("stars", "", "Star rating")
	cmd.Flags().String("room-type", "", "Room type")
	cmd.Flags().Uint64("rooms", 0, "Available rooms")
	cmd.Flags().String("distance-haram", "", "Distance to Haram (meters)")
	cmd.Flags().String("status", "", "Status (active, inactive)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
