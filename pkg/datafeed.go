package pkg

import (
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

type TickData struct {
	Side   Side
	Symbol string
	Price  decimal.Decimal
	Volume decimal.Decimal
	Time   int64
}

type DataFeed interface {
	Connect(ctx context.Context) error
	Subscribe(ctx context.Context, symbol string) error
	Run(ctx context.Context) error
	Listen() <-chan aggTrade
	Stop()
}

type aggTrade struct {
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Timestamp int64  `json:"T"`
	Volume    string `json:"q"`
	IsSeller  bool   `json:"m"`
}

func (a aggTrade) String() string {
	var buyer GreenRed
	if a.IsSeller {
		buyer = Red
	} else {
		buyer = Green
	}
	return fmt.Sprintf("%s - Price: %s, Volume: %s, Time: %s, Buyer: %s", a.Symbol, a.Price, a.Volume, time.UnixMilli(a.Timestamp).Format("15:04:05.00000"), buyer)
}

type WsDataFeed struct {
	c                 *websocket.Conn
	subscriptionCount int
	subscribedFeed    chan aggTrade
}

func NewWsDataFeed() *WsDataFeed {
	return &WsDataFeed{
		subscriptionCount: 0,
		subscribedFeed:    make(chan aggTrade, 100),
	}
}

func (ws *WsDataFeed) Stop() {
	log.Debug().Msg("Stop called in DataFeed, closing channel")
	close(ws.subscribedFeed)
}

func (ws *WsDataFeed) Listen() <-chan aggTrade {
	return ws.subscribedFeed
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
		ws.subscribedFeed <- agg
	}
	return nil
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
	c, _, err := websocket.Dial(ctx, "wss://stream.binancefuture.com/ws/", nil)
	//c, _, err := websocket.Dial(ctx, "wss://fstream.binance.com/ws/", nil)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to host")
		return err
	}
	ws.c = c
	return nil
}
