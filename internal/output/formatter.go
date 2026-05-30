package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"

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

// WriteTable writes v as an aligned table to w.
func WriteTable(w io.Writer, v any) error {
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

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	for i, col := range cols {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, strings.ToUpper(col))
	}
	fmt.Fprintln(tw)

	for _, item := range list {
		for i, col := range cols {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			val := item[col]
			if val == nil {
				val = ""
			}
			s := fmt.Sprintf("%v", val)
			if len(s) > 80 {
				s = s[:77] + "..."
			}
			fmt.Fprint(tw, s)
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

// WriteCommandPayload writes the output in the specified format with unwrap and filtering.
func WriteCommandPayload(format Format, data any, fields string, jqExpr string) error {
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
		return WriteTable(os.Stdout, data)
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
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
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
		fmt.Fprintf(tw, "%s:\t%v\n", name, val)
	}
	tw.Flush()
}
