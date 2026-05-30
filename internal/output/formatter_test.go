package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]any{"status": "ok", "count": 42}
	err := WriteJSON(&buf, data)
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", parsed["status"])
	}
}

func TestResolveFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"table", FormatTable},
		{"raw", FormatRaw},
		{"unknown", FormatJSON},
		{"", FormatJSON},
	}

	for _, tt := range tests {
		result := ResolveFormat(tt.input, FormatJSON)
		if result != tt.expected {
			t.Errorf("ResolveFormat(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractList(t *testing.T) {
	slice := []map[string]any{{"a": "1"}, {"a": "2"}}
	result := extractList(slice)
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}

	wrapped := map[string]any{"data": slice}
	result = extractList(wrapped)
	if len(result) != 2 {
		t.Errorf("expected 2 items from wrapped, got %d", len(result))
	}

	result = extractList(nil)
	if result != nil {
		t.Error("expected nil for nil input")
	}
}

func TestWriteTable(t *testing.T) {
	var buf bytes.Buffer
	data := []map[string]any{
		{"name": "张三", "status": "通过"},
		{"name": "李四", "status": "待审批"},
	}
	err := WriteTable(&buf, data)
	if err != nil {
		t.Fatal(err)
	}
	output := buf.String()
	if output == "" {
		t.Fatal("expected non-empty table output")
	}
	if !contains(output, "NAME") || !contains(output, "STATUS") {
		t.Errorf("expected NAME and STATUS in output, got: %s", output)
	}
}

func TestUnwrapResponse(t *testing.T) {
	// data → result → data → records
	nested := map[string]any{
		"data": map[string]any{
			"result": map[string]any{
				"data": []any{
					map[string]any{"name": "张三"},
					map[string]any{"name": "李四"},
				},
			},
		},
	}
	result := unwrapResponse(nested)
	list, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 records, got %d", len(list))
	}
}

func TestSelectFields(t *testing.T) {
	data := []map[string]any{
		{"name": "张三", "status": "通过", "id": "1"},
		{"name": "李四", "status": "待审批", "id": "2"},
	}
	result := selectFields(data, "name,status")
	list, ok := result.([]map[string]any)
	if !ok {
		t.Fatal("expected []map[string]any")
	}
	if len(list[0]) != 2 {
		t.Errorf("expected 2 fields, got %d", len(list[0]))
	}
	if _, exists := list[0]["id"]; exists {
		t.Error("id should be filtered out")
	}
}

func TestApplyJQ(t *testing.T) {
	data := []any{
		map[string]any{"name": "张三", "status": "通过"},
		map[string]any{"name": "李四", "status": "待审批"},
	}
	result, err := applyJQ(data, ".[] | {n: .name}")
	if err != nil {
		t.Fatal(err)
	}
	list, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 results, got %d", len(list))
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
