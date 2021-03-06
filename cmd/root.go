package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "stack",
	Short: "stack is a fat stack creator",
	Long:  `stack will do it's best to help you create useful trading strategies and make those strategies testable'`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Err(err)
		os.Exit(1)
	}
}
