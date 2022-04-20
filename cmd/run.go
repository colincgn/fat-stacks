package cmd

import (
	"context"
	"fat-stacks/pkg"
	"fmt"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the stack application",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithCancel(context.Background())
		var g run.Group
		datafeed := pkg.WsDataFeed{}
		{
			g.Add(func() error {
				return datafeed.Run(ctx)
			}, func(error) {
				cancel()
			})
		}
		{
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
			g.Add(func() error {
				<-signalCh
				fmt.Println("Shutting Down")
				return nil
			}, func(err error) {
				close(signalCh)
			})
		}
		g.Run()
		fmt.Println("Shut down complete.  Exiting...")
		return nil
	},
}
