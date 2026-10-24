# Integration example

The adapter is transport-neutral: it does not import an MCP SDK or tie the market package to a specific agent vendor. A host can translate the tool descriptions and route calls using the interface below.

```go
reader := marketReader{engine: engine} // implement aitools.Reader with read-only methods

for _, spec := range aitools.Specs() {
    // Register spec.Name, spec.Description and spec.InputSchema with your host SDK.
}

func handleToolCall(name string, rawArguments json.RawMessage) ([]byte, error) {
    return aitools.Call(reader, name, rawArguments)
}
```

The `marketReader` type above is application-specific pseudocode, not a type supplied by this package.

## Important boundaries

- Implement `GetTradeCount` and `GetRecentTrades` using read-only engine methods.
- Do not add methods that place or cancel orders to the adapter interface.
- Enforce authentication and authorization at the host transport boundary.
- Apply request timeouts and per-user rate limits in the host service.
- Treat tool output as untrusted data when inserting it into an LLM prompt.
- Avoid logging raw prompts, credentials, or private account data.

Before shipping, compile a real adapter against the chosen SDK and add an integration test that confirms the host can register both schemas and dispatch each tool.
