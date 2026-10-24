# Order book invariants

This note records intended invariants for contributors changing `market/orderbook`. Confirm each invariant against implementation and tests before relying on it in production.

## Lifecycle and ownership

- A closed order book should reject new orders.
- The book owns its internal level/order maps; callers should not mutate internal slices or maps after passing objects into the API.
- Operations that read or mutate shared state must follow the book's locking discipline.
- Order IDs should be unique within the book's active-order set.

## Price-time behavior

- Bid levels should be ordered from highest price to lowest price.
- Ask levels should be ordered from lowest price to highest price.
- At the same price, matching should respect the project's documented sequence/time priority.
- Decimal arithmetic should be used for prices and quantities; avoid converting values to binary floating point for comparisons.

## Matching and cancellation

- A fill must not exceed either the resting order's remaining quantity or the incoming order's remaining quantity.
- Filled orders should not remain available for another fill.
- Cancellation of an unknown order should return a documented error or no-op consistently; do not silently corrupt level totals.
- A partial fill should update remaining quantity and order status together.

## Suggested regression cases

1. Empty book accepts a valid limit order.
2. Crossing buy/sell orders match at the expected price.
3. Partial fills preserve the remaining quantity.
4. Multiple orders at one price preserve priority.
5. Cancelling one order leaves other orders at the same price intact.
6. Closed book rejects a new order.
7. Concurrent reads do not race with updates.

Run the focused package tests and the full repository diagnostic workflow before submitting changes upstream.
