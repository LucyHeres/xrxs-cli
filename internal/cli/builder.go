package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"text/template"

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
	cmd := &cobra.Command{
		Use:   t.Name,
		Short: t.Description,
		RunE:  b.createRunE(t),
	}

	for i := range t.Params {
		p := &t.Params[i]
		flagName := p.Name
		switch p.Type {
		case "string":
			def := ""
			if p.Default != nil {
				def = fmt.Sprintf("%v", p.Default)
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
		pageSize := toInt(values["page-size"])
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
	switch p.Type {
	case "int":
		v, _ := cmd.Flags().GetInt(p.Name)
		return fmt.Sprintf("%d", v)
	case "string":
		v, _ := cmd.Flags().GetString(p.Name)
		return v
	case "bool":
		v, _ := cmd.Flags().GetBool(p.Name)
		return fmt.Sprintf("%t", v)
	}
	return ""
}

func getFlagValue(cmd *cobra.Command, p schema.Param) any {
	switch p.Type {
	case "int":
		v, _ := cmd.Flags().GetInt(p.Name)
		return v
	case "string":
		v, _ := cmd.Flags().GetString(p.Name)
		return v
	case "bool":
		v, _ := cmd.Flags().GetBool(p.Name)
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
