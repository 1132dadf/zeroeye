package aitools

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/tent-of-trials/market/types"
)

type fakeReader struct {
	count int64
	trades []*types.Trade
	lastLimit int
}

func (f *fakeReader) GetTradeCount() int64 { return f.count }
func (f *fakeReader) GetRecentTrades(limit int) []*types.Trade {
	f.lastLimit = limit
	if limit > len(f.trades) {
		limit = len(f.trades)
	}
	return f.trades[len(f.trades)-limit:]
}

func TestMarketStatusReadOnlyResult(t *testing.T) {
	reader := &fakeReader{count: 42}
	raw, err := Call(reader, ToolMarketStatus, json.RawMessage(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"trade_count":42`) {
		t.Fatalf("unexpected result: %s", raw)
	}
}

func TestRecentTradesClampsLimitAndRedactsOrderIDs(t *testing.T) {
	reader := &fakeReader{trades: []*types.Trade{{
		ID: "trade-1", Symbol: types.Symbol("BTC-USD"),
		Price: decimal.NewFromInt(10), Quantity: decimal.NewFromInt(2),
		Timestamp: time.Date(2026, 1, 2, 3, 4, 5, 0, time.FixedZone("offset", 3600)),
	}}}
	raw, err := Call(reader, ToolRecentTrades, json.RawMessage(`{"limit":999}`))
	if err != nil {
		t.Fatal(err)
	}
	if reader.lastLimit != maxTradeLimit {
		t.Fatalf("limit not clamped: %d", reader.lastLimit)
	}
	if !strings.Contains(string(raw), `"symbol":"BTC-USD"`) || strings.Contains(string(raw), "buy_order_id") {
		t.Fatalf("unexpected trade projection: %s", raw)
	}
	if !strings.Contains(string(raw), "2026-01-02T02:04:05Z") {
		t.Fatalf("timestamp was not normalized to UTC: %s", raw)
	}
}

func TestUnknownToolIsRejected(t *testing.T) {
	_, err := Call(&fakeReader{}, "place_order", json.RawMessage(`{}`))
	if err == nil {
		t.Fatal("expected unknown tool error")
	}
}
