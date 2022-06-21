package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
	"time"
)

func fakeNow() time.Time {
	return time.Date(2019, time.January, 1, 1, 23, 1, 0, time.UTC)
}

func TestGet(t *testing.T) {
	d, _ := time.ParseDuration("5m")
	expected := fakeNow().Truncate(d)
	if expected.Second() != 0 {
		t.Fatal("Expected seconds to be zero")
	}
}

type testAgg struct {
	A            int    `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeId int64  `json:"f"`
	LastTradeId  int64  `json:"l"`
	Timestamp    int64  `json:"T"`
	MarketMaker  bool   `json:"m"`
}

func (t testAgg) String() string {
	return fmt.Sprintf("Price : %v - Time: %v", t.Price, time.UnixMilli(t.Timestamp))
}
func TestReadAggTrades(t *testing.T) {
	file, _ := os.Open("./tick_data.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	to, err := decoder.Token()
	if err != nil {
		t.Fatal("Failed to read token")
	}
	log.Info().Msgf("%T: %v\n", to, to)
	var msg testAgg
	for decoder.More() {
		decoder.Decode(&msg)
		log.Info().Msgf("%v", msg)
	}
	to, err = decoder.Token()

	for i := 0; i < 100; i++ {

	}
}
