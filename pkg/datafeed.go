package pkg

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type Side int

const (
	Buy = iota
	Sell
)

func (s Side) String() string {
	return [...]string{"Buy", "Sell"}[s]
}

type TickData struct {
	Side   Side
	Symbol string
	Price  decimal.Decimal
	Volume decimal.Decimal
	Time   int64
}

type DataFeed interface {
	Run(ctx context.Context) error
}

type aggTrade struct {
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Timestamp int64  `json:"T"`
	Volume    string `json:"q"`
}

func (a aggTrade) String() string {
	return fmt.Sprintf("%s - Price: %s, Volume: %s, Time: %s", a.Symbol, a.Price, a.Volume, time.UnixMilli(a.Timestamp).Format("15:04:05.00000"))
}

type WsDataFeed struct{}

func (ws WsDataFeed) Run(ctx context.Context) error {
	c, _, err := websocket.Dial(ctx, "wss://stream.binancefuture.com/ws/btcusdt@aggTrade", nil)
	if err != nil {
		fmt.Println("Error connecting to host", err)
		return err
	}
	defer c.Close(websocket.StatusInternalError, "Unable to process message")

	var v aggTrade
	for {
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			fmt.Println("Error reading json:", err)
			break
		}
		fmt.Printf("%s\n", v)
	}
	return nil
}
