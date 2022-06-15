package strategies

import (
	"fat-stacks/pkg"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type HLOC struct {
	High  decimal.Decimal
	Low   decimal.Decimal
	Open  decimal.Decimal
	Close decimal.Decimal
}

type NoOpStrategy struct {
	datafeed pkg.DataFeed
	hloc     *HLOC
}

func New(datafeed pkg.DataFeed) *NoOpStrategy {
	return &NoOpStrategy{
		hloc: &HLOC{
			High:  decimal.Zero,
			Low:   decimal.Zero,
			Open:  decimal.Zero,
			Close: decimal.Zero,
		},
		datafeed: datafeed,
	}
}

func (no *NoOpStrategy) Run() error {
	for {
		select {
		case tick, ok := <-no.datafeed.Listen():
			if !ok {
				return nil
			}
			log.Info().Msgf("%s", tick)
		case hloc, ok := <-no.datafeed.ListenHLOC():
			if !ok {
				return nil
			}
			log.Info().Msgf("%s", hloc)
		}
	}
}
