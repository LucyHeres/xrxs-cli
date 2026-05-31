package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"unicode"

	"github.com/LucyHeres/xrxs-cli/internal/client"
	"github.com/LucyHeres/xrxs-cli/internal/output"
	"github.com/LucyHeres/xrxs-cli/internal/schema"
	"github.com/spf13/cobra"
)

// Builder constructs Cobra commands from a Schema Manifest.
type Builder struct {
	ClientFactory func(cmd *cobra.Command) (*client.Client, error)
	FormatFunc    func(cmd *cobra.Command) (output.Format, string, string)
}

// BuildCommands creates Cobra commands for all products in the manifest.
func (b *Builder) BuildCommands(m *schema.Manifest) []*cobra.Command {
	var cmds []*cobra.Command
	for i := range m.Products {
		cmds = append(cmds, b.buildProduct(&m.Products[i]))
	}
	return cmds
}

func (b *Builder) buildProduct(p *schema.Product) *cobra.Command {
	cmd := &cobra.Command{
		Use:   p.Name,
		Short: p.Description,
		Long:  p.Description,
	}
	for i := range p.Tools {
		cmd.AddCommand(b.buildTool(&p.Tools[i]))
	}
	return cmd
}

func (b *Builder) buildTool(t *schema.Tool) *cobra.Command {
	if t.IsLeaf() {
		return b.buildLeaf(t)
	}
	cmd := &cobra.Command{
		Use:   t.Name,
		Short: t.Description,
		Long:  t.Description,
	}
	for i := range t.Subtools {
		cmd.AddCommand(b.buildTool(&t.Subtools[i]))
	}
	return cmd
}

func (b *Builder) buildLeaf(t *schema.Tool) *cobra.Command {
	var runE func(*cobra.Command, []string) error
	var params []schema.Param

	if t.Pipeline != nil {
		runE = b.createPipelineRunE(t)
		params = collectPipelineParams(t.Pipeline)
	} else {
		runE = b.createRunE(t)
		params = t.Params
	}

	cmd := &cobra.Command{
		Use:   t.Name,
		Short: t.Description,
		RunE:  runE,
	}

	for i := range params {
		p := &params[i]
		flagName := toKebabCase(p.Name) // CLI flag: kebab-case; API name: camelCase
		switch p.Type {
		case "string":
			def := ""
			if p.Default != nil {
				ds := fmt.Sprintf("%v", p.Default)
				if !strings.Contains(ds, "{{") {
					def = ds
				}
			}
			if p.Short != "" {
				cmd.Flags().StringP(flagName, p.Short, def, p.Description)
			} else {
				cmd.Flags().String(flagName, def, p.Description)
			}
		case "int":
			def := 0
			if p.Default != nil {
				if n, ok := p.Default.(float64); ok {
					def = int(n)
				}
			}
			cmd.Flags().Int(flagName, def, p.Description)
		case "bool":
			def := false
			if p.Default != nil {
				if bv, ok := p.Default.(bool); ok {
					def = bv
				}
			}
			cmd.Flags().Bool(flagName, def, p.Description)
		}
		if p.Required {
			cmd.MarkFlagRequired(flagName)
		}
	}

	return cmd
}

func (b *Builder) createRunE(t *schema.Tool) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cl, err := b.ClientFactory(cmd)
		if err != nil {
			return err
		}

		format, fields, jqExpr := b.FormatFunc(cmd)

		var resp *client.Response

		switch t.Encoding {
		case "form-nested":
			bodyMap := b.buildBodyMap(t, cmd)
			resp, err = cl.PostFormJSON(context.Background(), t.Path, bodyMap)
		case "json":
			bodyMap := b.buildBodyMap(t, cmd)
			resp, err = cl.PostJSON(context.Background(), t.Path, bodyMap)
		case "form":
			params := b.buildFormParams(t, cmd)
			if t.Method == "GET" {
				resp, err = cl.Get(context.Background(), t.Path, params)
			} else {
				resp, err = cl.Post(context.Background(), t.Path, params)
			}
		default:
			params := b.buildFormParams(t, cmd)
			if t.Method == "GET" {
				resp, err = cl.Get(context.Background(), t.Path, params)
			} else {
				resp, err = cl.Post(context.Background(), t.Path, params)
			}
		}

		if err != nil {
			return err
		}

		var parsed any
		json.Unmarshal(resp.Data, &parsed)
		if t.Output.Unwrap != "" {
			parsed = unwrapByPath(parsed, t.Output.Unwrap)
		}

		return output.WriteCommandPayload(format, parsed, fields, jqExpr, t.Output.Labels)
	}
}

// --- pipeline execution ---

// createPipelineRunE returns a RunE function that executes a multi-step API pipeline.
// Steps with FanOut configured iterate over an array from a previous step.
func (b *Builder) createPipelineRunE(t *schema.Tool) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cl, err := b.ClientFactory(cmd)
		if err != nil {
			return err
		}

		format, fields, jqExpr := b.FormatFunc(cmd)
		pipelineCtx := make(map[string]any)

		for _, step := range t.Pipeline.Steps {
			if step.Condition != "" {
				ok, err := evaluateCondition(step.Condition, pipelineCtx)
				if err != nil {
					return fmt.Errorf("步骤 %s 条件表达式错误: %w", step.ID, err)
				}
				if !ok {
					continue
				}
			}

			// Fan-out: iterate over an array from a previous step
			if step.FanOut != nil {
				results, err := b.executeFanOutStep(cl, step, cmd, pipelineCtx)
				if err != nil {
					return fmt.Errorf("步骤 %s (fan-out %s) 执行失败: %w", step.ID, step.Path, err)
				}
				key := step.OutputKey
				if key == "" {
					key = step.ID
				}
				pipelineCtx[key] = results
				continue
			}

			resolvedParams := b.resolveStepParams(step, cmd, pipelineCtx)
			resolvedBody := b.resolveStepBody(step, cmd, pipelineCtx)

			resp, err := b.executeStep(cl, step, resolvedParams, resolvedBody)
			if err != nil {
				return fmt.Errorf("步骤 %s (%s) 执行失败: %w", step.ID, step.Path, err)
			}

			var parsed any
			if len(resp.Data) > 0 {
				if err := json.Unmarshal(resp.Data, &parsed); err != nil {
					return fmt.Errorf("步骤 %s 解析响应失败: %w", step.ID, err)
				}
			} else {
				parsed = map[string]any{}
			}

			if step.Unwrap != "" {
				parsed = unwrapByPath(parsed, step.Unwrap)
			}

			key := step.OutputKey
			if key == "" {
				key = step.ID
			}
			pipelineCtx[key] = parsed
		}

		lastStep := t.Pipeline.Steps[len(t.Pipeline.Steps)-1]
		lastKey := lastStep.OutputKey
		if lastKey == "" {
			lastKey = lastStep.ID
		}
		finalResult := pipelineCtx[lastKey]

		if t.Output.Unwrap != "" {
			finalResult = unwrapByPath(finalResult, t.Output.Unwrap)
		}

		return output.WriteCommandPayload(format, finalResult, fields, jqExpr, t.Output.Labels)
	}
}

// executeStep performs a single API call using the same dispatch logic as createRunE.
func (b *Builder) executeStep(cl *client.Client, step schema.PipelineStep, resolvedParams map[string]any, bodyMap map[string]any) (*client.Response, error) {
	switch step.Encoding {
	case "form-nested":
		return cl.PostFormJSON(context.Background(), step.Path, bodyMap)
	case "json":
		return cl.PostJSON(context.Background(), step.Path, bodyMap)
	case "form":
		params := buildFormParamsFromResolved(resolvedParams)
		if step.Method == "GET" {
			return cl.Get(context.Background(), step.Path, params)
		}
		return cl.Post(context.Background(), step.Path, params)
	default:
		params := buildFormParamsFromResolved(resolvedParams)
		if step.Method == "GET" {
			return cl.Get(context.Background(), step.Path, params)
		}
		return cl.Post(context.Background(), step.Path, params)
	}
}

// --- fan-out execution ---

// executeFanOutStep iterates over a source array from a previous step, calls the API for
// each element, and collects results. Concurrency is limited to avoid overwhelming the backend.
func (b *Builder) executeFanOutStep(cl *client.Client, step schema.PipelineStep, cmd *cobra.Command, pipelineCtx map[string]any) ([]any, error) {
	fo := step.FanOut

	// Get the source array
	sourceVal, ok := pipelineCtx[fo.Source]
	if !ok {
		return nil, fmt.Errorf("fan-out source %q not found in pipeline context", fo.Source)
	}
	sourceArr, ok := sourceVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("fan-out source %q is not an array (type %T)", fo.Source, sourceVal)
	}
	if len(sourceArr) == 0 {
		return []any{}, nil
	}

	concurrency := fo.Concurrency
	if concurrency <= 0 {
		concurrency = 4 // default: max 4 concurrent requests
	}

	onError := fo.OnError
	if onError == "" {
		onError = "fail"
	}

	itemKey := fo.ItemKey
	if itemKey == "" {
		itemKey = "item"
	}
	resultKey := fo.ResultKey
	if resultKey == "" {
		resultKey = "result"
	}

	// Pre-resolve params/body that don't depend on .item (flag values from CLI)
	// We'll do per-item template resolution for template defaults that reference .item.
	hasItemTemplates := b.stepHasItemTemplates(step)

	type fanOutResult struct {
		index  int
		output map[string]any
		err    error
	}

	sem := make(chan struct{}, concurrency)
	results := make([]fanOutResult, len(sourceArr))
	var wg sync.WaitGroup
	var firstErr atomic.Value

	for i, item := range sourceArr {
		if onError == "fail" && firstErr.Load() != nil {
			break
		}

		wg.Add(1)
		go func(idx int, item any) {
			defer wg.Done()

			// Skip if we already have a fatal error
			if onError == "fail" && firstErr.Load() != nil {
				return
			}

			sem <- struct{}{}
			defer func() { <-sem }()

			var resolvedParams map[string]any
			var resolvedBody map[string]any

			if hasItemTemplates {
				resolvedParams = b.resolveStepParamsForItem(step, cmd, pipelineCtx, item)
				resolvedBody = b.resolveStepBodyForItem(step, cmd, pipelineCtx, item)
			} else {
				resolvedParams = b.resolveStepParams(step, cmd, pipelineCtx)
				resolvedBody = b.resolveStepBody(step, cmd, pipelineCtx)
			}

			resp, err := b.executeStep(cl, step, resolvedParams, resolvedBody)
			if err != nil {
				if onError == "skip" {
					results[idx] = fanOutResult{index: idx, err: err}
					return
				}
				firstErr.Store(err)
				return
			}

			var parsed any
			if len(resp.Data) > 0 {
				if err := json.Unmarshal(resp.Data, &parsed); err != nil {
					if onError == "skip" {
						results[idx] = fanOutResult{index: idx, err: fmt.Errorf("解析响应失败: %w", err)}
						return
					}
					firstErr.Store(fmt.Errorf("解析响应失败: %w", err))
					return
				}
			} else {
				parsed = map[string]any{}
			}

			if step.Unwrap != "" {
				parsed = unwrapByPath(parsed, step.Unwrap)
			}

			results[idx] = fanOutResult{
				index: idx,
				output: map[string]any{
					itemKey:   item,
					resultKey: parsed,
				},
			}
		}(i, item)
	}

	wg.Wait()

	if err, _ := firstErr.Load().(error); err != nil {
		return nil, err
	}

	// Collect non-error results in order
	var out []any
	for i := range sourceArr {
		r := results[i]
		if r.err != nil {
			continue // skipped
		}
		if r.output != nil {
			out = append(out, r.output)
		}
	}

	return out, nil
}

// stepHasItemTemplates checks whether any param default or body field references .item.
func (b *Builder) stepHasItemTemplates(step schema.PipelineStep) bool {
	for _, p := range step.Params {
		if p.Default != nil {
			if strings.Contains(fmt.Sprintf("%v", p.Default), ".item") {
				return true
			}
		}
	}
	if step.Body != nil {
		data, _ := json.Marshal(step.Body)
		if strings.Contains(string(data), ".item") {
			return true
		}
	}
	return false
}

// resolveStepParamsForItem is like resolveStepParams but provides .item in template context.
func (b *Builder) resolveStepParamsForItem(step schema.PipelineStep, cmd *cobra.Command, pipelineCtx map[string]any, item any) map[string]any {
	resolved := make(map[string]any)
	for _, p := range step.Params {
		flagName := toKebabCase(p.Name)
		userChanged := cmd.Flags().Changed(flagName)

		switch p.Type {
		case "string":
			v, _ := cmd.Flags().GetString(flagName)
			if !userChanged && p.Default != nil {
				ds := fmt.Sprintf("%v", p.Default)
				if strings.Contains(ds, "{{") {
					v = resolveTemplatesWithItem(ds, pipelineCtx, item)
				}
			}
			resolved[p.APIName()] = v
		case "int":
			v, _ := cmd.Flags().GetInt(flagName)
			resolved[p.APIName()] = v
		case "bool":
			v, _ := cmd.Flags().GetBool(flagName)
			resolved[p.APIName()] = v
		}
	}
	return resolved
}

// resolveStepBodyForItem is like resolveStepBody but provides .item in template context.
func (b *Builder) resolveStepBodyForItem(step schema.PipelineStep, cmd *cobra.Command, pipelineCtx map[string]any, item any) map[string]any {
	var templateMap map[string]any
	if step.Body != nil {
		data, _ := json.Marshal(step.Body)
		json.Unmarshal(data, &templateMap)
	}
	if templateMap == nil {
		templateMap = make(map[string]any)
	}

	values := b.buildPipelineTemplateValues(step, cmd, pipelineCtx)
	values["item"] = item
	return resolveTemplatesMap(templateMap, values)
}

// resolveTemplatesWithItem resolves Go templates with .steps and .item available.
func resolveTemplatesWithItem(raw string, pipelineCtx map[string]any, item any) string {
	if !strings.Contains(raw, "{{") {
		return raw
	}
	normalized := normalizePipelineContext(pipelineCtx)
	tmpl, err := template.New("x").Parse(raw)
	if err != nil {
		return raw
	}
	var buf bytes.Buffer
	data := map[string]any{"steps": normalized, "item": item}
	if err := tmpl.Execute(&buf, data); err != nil {
		return raw
	}
	return buf.String()
}

// resolveStepParams reads flag values and resolves template defaults.
func (b *Builder) resolveStepParams(step schema.PipelineStep, cmd *cobra.Command, pipelineCtx map[string]any) map[string]any {
	resolved := make(map[string]any)
	for _, p := range step.Params {
		flagName := toKebabCase(p.Name)
		userChanged := cmd.Flags().Changed(flagName)

		switch p.Type {
		case "string":
			v, _ := cmd.Flags().GetString(flagName)
			if !userChanged && p.Default != nil {
				ds := fmt.Sprintf("%v", p.Default)
				if strings.Contains(ds, "{{") {
					v = resolvePipelineTemplates(ds, pipelineCtx)
				}
			}
			resolved[p.APIName()] = v
		case "int":
			v, _ := cmd.Flags().GetInt(flagName)
			resolved[p.APIName()] = v
		case "bool":
			v, _ := cmd.Flags().GetBool(flagName)
			resolved[p.APIName()] = v
		}
	}
	return resolved
}

// resolveStepBody builds and resolves the body template for a pipeline step.
func (b *Builder) resolveStepBody(step schema.PipelineStep, cmd *cobra.Command, pipelineCtx map[string]any) map[string]any {
	var templateMap map[string]any
	if step.Body != nil {
		data, _ := json.Marshal(step.Body)
		json.Unmarshal(data, &templateMap)
	}
	if templateMap == nil {
		templateMap = make(map[string]any)
	}

	values := b.buildPipelineTemplateValues(step, cmd, pipelineCtx)
	return resolveTemplatesMap(templateMap, values)
}

// buildPipelineTemplateValues builds the template data context merging CLI flags and pipeline step results.
func (b *Builder) buildPipelineTemplateValues(step schema.PipelineStep, cmd *cobra.Command, pipelineCtx map[string]any) map[string]any {
	values := make(map[string]any)

	normalized := normalizePipelineContext(pipelineCtx)
	values["steps"] = normalized

	for _, p := range step.Params {
		flagName := toKebabCase(p.Name)
		switch p.Type {
		case "int":
			v, _ := cmd.Flags().GetInt(flagName)
			values[p.APIName()] = v
		case "string":
			v, _ := cmd.Flags().GetString(flagName)
			values[p.APIName()] = v
		case "bool":
			v, _ := cmd.Flags().GetBool(flagName)
			values[p.APIName()] = v
		}
	}

	if _, hasPage := values["page"]; hasPage {
		page := toInt(values["page"])
		pageSize := toInt(values["pageSize"])
		if pageSize < 1 {
			pageSize = 20
		}
		if page < 1 {
			page = 1
		}
		values["start"] = (page - 1) * pageSize
		values["pageSize"] = pageSize
	}

	return values
}

// collectPipelineParams gathers all params from all pipeline steps, deduplicated by name.
func collectPipelineParams(pipe *schema.PipelineSpec) []schema.Param {
	seen := make(map[string]bool)
	var result []schema.Param
	for _, step := range pipe.Steps {
		for _, p := range step.Params {
			if !seen[p.Name] {
				seen[p.Name] = true
				result = append(result, p)
			}
		}
	}
	return result
}

// buildFormParamsFromResolved converts resolved param values to url.Values.
func buildFormParamsFromResolved(params map[string]any) url.Values {
	result := url.Values{}
	for k, v := range params {
		s := fmt.Sprintf("%v", v)
		if s != "" {
			result.Set(k, s)
		}
	}
	return result
}

// normalizePipelineContext passes through pipeline step results for template access.
// Array values remain as []interface{} — use {{index .steps.X 0}} to access items.
func normalizePipelineContext(ctx map[string]any) map[string]any {
	return ctx
}

// resolvePipelineTemplates resolves {{.steps.X.field}} and {{index .steps.X N}} templates.
func resolvePipelineTemplates(raw string, pipelineCtx map[string]any) string {
	if !strings.Contains(raw, "{{") {
		return raw
	}
	normalized := normalizePipelineContext(pipelineCtx)
	tmpl, err := template.New("x").Parse(raw)
	if err != nil {
		return raw
	}
	var buf bytes.Buffer
	data := map[string]any{"steps": normalized}
	if err := tmpl.Execute(&buf, data); err != nil {
		return raw
	}
	return buf.String()
}

// evaluateCondition evaluates a Go template condition string against the pipeline context.
func evaluateCondition(cond string, pipelineCtx map[string]any) (bool, error) {
	if strings.TrimSpace(cond) == "" {
		return true, nil
	}
	normalized := normalizePipelineContext(pipelineCtx)
	tmpl, err := template.New("cond").Parse(cond)
	if err != nil {
		return false, fmt.Errorf("parse condition: %w", err)
	}
	var buf bytes.Buffer
	data := map[string]any{"steps": normalized}
	if err := tmpl.Execute(&buf, data); err != nil {
		return false, fmt.Errorf("execute condition: %w", err)
	}
	result := strings.TrimSpace(buf.String())
	return result != "" && result != "false" && result != "0", nil
}

// buildBodyMap constructs request body from template + computed values.
func (b *Builder) buildBodyMap(t *schema.Tool, cmd *cobra.Command) map[string]any {
	var templateMap map[string]any
	if t.Body != nil {
		data, _ := json.Marshal(t.Body)
		json.Unmarshal(data, &templateMap)
	}
	if templateMap == nil {
		templateMap = make(map[string]any)
	}

	values := b.buildTemplateValues(t, cmd)
	templateMap = resolveTemplatesMap(templateMap, values)
	return templateMap
}

// buildTemplateValues extracts flag values + computed fields for template rendering.
func (b *Builder) buildTemplateValues(t *schema.Tool, cmd *cobra.Command) map[string]any {
	values := make(map[string]any)
	for _, p := range t.Params {
		values[p.APIName()] = getFlagValue(cmd, p)
	}

	// Always compute start from page+pageSize (even if using defaults)
	if _, hasPage := values["page"]; hasPage {
		page := toInt(values["page"])
		pageSize := toInt(values["pageSize"])
		if pageSize < 1 {
			pageSize = 20
		}
		if page < 1 {
			page = 1
		}
		values["start"] = (page - 1) * pageSize
		values["pageSize"] = pageSize
	}

	return values
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	default:
		return 0
	}
}

// buildFormParams builds url.Values from schema params.
func (b *Builder) buildFormParams(t *schema.Tool, cmd *cobra.Command) url.Values {
	params := url.Values{}
	for _, p := range t.Params {
		val := getParamValue(cmd, p)
		if val != "" || p.Required {
			params.Set(p.APIName(), val)
		}
	}
	return params
}

// --- template engine ---

func resolveTemplatesMap(m map[string]any, values map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		result[k] = resolveTemplatesAny(v, values)
	}
	return result
}

func resolveTemplatesAny(v any, values map[string]any) any {
	switch val := v.(type) {
	case map[string]any:
		return resolveTemplatesMap(val, values)
	case []any:
		var result []any
		for _, item := range val {
			result = append(result, resolveTemplatesAny(item, values))
		}
		return result
	case string:
		if !strings.Contains(val, "{{") {
			return val
		}
		tmpl, err := template.New("x").Parse(val)
		if err != nil {
			return val
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, values); err != nil {
			return val
		}
		return buf.String()
	default:
		return v
	}
}

// --- helpers ---

func getParamValue(cmd *cobra.Command, p schema.Param) string {
	flagName := toKebabCase(p.Name)
	switch p.Type {
	case "int":
		v, _ := cmd.Flags().GetInt(flagName)
		return fmt.Sprintf("%d", v)
	case "string":
		v, _ := cmd.Flags().GetString(flagName)
		return v
	case "bool":
		v, _ := cmd.Flags().GetBool(flagName)
		return fmt.Sprintf("%t", v)
	}
	return ""
}

func getFlagValue(cmd *cobra.Command, p schema.Param) any {
	flagName := toKebabCase(p.Name)
	switch p.Type {
	case "int":
		v, _ := cmd.Flags().GetInt(flagName)
		return v
	case "string":
		v, _ := cmd.Flags().GetString(flagName)
		return v
	case "bool":
		v, _ := cmd.Flags().GetBool(flagName)
		return v
	}
	return nil
}

func getIntValue(cmd *cobra.Command, name string) (int, bool) {
	if cmd.Flags().Changed(name) {
		v, _ := cmd.Flags().GetInt(name)
		return v, true
	}
	return 0, false
}

// toKebabCase converts a camelCase string to kebab-case for use as a CLI flag name.
// e.g. "flowStepId" -> "flow-step-id", "sid" -> "sid".
func toKebabCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func unwrapByPath(v any, path string) any {
	parts := strings.Split(path, ".")
	for _, part := range parts {
		m, ok := v.(map[string]any)
		if !ok {
			return v
		}
		child, exists := m[part]
		if !exists {
			return v
		}
		v = child
	}
	return v
}
