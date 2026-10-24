package aitools

import "testing"

func TestSpecsExposeOnlyReadOperations(t *testing.T) {
	specs := Specs()
	if len(specs) != 2 {
		t.Fatalf("Specs() returned %d tools, want 2", len(specs))
	}
	seen := map[string]bool{}
	for _, spec := range specs {
		if spec.Name == "" || spec.Description == "" {
			t.Fatalf("tool has empty name or description: %#v", spec)
		}
		if seen[spec.Name] {
			t.Fatalf("duplicate tool name %q", spec.Name)
		}
		seen[spec.Name] = true
		if spec.InputSchema["type"] != "object" {
			t.Errorf("%s schema type = %v, want object", spec.Name, spec.InputSchema["type"])
		}
		if spec.InputSchema["additionalProperties"] != false {
			t.Errorf("%s should reject unknown properties", spec.Name)
		}
	}
	if !seen[ToolMarketStatus] || !seen[ToolRecentTrades] {
		t.Fatalf("missing expected read tools: %v", seen)
	}
	if seen["place_order"] || seen["cancel_order"] {
		t.Fatal("mutating trading tools must not be exposed by the read-only adapter")
	}
}

func TestRecentTradesSchemaBoundsLimit(t *testing.T) {
	var recent *ToolSpec
	for _, spec := range Specs() {
		if spec.Name == ToolRecentTrades {
			recent = &spec
			break
		}
	}
	if recent == nil {
		t.Fatal("recent_trades spec missing")
	}
	props, ok := recent.InputSchema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("properties has unexpected type %T", recent.InputSchema["properties"])
	}
	limit, ok := props["limit"].(map[string]any)
	if !ok {
		t.Fatalf("limit schema has unexpected type %T", props["limit"])
	}
	if limit["minimum"] != 1 || limit["maximum"] != maxTradeLimit {
		t.Errorf("limit bounds = (%v, %v), want (1, %d)", limit["minimum"], limit["maximum"], maxTradeLimit)
	}
}
