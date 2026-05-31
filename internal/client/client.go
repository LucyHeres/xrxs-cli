package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/LucyHeres/xrxs-cli/internal/auth"
	"github.com/LucyHeres/xrxs-cli/pkg/config"
)

// Response is the standard API response format.
type Response struct {
	Code    SafeCode        `json:"code"`
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Client wraps an HTTP client with session cookies and CSRF token.
type Client struct {
	BaseURL    string
	Session    *auth.Session
	HTTPClient *http.Client
	Verbose    bool
	DryRun     bool
}

// NewClient creates an API client from a session.
func NewClient(session *auth.Session, verbose, dryRun bool) *Client {
	jar := session.CookieJar()
	return &Client{
		BaseURL:    session.BaseURL,
		Session:    session,
		HTTPClient: &http.Client{Jar: jar, Timeout: config.HTTPTimeout},
		Verbose:    verbose,
		DryRun:     dryRun,
	}
}

// Post sends a POST request with form-encoded parameters.
func (c *Client) Post(ctx context.Context, path string, params url.Values) (*Response, error) {
	return c.doFormRequest(ctx, "POST", path, params)
}

// PostJSON sends a POST request with a JSON body.
func (c *Client) PostJSON(ctx context.Context, path string, body any) (*Response, error) {
	return c.doJSONRequest(ctx, "POST", path, body)
}

// PostFormJSON sends a POST with form-urlencoded data, encoding nested maps/structs
// like qs.stringify with {arrayFormat: 'brackets'}.
func (c *Client) PostFormJSON(ctx context.Context, path string, body map[string]any) (*Response, error) {
	params := flattenParams(body, "")
	bodyStr := params.Encode()
	c.logf("=> Body: %s", bodyStr)
	return c.doRequest(ctx, "POST", path, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(bodyStr))
}

// Get sends a GET request with query parameters.
func (c *Client) Get(ctx context.Context, path string, params url.Values) (*Response, error) {
	return c.doFormRequest(ctx, "GET", path, params)
}

func (c *Client) doFormRequest(ctx context.Context, method, path string, params url.Values) (*Response, error) {
	var body io.Reader
	var bodyStr string
	fullPath := path
	if method == "GET" && params != nil {
		fullPath = path + "?" + encodeParams(params)
	} else if params != nil {
		bodyStr = encodeParams(params)
		body = strings.NewReader(bodyStr)
	}
	c.logf("=> Body: %s", bodyStr)
	return c.doRequest(ctx, method, fullPath, "application/x-www-form-urlencoded; charset=UTF-8", body)
}

func (c *Client) doJSONRequest(ctx context.Context, method, path string, body any) (*Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}
	c.logf("=> Body: %s", string(data))
	return c.doRequest(ctx, method, path, "application/json; charset=UTF-8", strings.NewReader(string(data)))
}

func (c *Client) doRequest(ctx context.Context, method, path, contentType string, body io.Reader) (*Response, error) {
	fullURL := strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")

	if c.DryRun {
		c.logf("[DRY RUN] %s %s", method, fullURL)
		return &Response{Code: "0", Status: true, Message: "dry run"}, nil
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	if c.Session.CSRFToken != "" {
		req.Header.Set("X-CSRF-TOKEN", c.Session.CSRFToken)
	}
	req.Header.Set("XRXS-Language", "zh")

	c.logf("=> %s %s", method, fullURL)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, config.MaxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	c.logf("<= %d (%d bytes)", resp.StatusCode, len(bodyBytes))
	c.logf("<= Body: %s", string(bodyBytes[:min(len(bodyBytes), 500)]))

	var apiResp Response
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w (body: %s)", err, string(bodyBytes[:min(len(bodyBytes), 200)]))
	}

	if !apiResp.Code.IsZero() || !apiResp.Status {
		return &apiResp, fmt.Errorf("API 错误 [%s]: %s", apiResp.Code, apiResp.Message)
	}

	return &apiResp, nil
}

func (c *Client) logf(format string, args ...any) {
	if c.Verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

// encodeParams encodes url.Values with bracket notation for arrays.
// Go's default url.Values.Encode() uses key=val1&key=val2 for arrays,
// but the server expects key[]=val1&key[]=val2 (brackets format).
func encodeParams(v url.Values) string {
	if v == nil {
		return ""
	}
	var parts []string
	for key, values := range v {
		if len(values) == 1 {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(values[0]))
		} else {
			for _, val := range values {
				parts = append(parts, url.QueryEscape(key)+"[]="+url.QueryEscape(val))
			}
		}
	}
	return strings.Join(parts, "&")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// flattenParams converts a nested map to url.Values with bracket notation,
// matching the behavior of qs.stringify(data, {arrayFormat: 'brackets'}).
// e.g. {order: {0: {field: "", dir: ""}}} becomes order[0][field]=&order[0][dir]=
func flattenParams(m map[string]any, prefix string) url.Values {
	result := url.Values{}
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "[" + k + "]"
		}
		switch val := v.(type) {
		case map[string]any:
			for sk, sv := range flattenParams(val, key) {
				result[sk] = sv
			}
		case map[string]string:
			for sk, sv := range val {
				nestedKey := key + "[" + sk + "]"
				result.Set(nestedKey, sv)
			}
		case string:
			result.Set(key, val)
		case int:
			result.Set(key, fmt.Sprintf("%d", val))
		case bool:
			if val {
				result.Set(key, "true")
			} else {
				result.Set(key, "false")
			}
		case float64:
			result.Set(key, fmt.Sprintf("%v", val))
		case []any:
			for i, item := range val {
				arrKey := key + "[" + fmt.Sprintf("%d", i) + "]"
				switch iv := item.(type) {
				case map[string]string:
					for sk, sv := range iv {
						result.Set(arrKey+"["+sk+"]", sv)
					}
				default:
					result.Set(arrKey, fmt.Sprintf("%v", iv))
				}
			}
		default:
			result.Set(key, fmt.Sprintf("%v", val))
		}
	}
	return result
}
