package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LucyHeres/xrxs-cli/internal/auth"
	"github.com/LucyHeres/xrxs-cli/internal/cli"
	"github.com/LucyHeres/xrxs-cli/internal/client"
	"github.com/LucyHeres/xrxs-cli/internal/output"
	"github.com/LucyHeres/xrxs-cli/internal/schema"
	"github.com/LucyHeres/xrxs-cli/pkg/config"
	"github.com/spf13/cobra"
)

func loadManifest() (*schema.Manifest, error) {
	// 1. Env override: XRXS_SCHEMA_DIR
	if d := os.Getenv("XRXS_SCHEMA_DIR"); d != "" {
		return schema.LoadAllManifests(d)
	}

	// 2. Try embedded schemas (production binary)
	if m, err := schema.LoadFromEmbed(); err == nil && len(m.Products) > 0 {
		return m, nil
	}

	// 3. Development fallback: look relative to cwd
	dir, _ := os.Getwd()
	for {
		schemaDir := filepath.Join(dir, "config", "schemas")
		if _, err := os.Stat(schemaDir); err == nil {
			return schema.LoadAllManifests(schemaDir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return nil, fmt.Errorf("未找到 Schema 文件")
}

// Execute is the main entry point, called from cmd/main.go.
func Execute() int {
	root := newRootCommand()
	err := root.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		return 1
	}
	return 0
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "xrxs",
		Short: "薪人薪事 HR SaaS 命令行工具",
		Long: `xrxs 是薪人薪事 HR SaaS 系统的命令行工具，支持审批、人事等模块的操作。

命令通过 Schema 文件自动生成，添加新模块只需编写 JSON 配置。

示例:
  xrxs auth login --base-url https://your-instance.example.com
  xrxs approval list search -f table
  xrxs approval list search --status 0 -f table
  xrxs approval detail get --sid 12345 --fields employeeName,statusName`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().String("config", config.DefaultConfigFile(), "配置文件路径")
	root.PersistentFlags().String("base-url", "", "API 服务地址")
	root.PersistentFlags().StringP("format", "f", "json", "输出格式: json|table|raw")
	root.PersistentFlags().String("jq", "", "jq 表达式过滤输出")
	root.PersistentFlags().String("fields", "", "只输出指定字段 (逗号分隔)")
	root.PersistentFlags().BoolP("verbose", "v", false, "显示详细请求日志")
	root.PersistentFlags().BoolP("yes", "y", false, "跳过确认提示")
	root.PersistentFlags().Bool("dry-run", false, "预览操作但不执行")

	root.AddCommand(newAuthCommand(), newVersionCommand(), newUpgradeCommand(), newUninstallCommand())

	// Load schema and build dynamic commands
	manifest, err := loadManifest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 加载 Schema 失败: %v\n", err)
	} else {
		builder := &cli.Builder{
			ClientFactory: makeClient,
			FormatFunc:    resolveFormatOpts,
		}
		for _, cmd := range builder.BuildCommands(manifest) {
			root.AddCommand(cmd)
		}
	}

	return root
}

// makeClient creates an API client from stored session and cmd flags.
func makeClient(cmd *cobra.Command) (*client.Client, error) {
	session, err := loadSession(cmd)
	if err != nil {
		return nil, err
	}

	if url, _ := cmd.Flags().GetString("base-url"); url != "" {
		session.BaseURL = strings.TrimRight(url, "/")
	} else if envURL := os.Getenv("XRXS_BASE_URL"); envURL != "" {
		session.BaseURL = strings.TrimRight(envURL, "/")
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	return client.NewClient(session, verbose, dryRun), nil
}

// resolveFormatOpts extracts format, fields, jq from command flags.
func resolveFormatOpts(cmd *cobra.Command) (output.Format, string, string) {
	formatStr, _ := cmd.Flags().GetString("format")
	format := output.ResolveFormat(formatStr, output.FormatJSON)
	fields, _ := cmd.Flags().GetString("fields")
	jqExpr, _ := cmd.Flags().GetString("jq")
	return format, fields, jqExpr
}

// --- shared helpers for auth subcommands ---

func resolveBaseURL(cmd *cobra.Command) (string, error) {
	if url, _ := cmd.Flags().GetString("base-url"); url != "" {
		return strings.TrimRight(url, "/"), nil
	}
	if envURL := os.Getenv("XRXS_BASE_URL"); envURL != "" {
		return strings.TrimRight(envURL, "/"), nil
	}

	configPath := resolveConfigPath(cmd)
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return "", fmt.Errorf("无法加载配置: %w", err)
	}
	if cfg.BaseURL == "" {
		return "", fmt.Errorf("请指定 --base-url 或设置 XRXS_BASE_URL 环境变量")
	}
	return strings.TrimRight(cfg.BaseURL, "/"), nil
}

func resolveConfigPath(cmd *cobra.Command) string {
	if path, _ := cmd.Flags().GetString("config"); path != "" {
		return path
	}
	return config.DefaultConfigFile()
}

func resolveCookiesPath(cmd *cobra.Command) string {
	configPath := resolveConfigPath(cmd)
	dir := config.DefaultConfigDir()
	if strings.Contains(configPath, "/") {
		lastSep := strings.LastIndex(configPath, "/")
		if lastSep >= 0 {
			dir = configPath[:lastSep]
		}
	}
	return dir + "/" + config.CookiesFileName
}

func loadSession(cmd *cobra.Command) (*auth.Session, error) {
	cookiesPath := resolveCookiesPath(cmd)
	keyring, err := auth.NewKeyring()
	if err != nil {
		return nil, fmt.Errorf("初始化密钥: %w", err)
	}

	session, err := auth.LoadSession(cookiesPath, keyring)
	if err != nil {
		return nil, fmt.Errorf("未登录，请先运行: xrxs auth login --base-url <地址>")
	}
	if session.IsExpired() {
		return nil, fmt.Errorf("会话已过期，请重新登录: xrxs auth login --base-url <地址>")
	}
	return session, nil
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(FullVersion())
		},
	}
}
