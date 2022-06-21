package cmd

import (
	"context"
	"fat-stacks/pkg"
	tradingEngine "fat-stacks/pkg/engine"
	"fat-stacks/pkg/strategies"
	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
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
		defer cancel()
		var g run.Group

		datafeed := pkg.NewWsDataFeed()
		datafeed.Connect(ctx)
		datafeed.Subscribe(ctx, "btcusdt")
		engine := tradingEngine.New()
		noop := strategies.New(datafeed)
		engine.Run()
		{
			g.Add(func() error {
				return datafeed.Run(ctx)
			}, func(error) {
				datafeed.Stop()
			})
		}
		{
			g.Add(func() error {
				return noop.Run()
			}, func(err error) {
				log.Debug().Msg("Stopping strategy in run group, nothing to do here since it's stopped by a closing channel")
			})
		}
		{
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
			g.Add(func() error {
				<-signalCh
				log.Debug().Msg("Shutting Down Signal received.  Preparing to exit")
				cancel()
				return nil
			}, func(err error) {
				close(signalCh)
			})
		}
		g.Run()
		log.Info().Msg("Shut down complete.")
		return nil
	},
}
