package client

import (
	"encoding/json"
	"fmt"
)

// SafeCode handles both string and int JSON code values.
type SafeCode string

func (c *SafeCode) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*c = SafeCode(s)
		return nil
	}
	// Try number
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*c = SafeCode(fmt.Sprintf("%d", n))
		return nil
	}
	return fmt.Errorf("code must be string or int, got %s", string(data))
}

func (c SafeCode) IsZero() bool {
	return string(c) == "0"
}

func (c SafeCode) String() string {
	return string(c)
}
