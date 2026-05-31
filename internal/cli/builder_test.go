package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/LucyHeres/xrxs-cli/internal/auth"
	"github.com/LucyHeres/xrxs-cli/internal/client"
	"github.com/LucyHeres/xrxs-cli/internal/output"
	"github.com/LucyHeres/xrxs-cli/internal/schema"
	"github.com/spf13/cobra"
)

func TestIsLeaf_Pipeline(t *testing.T) {
	tool := schema.Tool{
		Name: "test",
		Pipeline: &schema.PipelineSpec{
			Steps: []schema.PipelineStep{
				{ID: "s1", Path: "/api/test", Method: "GET", Encoding: "form"},
			},
		},
	}
	if !tool.IsLeaf() {
		t.Error("tool with Pipeline should be a leaf")
	}
}

func TestIsLeaf_SingleAPI(t *testing.T) {
	tool := schema.Tool{
		Name:   "test",
		Path:   "/api/test",
		Method: "GET",
	}
	if !tool.IsLeaf() {
		t.Error("tool with Path should be a leaf")
	}
}

func TestIsLeaf_Group(t *testing.T) {
	tool := schema.Tool{
		Name:     "group",
		Subtools: []schema.Tool{{Name: "sub", Path: "/api/test", Method: "GET"}},
	}
	if tool.IsLeaf() {
		t.Error("tool with Subtools should not be a leaf")
	}
}

func TestCollectPipelineParams(t *testing.T) {
	pipe := &schema.PipelineSpec{
		Steps: []schema.PipelineStep{
			{
				ID: "s1",
				Params: []schema.Param{
					{Name: "groupId", Type: "string"},
					{Name: "page", Type: "int"},
				},
			},
			{
				ID: "s2",
				Params: []schema.Param{
					{Name: "groupId", Type: "string"}, // duplicate
					{Name: "sid", Type: "string"},
				},
			},
		},
	}

	params := collectPipelineParams(pipe)
	if len(params) != 3 {
		t.Fatalf("expected 3 params (deduplicated), got %d", len(params))
	}

	names := make(map[string]bool)
	for _, p := range params {
		names[p.Name] = true
	}
	if !names["groupId"] || !names["page"] || !names["sid"] {
		t.Errorf("unexpected param names: %v", names)
	}
}

func TestCollectPipelineParams_Empty(t *testing.T) {
	pipe := &schema.PipelineSpec{Steps: []schema.PipelineStep{}}
	params := collectPipelineParams(pipe)
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestNormalizePipelineContext_PassThrough(t *testing.T) {
	raw := `[{"groupId": 5, "groupName": "考勤"}]`
	var parsed any
	json.Unmarshal([]byte(raw), &parsed)

	ctx := map[string]any{"groups": parsed}
	result := normalizePipelineContext(ctx)

	// Pass-through: result should have the same raw data
	arr, ok := result["groups"].([]interface{})
	if !ok {
		t.Fatal("groups should be []interface{} (pass-through)")
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 element, got %d", len(arr))
	}
}

func TestResolvePipelineTemplates_WithIndex(t *testing.T) {
	raw := `[{"groupId": 5, "groupName": "考勤"}]`
	var parsed any
	json.Unmarshal([]byte(raw), &parsed)

	ctx := map[string]any{"groups": parsed}

	result := resolvePipelineTemplates("{{(index .steps.groups 0).groupId}}", ctx)
	if result != "5" {
		t.Errorf("expected '5', got %q", result)
	}
}

func TestResolvePipelineTemplates_NoTemplate(t *testing.T) {
	result := resolvePipelineTemplates("plain-value", nil)
	if result != "plain-value" {
		t.Errorf("expected 'plain-value', got %q", result)
	}
}

func TestResolvePipelineTemplates_FieldAccess(t *testing.T) {
	raw := `{"name": "请假", "flowType": 4}`
	var parsed any
	json.Unmarshal([]byte(raw), &parsed)

	ctx := map[string]any{"types": parsed}

	result := resolvePipelineTemplates("{{.steps.types.name}}", ctx)
	if result != "请假" {
		t.Errorf("expected '请假', got %q", result)
	}
}

func TestEvaluateCondition_Truthy(t *testing.T) {
	ok, err := evaluateCondition("true", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("'true' should be truthy")
	}
}

func TestEvaluateCondition_Falsy(t *testing.T) {
	ok, err := evaluateCondition("false", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("'false' should be falsy")
	}
}

func TestEvaluateCondition_Empty(t *testing.T) {
	ok, err := evaluateCondition("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("empty condition should be truthy (skip nothing)")
	}
}

func TestBuildFormParamsFromResolved(t *testing.T) {
	params := map[string]any{
		"groupId": "5",
		"page":    1,
	}
	result := buildFormParamsFromResolved(params)
	if result.Get("groupId") != "5" {
		t.Errorf("expected groupId=5, got %s", result.Get("groupId"))
	}
	if result.Get("page") != "1" {
		t.Errorf("expected page=1, got %s", result.Get("page"))
	}
}

// --- integration tests with mock HTTP server ---

func TestPipeline_FullExecution(t *testing.T) {
	// Step counters to verify both API calls are made
	var step1Called, step2Called bool
	var step2GroupID string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.URL.Path {
		case "/approve/service/ajax-get-flow-group":
			step1Called = true
			// Return groups list — simulate real API response
			json.NewEncoder(w).Encode(map[string]any{
				"code":    "0",
				"status":  true,
				"message": "",
				"data": []map[string]any{
					{"groupId": 5, "groupName": "考勤", "orderNum": 1},
					{"groupId": 7, "groupName": "员工", "orderNum": 2},
				},
			})
		case "/approve/service/ajax-get-flow-setting":
			step2Called = true
			step2GroupID = r.URL.Query().Get("groupId")
			// Return types for the given group — simulate real API response
			json.NewEncoder(w).Encode(map[string]any{
				"code":    "0",
				"status":  true,
				"message": "",
				"data": []map[string]any{
					{"flowSettingId": 101, "flowType": 1, "name": "请假", "openStatus": 1, "type": 1, "complete": true},
					{"flowSettingId": 102, "flowType": 2, "name": "加班", "openStatus": 1, "type": 1, "complete": true},
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	// Build a pipeline tool that mirrors the real list-all-types schema
	pipelineTool := &schema.Tool{
		Name:        "list-all-types",
		Description: "列出所有审批分组下的类型",
		Pipeline: &schema.PipelineSpec{
			Steps: []schema.PipelineStep{
				{
					ID:        "fetchGroups",
					Path:      "/approve/service/ajax-get-flow-group",
					Method:    "GET",
					Encoding:  "form",
					OutputKey: "groups",
				},
				{
					ID:       "fetchTypes",
					Path:     "/approve/service/ajax-get-flow-setting",
					Method:   "GET",
					Encoding: "form",
					Params: []schema.Param{
						{
							Name:    "groupId",
							Type:    "string",
							Default: "{{(index .steps.groups 0).groupId}}",
						},
					},
					OutputKey: "types",
				},
			},
		},
		Output: schema.OutputSpec{
			Labels: map[string]string{
				"flowSettingId": "类型ID",
				"name":          "名称",
			},
		},
	}

	var capturedOutput bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Build the cobra command
	b := &Builder{
		ClientFactory: func(cmd *cobra.Command) (*client.Client, error) {
			return &client.Client{
				BaseURL:    srv.URL,
				Session:    &auth.Session{},
				HTTPClient: &http.Client{},
			}, nil
		},
		FormatFunc: func(cmd *cobra.Command) (output.Format, string, string) {
			return output.FormatJSON, "", ""
		},
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(b.buildLeaf(pipelineTool))
	rootCmd.SetArgs([]string{"list-all-types"})

	// Execute the command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("pipeline execution failed: %v", err)
	}

	w.Close()
	capturedOutput.ReadFrom(r)
	os.Stdout = oldStdout

	// Verify both steps were called
	if !step1Called {
		t.Error("step 1 (fetchGroups) was not called")
	}
	if !step2Called {
		t.Error("step 2 (fetchTypes) was not called")
	}

	// Verify step 2 used the groupId from step 1's response (template resolution)
	if step2GroupID != "5" {
		t.Errorf("step 2 should use groupId=5 from step 1, got %q", step2GroupID)
	}

	// Verify the output contains types data
	var outputData []map[string]any
	if err := json.Unmarshal(capturedOutput.Bytes(), &outputData); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, capturedOutput.String())
	}
	if len(outputData) != 2 {
		t.Fatalf("expected 2 types in output, got %d", len(outputData))
	}
	if outputData[0]["name"] != "请假" {
		t.Errorf("expected first type name '请假', got %v", outputData[0]["name"])
	}
	if outputData[1]["flowSettingId"] != float64(102) {
		t.Errorf("expected second type flowSettingId 102, got %v", outputData[1]["flowSettingId"])
	}
}

func TestPipeline_UserOverridesDefault(t *testing.T) {
	var step2GroupID string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.URL.Path {
		case "/approve/service/ajax-get-flow-group":
			json.NewEncoder(w).Encode(map[string]any{
				"code":   "0",
				"status": true,
				"data": []map[string]any{
					{"groupId": 5, "groupName": "考勤"},
				},
			})
		case "/approve/service/ajax-get-flow-setting":
			step2GroupID = r.URL.Query().Get("groupId")
			json.NewEncoder(w).Encode(map[string]any{
				"code":   "0",
				"status": true,
				"data":   []map[string]any{},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	pipelineTool := &schema.Tool{
		Name: "test-override",
		Pipeline: &schema.PipelineSpec{
			Steps: []schema.PipelineStep{
				{ID: "s1", Path: "/approve/service/ajax-get-flow-group", Method: "GET", Encoding: "form", OutputKey: "groups"},
				{
					ID: "s2", Path: "/approve/service/ajax-get-flow-setting", Method: "GET", Encoding: "form",
					Params:    []schema.Param{{Name: "groupId", Type: "string", Default: "{{(index .steps.s1 0).groupId}}"}},
					OutputKey: "types",
				},
			},
		},
		Output: schema.OutputSpec{},
	}

	b := &Builder{
		ClientFactory: func(cmd *cobra.Command) (*client.Client, error) {
			return &client.Client{BaseURL: srv.URL, Session: &auth.Session{}, HTTPClient: &http.Client{}}, nil
		},
		FormatFunc: func(cmd *cobra.Command) (output.Format, string, string) {
			return output.FormatJSON, "", ""
		},
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(b.buildLeaf(pipelineTool))
	rootCmd.SetArgs([]string{"test-override", "--group-id", "99"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("pipeline execution failed: %v", err)
	}

	// Verify user override took effect, not the template default
	if step2GroupID != "99" {
		t.Errorf("step 2 should use user-provided groupId=99, got %q", step2GroupID)
	}
}

func TestPipeline_ConditionalStep(t *testing.T) {
	var step2Called bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.URL.Path {
		case "/api/step1":
			json.NewEncoder(w).Encode(map[string]any{
				"code": "0", "status": true,
				"data": map[string]any{"skippable": true},
			})
		case "/api/step2":
			step2Called = true
			json.NewEncoder(w).Encode(map[string]any{"code": "0", "status": true, "data": "should not run"})
		case "/api/step3":
			json.NewEncoder(w).Encode(map[string]any{"code": "0", "status": true, "data": "final"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	pipelineTool := &schema.Tool{
		Name: "test-condition",
		Pipeline: &schema.PipelineSpec{
			Steps: []schema.PipelineStep{
				{ID: "s1", Path: "/api/step1", Method: "GET", Encoding: "form", OutputKey: "s1"},
				{ID: "s2", Path: "/api/step2", Method: "GET", Encoding: "form", OutputKey: "s2", Condition: "false"},
				{ID: "s3", Path: "/api/step3", Method: "GET", Encoding: "form", OutputKey: "s3"},
			},
		},
		Output: schema.OutputSpec{},
	}

	b := &Builder{
		ClientFactory: func(cmd *cobra.Command) (*client.Client, error) {
			return &client.Client{BaseURL: srv.URL, Session: &auth.Session{}, HTTPClient: &http.Client{}}, nil
		},
		FormatFunc: func(cmd *cobra.Command) (output.Format, string, string) {
			return output.FormatJSON, "", ""
		},
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(b.buildLeaf(pipelineTool))
	rootCmd.SetArgs([]string{"test-condition"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("pipeline execution failed: %v", err)
	}

	if step2Called {
		t.Error("step 2 should have been skipped (condition=false)")
	}
}

// --- fan-out tests ---

func TestPipeline_FanOut_AllGroups(t *testing.T) {
	var mu sync.Mutex
	groupCalls := make(map[string]int)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.URL.Path {
		case "/approve/service/ajax-get-flow-group":
			json.NewEncoder(w).Encode(map[string]any{
				"code":   "0",
				"status": true,
				"data": []map[string]any{
					{"groupId": 5, "groupName": "考勤"},
					{"groupId": 7, "groupName": "员工"},
					{"groupId": 8, "groupName": "薪酬"},
				},
			})
		case "/approve/service/ajax-get-flow-setting":
			gid := r.URL.Query().Get("groupId")
			mu.Lock()
			groupCalls[gid]++
			mu.Unlock()
			json.NewEncoder(w).Encode(map[string]any{
				"code":   "0",
				"status": true,
				"data":   []map[string]any{{"flowSettingId": 1, "name": "审批-" + gid}},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	pipelineTool := &schema.Tool{
		Name: "list-all",
		Pipeline: &schema.PipelineSpec{
			Steps: []schema.PipelineStep{
				{ID: "fetchGroups", Path: "/approve/service/ajax-get-flow-group", Method: "GET", Encoding: "form", OutputKey: "groups"},
				{
					ID: "fetchTypes", Path: "/approve/service/ajax-get-flow-setting", Method: "GET", Encoding: "form",
					Params:  []schema.Param{{Name: "groupId", Type: "string", Default: "{{.item.groupId}}"}},
					FanOut: &schema.FanOutSpec{Source: "groups", Concurrency: 4, OnError: "fail", ItemKey: "group", ResultKey: "types"},
				},
			},
		},
		Output: schema.OutputSpec{},
	}

	var capturedOutput bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	b := &Builder{
		ClientFactory: func(cmd *cobra.Command) (*client.Client, error) {
			return &client.Client{BaseURL: srv.URL, Session: &auth.Session{}, HTTPClient: &http.Client{}}, nil
		},
		FormatFunc: func(cmd *cobra.Command) (output.Format, string, string) {
			return output.FormatJSON, "", ""
		},
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(b.buildLeaf(pipelineTool))
	rootCmd.SetArgs([]string{"list-all"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("fan-out pipeline failed: %v", err)
	}

	w.Close()
	capturedOutput.ReadFrom(r)
	os.Stdout = oldStdout

	// Verify all 3 groups were called
	if len(groupCalls) != 3 {
		t.Errorf("expected 3 group calls, got %d: %v", len(groupCalls), groupCalls)
	}
	for _, gid := range []string{"5", "7", "8"} {
		if groupCalls[gid] != 1 {
			t.Errorf("group %s should be called once, got %d", gid, groupCalls[gid])
		}
	}

	// Verify output: [{group: {groupId:...}, types: [...]}, ...]
	var outputData []map[string]any
	if err := json.Unmarshal(capturedOutput.Bytes(), &outputData); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, capturedOutput.String())
	}
	if len(outputData) != 3 {
		t.Fatalf("expected 3 items in output, got %d", len(outputData))
	}
	for i, item := range outputData {
		if item["group"] == nil {
			t.Errorf("item %d: missing 'group' key", i)
		}
		if item["types"] == nil {
			t.Errorf("item %d: missing 'types' key", i)
		}
	}
}

func TestPipeline_FanOut_SkipErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.URL.Path {
		case "/api/groups":
			json.NewEncoder(w).Encode(map[string]any{
				"code": "0", "status": true,
				"data": []map[string]any{
					{"id": 1, "name": "good"},
					{"id": 2, "name": "bad"},
					{"id": 3, "name": "good2"},
				},
			})
		case "/api/types":
			gid := r.URL.Query().Get("id")
			if gid == "2" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]any{
				"code": "0", "status": true,
				"data": []map[string]any{{"name": "type-of-" + gid}},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	pipelineTool := &schema.Tool{
		Name: "fanout-skip",
		Pipeline: &schema.PipelineSpec{
			Steps: []schema.PipelineStep{
				{ID: "s1", Path: "/api/groups", Method: "GET", Encoding: "form", OutputKey: "groups"},
				{
					ID: "s2", Path: "/api/types", Method: "GET", Encoding: "form",
					Params:  []schema.Param{{Name: "id", Type: "string", Default: "{{.item.id}}"}},
					FanOut: &schema.FanOutSpec{Source: "groups", Concurrency: 2, OnError: "skip", ItemKey: "src", ResultKey: "data"},
				},
			},
		},
		Output: schema.OutputSpec{},
	}

	var capturedOutput bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	b := &Builder{
		ClientFactory: func(cmd *cobra.Command) (*client.Client, error) {
			return &client.Client{BaseURL: srv.URL, Session: &auth.Session{}, HTTPClient: &http.Client{}}, nil
		},
		FormatFunc: func(cmd *cobra.Command) (output.Format, string, string) {
			return output.FormatJSON, "", ""
		},
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(b.buildLeaf(pipelineTool))
	rootCmd.SetArgs([]string{"fanout-skip"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("fan-out with skip should not fail: %v", err)
	}

	w.Close()
	capturedOutput.ReadFrom(r)
	os.Stdout = oldStdout

	var outputData []map[string]any
	if err := json.Unmarshal(capturedOutput.Bytes(), &outputData); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}
	if len(outputData) != 2 {
		t.Fatalf("expected 2 results (1 skipped), got %d", len(outputData))
	}
}

func TestPipeline_FanOut_EmptySource(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]any{
			"code": "0", "status": true, "data": []map[string]any{},
		})
	}))
	defer srv.Close()

	pipelineTool := &schema.Tool{
		Name: "fanout-empty",
		Pipeline: &schema.PipelineSpec{
			Steps: []schema.PipelineStep{
				{ID: "s1", Path: "/api/groups", Method: "GET", Encoding: "form", OutputKey: "groups"},
				{
					ID: "s2", Path: "/api/types", Method: "GET", Encoding: "form",
					FanOut: &schema.FanOutSpec{Source: "groups", Concurrency: 2},
				},
			},
		},
		Output: schema.OutputSpec{},
	}

	var capturedOutput bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	b := &Builder{
		ClientFactory: func(cmd *cobra.Command) (*client.Client, error) {
			return &client.Client{BaseURL: srv.URL, Session: &auth.Session{}, HTTPClient: &http.Client{}}, nil
		},
		FormatFunc: func(cmd *cobra.Command) (output.Format, string, string) {
			return output.FormatJSON, "", ""
		},
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(b.buildLeaf(pipelineTool))
	rootCmd.SetArgs([]string{"fanout-empty"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("fan-out with empty source should not fail: %v", err)
	}

	w.Close()
	capturedOutput.ReadFrom(r)
	os.Stdout = oldStdout

	var outputData []any
	json.Unmarshal(capturedOutput.Bytes(), &outputData)
	if len(outputData) != 0 {
		t.Errorf("expected empty array, got %v", outputData)
	}
}
