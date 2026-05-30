package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/itchyny/gojq"
)

// Format represents an output format.
type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatRaw   Format = "raw"
)

// ResolveFormat determines the output format from the --format flag value.
func ResolveFormat(formatFlag string, fallback Format) Format {
	switch strings.ToLower(formatFlag) {
	case "json":
		return FormatJSON
	case "table":
		return FormatTable
	case "raw":
		return FormatRaw
	default:
		return fallback
	}
}

// WriteJSON writes v as pretty-printed JSON to w.
func WriteJSON(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(v)
}

// WriteRaw writes v as compact JSON to w.
func WriteRaw(w io.Writer, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// WriteTable writes v as an aligned table to w, using labels for column headers.
func WriteTable(w io.Writer, v any, labels map[string]string) error {
	list := extractList(v)
	if len(list) == 0 {
		writeKeyValue(w, v)
		return nil
	}

	colSet := make(map[string]struct{})
	for _, item := range list {
		for k := range item {
			colSet[k] = struct{}{}
		}
	}

	cols := make([]string, 0, len(colSet))
	for k := range colSet {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	// Build header and all rows as strings
	rows := make([][]string, 0, len(list)+1)
	header := make([]string, len(cols))
	for i, col := range cols {
		if label, ok := labels[col]; ok && label != "" {
			header[i] = label
		} else {
			header[i] = col
		}
	}
	rows = append(rows, header)

	for _, item := range list {
		row := make([]string, len(cols))
		for i, col := range cols {
			val := item[col]
			if val == nil {
				val = ""
			}
			s := fmt.Sprintf("%v", val)
			row[i] = s
		}
		rows = append(rows, row)
	}

	// Calculate display width for each column
	colWidths := make([]int, len(cols))
	for _, row := range rows {
		for i, cell := range row {
			w := displayWidth(cell)
			if w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	// Write with proper padding
	sep := "  "
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(w, sep)
			}
			fmt.Fprint(w, cell)
			pad := colWidths[i] - displayWidth(cell)
			if i < len(row)-1 {
				fmt.Fprint(w, strings.Repeat(" ", pad))
			}
		}
		fmt.Fprintln(w)
	}

	return nil
}

// displayWidth returns the visual width: ASCII=1, CJK=2.
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if isWide(r) {
			w += 2
		} else {
			w++
		}
	}
	if w == 0 && s != "" {
		// Fallback for non-UTF8
		return utf8.RuneCountInString(s)
	}
	return w
}

func isWide(r rune) bool {
	if r < 128 {
		return false
	}
	// CJK Unified Ideographs, CJK Compatibility, Hiragana, Katakana, Fullwidth
	return (r >= 0x1100 && r <= 0x115F) || // Hangul
		(r >= 0x2E80 && r <= 0xA4CF) || // CJK Radicals ~ Yi
		(r >= 0xAC00 && r <= 0xD7A3) || // Hangul Syllables
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility
		(r >= 0xFE10 && r <= 0xFE19) || // Vertical forms
		(r >= 0xFE30 && r <= 0xFE6F) || // CJK Compatibility Forms
		(r >= 0xFF00 && r <= 0xFF60) || // Fullwidth Forms
		(r >= 0xFFE0 && r <= 0xFFE6) || // Fullwidth Signs
		(r >= 0x20000 && r <= 0x2FFFD) || // CJK Extension B+
		(r >= 0x30000 && r <= 0x3FFFD) // CJK Extension G+
}

// WriteCommandPayload writes the output in the specified format with unwrap and filtering.
func WriteCommandPayload(format Format, data any, fields string, jqExpr string, labels map[string]string) error {
	// Step 1: Unwrap nested response
	data = unwrapResponse(data)

	// Step 2: Apply --fields filter
	if fields != "" {
		data = selectFields(data, fields)
	}

	// Step 3: Apply --jq filter
	if jqExpr != "" {
		var err error
		data, err = applyJQ(data, jqExpr)
		if err != nil {
			return fmt.Errorf("jq 表达式错误: %w", err)
		}
	}

	// Step 4: Format output
	switch format {
	case FormatJSON:
		return WriteJSON(os.Stdout, data)
	case FormatTable:
		return WriteTable(os.Stdout, data, labels)
	case FormatRaw:
		return WriteRaw(os.Stdout, data)
	default:
		return WriteJSON(os.Stdout, data)
	}
}

// unwrapResponse extracts the actual data from nested API response wrappers.
// Handles patterns like:
//   {"data": {"result": {"data": [...]}}}  → [...]
//   {"data": {"data": [...]}}              → [...]
//   {"data": [...]}                        → [...]
//   {"result": {...}}                      → {...}
func unwrapResponse(v any) any {
	m, ok := v.(map[string]any)
	if !ok {
		return v
	}

	// If there's a "data" key, recurse into it
	if dataVal, exists := m["data"]; exists {
		return unwrapResponse(dataVal)
	}

	// If there's a "result" key, use it
	if resultVal, exists := m["result"]; exists {
		switch rv := resultVal.(type) {
		case map[string]any:
			// result might contain {data: [...], recordsTotal: N}
			if innerData, ok := rv["data"]; ok {
				return innerData
			}
			return rv
		default:
			return rv
		}
	}

	return v
}

// selectFields filters a list of maps to only include specified fields.
func selectFields(v any, fields string) any {
	list := extractList(v)
	if len(list) == 0 {
		return v
	}
	fieldSet := make(map[string]bool)
	for _, f := range strings.Split(fields, ",") {
		fieldSet[strings.TrimSpace(f)] = true
	}
	var result []map[string]any
	for _, item := range list {
		filtered := make(map[string]any)
		for k, val := range item {
			if fieldSet[k] {
				filtered[k] = val
			}
		}
		result = append(result, filtered)
	}
	return result
}

// applyJQ applies a jq expression to the data.
func applyJQ(v any, expr string) (any, error) {
	query, err := gojq.Parse(expr)
	if err != nil {
		return nil, err
	}
	iter := query.Run(v)
	var results []any
	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := val.(error); ok {
			return nil, err
		}
		results = append(results, val)
	}
	if len(results) == 0 {
		return nil, nil
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

// extractList tries to find a list of maps in a nested structure.
func extractList(v any) []map[string]any {
	if v == nil {
		return nil
	}
	// Direct slice types
	switch typed := v.(type) {
	case []map[string]any:
		return typed
	case []any:
		var result []map[string]any
		for _, item := range typed {
			if im, ok := item.(map[string]any); ok {
				result = append(result, im)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	// Check nested wrappers
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	listKeys := []string{"data", "items", "list", "records", "result"}
	for _, key := range listKeys {
		if val, exists := m[key]; exists {
			switch typed := val.(type) {
			case []map[string]any:
				return typed
			case []any:
				var result []map[string]any
				for _, item := range typed {
					if im, ok := item.(map[string]any); ok {
						result = append(result, im)
					}
				}
				return result
			}
		}
	}
	return nil
}

func writeKeyValue(w io.Writer, v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		fmt.Fprintf(w, "%v\n", v)
		return
	}
	rt := rv.Type()
	maxKeyWidth := 0
	type kv struct {
		key string
		val string
	}
	var pairs []kv
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		val := rv.Field(i).Interface()
		name := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				name = parts[0]
			}
		}
		valStr := fmt.Sprintf("%v", val)
		pairs = append(pairs, kv{key: name, val: valStr})
		w := displayWidth(name)
		if w > maxKeyWidth {
			maxKeyWidth = w
		}
	}
	for _, p := range pairs {
		pad := maxKeyWidth - displayWidth(p.key)
		fmt.Fprintf(w, "%s:%s  %s\n", p.key, strings.Repeat(" ", pad), p.val)
	}
}
