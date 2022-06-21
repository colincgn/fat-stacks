package pkg

import (
	"container/list"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type Side int
type GreenRed int

const (
	Buy = iota
	Sell
)

const (
	Green = iota
	Red
)

func (s Side) String() string {
	return [...]string{"Buy", "Sell"}[s]
}

func (g GreenRed) String() string {
	return [...]string{"Green", "Red"}[g]
}

type Engine struct {
}

type TickData struct {
	Side    Side
	Symbol  string
	Price   decimal.Decimal
	Volume  decimal.Decimal
	Time    int64
	IsBuyer bool
}

type DataFeed interface {
	Connect(ctx context.Context) error
	Subscribe(ctx context.Context, symbol string) error
	Run(ctx context.Context) error
	Listen() <-chan TickData
	Stop()
}

type aggTrade struct {
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Timestamp int64  `json:"T"`
	Volume    string `json:"q"`
	IsSeller  bool   `json:"m"`
}

type CandleData struct {
	Historical map[int64]list.List
}

type HLOC struct {
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	OpenTime  int64
	CloseTime int64
}

func (td TickData) String() string {
	var buyer GreenRed
	if td.IsBuyer {
		buyer = Green
	} else {
		buyer = Red
	}
	return fmt.Sprintf("%s - Price: %s, Volume: %s, Time: %s, Buyer: %s", td.Symbol, td.Price, td.Volume, time.UnixMilli(td.Time).Format("15:04:05.00000"), buyer)
}

type WsDataFeed struct {
	c                 *websocket.Conn
	subscriptionCount int
	tickDataFeed      chan TickData
	candlestickFeed   chan HLOC
}

func NewWsDataFeed() *WsDataFeed {
	return &WsDataFeed{
		subscriptionCount: 0,
		tickDataFeed:      make(chan TickData, 5),
	}
}

func (ws *WsDataFeed) Stop() {
	log.Debug().Msg("Stop called in DataFeed, closing channel")
	close(ws.tickDataFeed)
}

func (ws *WsDataFeed) Listen() <-chan TickData {
	return ws.tickDataFeed
}

func (ws *WsDataFeed) ListenHLOC() <-chan HLOC {
	return ws.candlestickFeed
}

func (ws *WsDataFeed) Run(ctx context.Context) error {
	var agg aggTrade
	for {
		err := wsjson.Read(ctx, ws.c, &agg)
		if err != nil {
			if ctx.Err() == context.Canceled {
				log.Debug().Msg("Websocket datafeed closed, shutting down reader.")
			} else {
				log.Error().Err(err).Msg("expected to be disconnected with a context cancel error but got")

			}
			break
		}

		tick := generateTickEvent(agg)
		select {
		case ws.tickDataFeed <- tick:
		default:
			log.Warn().Msg("Failed to write to subscribed feed, dropping message")
		}
		//select {
		//case ws.candlestickFeed <- generateHLOCEvent(tick):
		//default:
		//	log.Warn().Msg("Failed to write to candlestick feed, dropping message")
		//}
	}
	return nil
}

func generateTickEvent(agg aggTrade) TickData {
	price, _ := decimal.NewFromString(agg.Price)
	volume, _ := decimal.NewFromString(agg.Volume)
	return TickData{
		Symbol:  agg.Symbol,
		Price:   price,
		Volume:  volume,
		Time:    agg.Timestamp,
		IsBuyer: !agg.IsSeller,
	}
}

func generateHLOCEvent(tick TickData) HLOC {
	return HLOC{}
}

func (ws *WsDataFeed) Subscribe(ctx context.Context, symbol string) error {
	ws.subscriptionCount += 1
	subscription := fmt.Sprintf(`{"method": "SUBSCRIBE", "params":["%s@aggTrade"],"id":%d}`, symbol, ws.subscriptionCount)
	err := ws.c.Write(ctx, websocket.MessageText, []byte(subscription))
	if err != nil {
		log.Error().Err(err).Msgf("failed to subscribe to symbol %s", symbol)
		return err
	}
	_, msg, err := ws.c.Read(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read from exchange")
		return err
	}
	log.Info().Msgf("Successfully subscribed to channel: %s - %s\n", symbol, string(msg))
	return nil
}

func (ws *WsDataFeed) Connect(ctx context.Context) error {
	//Test System
	c, _, err := websocket.Dial(ctx, "wss://stream.binancefuture.com/ws/", nil)

	//Live System
	//c, _, err := websocket.Dial(ctx, "wss://fstream.binance.com/ws/", nil)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to host")
		return err
	}
	ws.c = c
	return nil
}
