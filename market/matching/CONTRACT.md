# Matching engine safety expectations

The matching engine is part of a trading/risk platform. Changes should preserve deterministic matching behavior and must not weaken validation to make tests pass.

## Input and configuration

- Validate nil orders before dereferencing them.
- Validate symbol, side, order type, price, and quantity before forwarding an order to a book.
- Reject non-positive quantities and invalid limit prices.
- Treat fee rates and order limits as explicit configuration; document defaults and invalid values.
- Keep errors actionable but do not include secrets or private account data.

## Matching behavior

- A trade quantity must never exceed either order's remaining quantity.
- Decimal values should remain decimal through comparison and accounting.
- Update trade history and counters consistently with successful fills.
- Protect shared trade history and order-book maps with the established locking strategy.
- Ensure cancellation and timeout paths cannot resurrect completed orders.

## Suggested tests

- Nil order returns a safe error rather than panicking.
- Unknown symbol is rejected.
- Non-positive quantity is rejected.
- A crossing order produces a trade with bounded quantity.
- A non-crossing limit order remains open.
- Concurrent reads of trade count and recent trades are race-free.
- Fee configuration edge cases are handled consistently.

This file is a contributor checklist, not evidence that each item is already implemented. Validate with focused Go tests, race tests where practical, and the repository's full build/diagnostic process.
