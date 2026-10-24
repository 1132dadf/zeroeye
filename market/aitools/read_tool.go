// Package aitools exposes a small, read-only tool adapter for AI integrations.
// It deliberately depends on an interface rather than a concrete engine so it
// can be registered with MCP, an agent SDK, or a local tool router.
package aitools

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/tent-of-trials/market/types"
)

const (
	ToolMarketStatus = "market_status"
	ToolRecentTrades = "recent_trades"
	maxTradeLimit    = 50
)

type Reader interface {
	GetTradeCount() int64
	GetRecentTrades(limit int) []*types.Trade
}

// ToolSpec is a transport-neutral tool description. An MCP or other AI adapter
// can translate this into the schema format expected by its SDK.
type ToolSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

type readRequest struct {
	Limit int `json:"limit"`
}

type tradeSummary struct {
	ID        string `json:"id"`
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Timestamp string `json:"timestamp"`
}

type result struct {
	Tool        string         `json:"tool"`
	TradeCount  *int64         `json:"trade_count,omitempty"`
	Returned    *int           `json:"returned,omitempty"`
	Trades      []tradeSummary `json:"trades,omitempty"`
}

// Specs returns only read operations. No order placement, cancellation, or
// account mutation tool is exposed through this adapter.
func Specs() []ToolSpec {
	return []ToolSpec{
		{
			Name: ToolMarketStatus,
			Description: "Read the market engine's total processed trade count. Does not place or modify orders.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{},
				"additionalProperties": false,
			},
		},
		{
			Name: ToolRecentTrades,
			Description: "Read recent public trades. The limit is clamped to 1..50. Does not place or modify orders.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit": map[string]any{
						"type": "integer",
						"description": "Maximum number of recent trades to return (1-50).",
						"minimum": 1,
						"maximum": maxTradeLimit,
					},
				},
				"additionalProperties": false,
			},
		},
	}
}

// Call dispatches a tool by name and returns JSON suitable for an AI tool result.
// Unknown tools and malformed arguments return errors instead of silently
// falling through to a mutating operation.
func Call(reader Reader, name string, arguments json.RawMessage) (json.RawMessage, error) {
	if reader == nil {
		return nil, errors.New("market read tool: reader is nil")
	}

	switch name {
	case ToolMarketStatus:
		if len(arguments) > 0 && strings.TrimSpace(string(arguments)) != "null" && strings.TrimSpace(string(arguments)) != "{}" {
			var extra map[string]json.RawMessage
			if err := json.Unmarshal(arguments, &extra); err != nil {
				return nil, fmt.Errorf("market_status arguments: %w", err)
			}
			if len(extra) != 0 {
				return nil, errors.New("market_status accepts no arguments")
			}
		}
		count := reader.GetTradeCount()
		return json.Marshal(result{Tool: name, TradeCount: &count})

	case ToolRecentTrades:
		var request readRequest
		if len(arguments) > 0 && strings.TrimSpace(string(arguments)) != "null" {
			if err := json.Unmarshal(arguments, &request); err != nil {
				return nil, fmt.Errorf("recent_trades arguments: %w", err)
			}
		}
		if request.Limit <= 0 {
			request.Limit = 10
		}
		if request.Limit > maxTradeLimit {
			request.Limit = maxTradeLimit
		}
		trades := reader.GetRecentTrades(request.Limit)
		summaries := make([]tradeSummary, 0, len(trades))
		for _, trade := range trades {
			if trade == nil {
				continue
			}
			summaries = append(summaries, tradeSummary{
				ID: trade.ID,
				Symbol: string(trade.Symbol),
				Price: trade.Price.String(),
				Quantity: trade.Quantity.String(),
				Timestamp: trade.Timestamp.UTC().Format("2006-01-02T15:04:05.999999999Z"),
			})
		}
		return json.Marshal(result{Tool: name, Returned: intPointer(len(summaries)), Trades: summaries})

	default:
		return nil, fmt.Errorf("market read tool: unknown tool %q", name)
	}
}

func intPointer(value int) *int { return &value }
