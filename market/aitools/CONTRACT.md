# Read-only market AI tool contract

## Supported operations

| Tool | Input | Output | Side effects |
|---|---|---|---|
| `market_status` | Empty object | `trade_count` | None |
| `recent_trades` | Optional `limit` | Public trade summaries | None |

## Input rules

- `market_status` accepts no non-empty argument object.
- `recent_trades` defaults to 10 records when the limit is missing or non-positive.
- The adapter caps the requested limit at 50.
- Unknown tool names and malformed JSON return errors.
- Unknown JSON fields are rejected by a schema-aware caller; callers should validate against `Specs()`.

## Output rules

- Trade timestamps are normalized to UTC.
- Prices and quantities are serialized as decimal strings to avoid floating-point precision loss.
- The response includes only the trade identifier, symbol, price, quantity, and timestamp.
- Nil trade entries are skipped.
- Internal order identifiers and account state must not be added to the projection without an explicit security review.

## Integration guidance

Translate `Specs()` into the schema format required by the host agent protocol, then route tool calls to `Call(reader, name, arguments)`. Keep the `Reader` implementation read-only. Do not expose order placement, cancellation, account changes, or other mutating methods through this adapter.

## Validation

From the repository root:

```bash
cd market
go test ./aitools
go test ./...
cd ..
python3 build.py
```

Do not claim checks passed unless they were actually run. The full repository workflow also requires the real diagnostic artifact described in the root contribution guide.
