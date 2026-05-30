package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/LucyHeres/xrxs-cli/internal/auth"
	"github.com/LucyHeres/xrxs-cli/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "登录认证管理",
		Long:  `管理 xrxs CLI 的登录认证状态。`,
	}

	cmd.AddCommand(
		newAuthLoginCommand(),
		newAuthLogoutCommand(),
		newAuthStatusCommand(),
	)
	return cmd
}

func newAuthLoginCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "登录到薪人薪事系统",
		Long: `使用用户名和密码登录薪人薪事系统，保存认证会话供后续命令使用。

示例:
  xrxs auth login --base-url https://your-instance.example.com`,
		RunE: authLoginRunE,
	}
}

func newAuthLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "退出登录，清除认证会话",
		RunE: authLogoutRunE,
	}
}

func newAuthStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "查看当前登录状态",
		RunE: authStatusRunE,
	}
}

func authLoginRunE(cmd *cobra.Command, args []string) error {
	baseURL, err := resolveBaseURL(cmd)
	if err != nil {
		return err
	}

	fmt.Print("用户名: ")
	username, err := readLine()
	if err != nil {
		return fmt.Errorf("读取用户名失败: %w", err)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	fmt.Print("密码: ")
	pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("读取密码失败: %w", err)
	}
	fmt.Println()
	pwdStr := strings.TrimSpace(string(pwdBytes))
	if pwdStr == "" {
		return fmt.Errorf("密码不能为空")
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Printf("[DRY RUN] 将登录到 %s (用户: %s)\n", baseURL, username)
		return nil
	}

	fmt.Printf("正在登录 %s ...\n", baseURL)
	session, err := auth.Login(baseURL, username, pwdStr)
	if err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	keyring, err := auth.NewKeyring()
	if err != nil {
		return fmt.Errorf("初始化密钥: %w", err)
	}

	cookiesPath := resolveCookiesPath(cmd)
	if err := os.MkdirAll(config.DefaultConfigDir(), config.DirPerm); err != nil {
		return fmt.Errorf("创建配置目录: %w", err)
	}

	if err := session.Save(cookiesPath, keyring); err != nil {
		return fmt.Errorf("保存会话: %w", err)
	}

	fmt.Println("登录成功！")
	return nil
}

func authLogoutRunE(cmd *cobra.Command, args []string) error {
	cookiesPath := resolveCookiesPath(cmd)
	if err := os.Remove(cookiesPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("清除会话失败: %w", err)
	}
	fmt.Println("已退出登录。")
	return nil
}

func authStatusRunE(cmd *cobra.Command, args []string) error {
	cookiesPath := resolveCookiesPath(cmd)
	keyring, err := auth.NewKeyring()
	if err != nil {
		return fmt.Errorf("初始化密钥: %w", err)
	}

	session, err := auth.LoadSession(cookiesPath, keyring)
	if err != nil {
		fmt.Println("未登录。请运行: xrxs auth login")
		return nil
	}

	baseURL, _ := resolveBaseURL(cmd)
	fmt.Printf("状态:    已登录\n")
	fmt.Printf("服务器:   %s\n", baseURL)
	fmt.Printf("登录时间: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05"))
	if session.IsExpired() {
		fmt.Println("警告: 会话已过期，请重新登录。")
	} else {
		fmt.Println("会话有效。")
	}
	return nil
}
