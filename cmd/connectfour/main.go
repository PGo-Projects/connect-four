package main

import (
	"github.com/PGo-Projects/connect-four/internal/game"
	"github.com/spf13/cobra"
)

func main() {
	var cmdPlay = &cobra.Command{
		Use:   "play",
		Short: "Start a connect four game",
		Long:  "Play a connect four game against the computer or another player",
		Args:  cobra.ExactArgs(0),
		Run:   startGame,
	}

	var rootCmd = &cobra.Command{Use: "connectfour"}
	rootCmd.AddCommand(cmdPlay)
	rootCmd.Execute()
}

func startGame(cmd *cobra.Command, args []string) {
	connectfour := game.New()
	connectfour.Start()
}
