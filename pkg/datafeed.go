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
	Run(ctx context.Context) error
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
