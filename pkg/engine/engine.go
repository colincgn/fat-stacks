package engine

type Engine interface {
	Run()
}

type TradingEngine struct {
}

func (b TradingEngine) Run() {

}

func New() Engine {
	return TradingEngine{}
}
