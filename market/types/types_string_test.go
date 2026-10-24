package types

import "testing"

func TestOrderSideString(t *testing.T) {
	tests := []struct {
		value OrderSide
		want  string
	}{
		{Buy, "buy"},
		{Sell, "sell"},
		{OrderSide(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.value.String(); got != tt.want {
			t.Errorf("OrderSide(%d).String() = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestOrderTypeString(t *testing.T) {
	tests := []struct {
		value OrderType
		want  string
	}{
		{Limit, "limit"},
		{Market, "market"},
		{StopLimit, "stop_limit"},
		{StopMarket, "stop_market"},
		{TrailingStop, "trailing_stop"},
		{Iceberg, "iceberg"},
		{OrderType(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.value.String(); got != tt.want {
			t.Errorf("OrderType(%d).String() = %q, want %q", tt.value, got, tt.want)
		}
	}
}
