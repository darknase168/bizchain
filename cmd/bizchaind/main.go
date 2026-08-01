package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/bizchain/blockchain/app"
	"github.com/bizchain/blockchain/cmd/bizchaind/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "BIZCHAIN", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
