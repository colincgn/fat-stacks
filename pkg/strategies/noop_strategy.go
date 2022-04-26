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

	for msg := range no.datafeed.Listen() {
		log.Info().Msgf("%s", msg)
	}
	log.Debug().Msg("Stopping strategy, datafeed listening channel must have been closed.")
	return nil
}
