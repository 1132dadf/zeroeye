# Market AI read-tool adapter

Package `aitools` provides a transport-neutral read-only adapter for agent integrations. It exposes two tool specifications:

- `market_status`: total processed trade count.
- `recent_trades`: a bounded projection of recent public trades.

An MCP or agent-SDK integration can translate `Specs()` into its native tool-registration format and route tool calls to `Call(reader, name, arguments)`. The adapter depends on a narrow `Reader` interface, so it can be tested without starting the WebSocket server.

Security boundaries:
- Only read operations are registered; order placement and cancellation are not exposed.
- Recent trade requests default to 10 and are capped at 50.
- Responses intentionally omit buy/sell order IDs and internal order state.
- Unknown tools and invalid JSON return errors.

Validation from the repository root:

    cd market
    go test ./aitools
    cd ..
    python3 build.py

The full build and required encrypted diagnostic artifact still need to be generated in the target environment before final submission.
