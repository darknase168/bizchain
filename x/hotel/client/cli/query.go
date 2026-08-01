package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bizchain/blockchain/x/hotel/types"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the hotel module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdListHotel(),
		GetCmdShowHotel(),
	)

	return queryCmd
}

// GetCmdListHotel queries all hotels
func GetCmdListHotel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-hotel",
		Short: "List all hotels",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.HotelAll(cmd.Context(), &types.QueryAllHotelRequest{})
			if err != nil {
				return err
			}

			for _, h := range resp.Hotel {
				fmt.Printf("ID: %d | %s | %s | %s stars | %s/night | Rooms: %d | Status: %s\n",
					h.Id, h.Name, h.City, h.StarRating, h.PricePerNight, h.AvailableRooms, h.Status)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdShowHotel queries a specific hotel
func GetCmdShowHotel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-hotel [id]",
		Short: "Show a hotel by ID",
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
			resp, err := queryClient.Hotel(cmd.Context(), &types.QueryGetHotelRequest{Id: id})
			if err != nil {
				return err
			}

			h := resp.Hotel
			fmt.Printf("ID: %d\n", h.Id)
			fmt.Printf("Name: %s\n", h.Name)
			fmt.Printf("City: %s\n", h.City)
			fmt.Printf("Address: %s\n", h.Address)
			fmt.Printf("Star Rating: %s\n", h.StarRating)
			fmt.Printf("Price/Night: %s\n", h.PricePerNight)
			fmt.Printf("Room Type: %s\n", h.RoomType)
			fmt.Printf("Available Rooms: %d\n", h.AvailableRooms)
			fmt.Printf("Distance to Haram: %s m\n", h.DistanceHaram)
			fmt.Printf("Status: %s\n", h.Status)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
