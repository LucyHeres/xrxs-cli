package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "卸载 xrxs CLI",
		Long:  "删除 xrxs 可执行文件、配置目录和 AI Skill 文件。",
		RunE: func(cmd *cobra.Command, args []string) error {
			yes, _ := cmd.Flags().GetBool("yes")
			if !yes {
				fmt.Print("确认卸载 xrxs CLI？将删除程序文件和配置 (y/N): ")
				var answer string
				fmt.Scanln(&answer)
				if strings.ToLower(answer) != "y" {
					fmt.Println("已取消。")
					return nil
				}
			}

			removed := 0

			// 1. Remove binary
			binaryPath, err := os.Executable()
			if err == nil {
				if resolved, err := filepath.EvalSymlinks(binaryPath); err == nil {
					binaryPath = resolved
				}
				if err := os.Remove(binaryPath); err == nil {
					fmt.Printf("已删除: %s\n", binaryPath)
					removed++
				}
			}

			// 2. Remove config directory
			home, _ := os.UserHomeDir()
			configDir := filepath.Join(home, ".xrxs")
			if err := os.RemoveAll(configDir); err == nil {
				fmt.Printf("已删除: %s\n", configDir)
				removed++
			}

			// 3. Remove skills from agent directories
			agentDirs := []string{
				".agents/skills/xrxs",
				".claude/skills/xrxs",
				".cursor/skills/xrxs",
				".gemini/skills/xrxs",
				".codex/skills/xrxs",
				".github/skills/xrxs",
				".windsurf/skills/xrxs",
				".augment/skills/xrxs",
				".cline/skills/xrxs",
				".amp/skills/xrxs",
				".kiro/skills/xrxs",
				".trae/skills/xrxs",
				".openclaw/skills/xrxs",
				".hermes/skills/xrxs",
				".qoder/skills/xrxs",
				".opencode/skills/xrxs",
			}
			for _, d := range agentDirs {
				skillDir := filepath.Join(home, d)
				os.RemoveAll(skillDir)
			}

			// 4. Clean PATH from shell configs
			for _, rc := range []string{".zshenv", ".zshrc", ".bashrc", ".profile"} {
				rcPath := filepath.Join(home, rc)
				if data, err := os.ReadFile(rcPath); err == nil {
					content := string(data)
					newContent := removePathLine(content, ".local/bin")
					if newContent != content {
						os.WriteFile(rcPath, []byte(newContent), 0o644)
						fmt.Printf("已清理: %s\n", rcPath)
					}
				}
			}

			fmt.Println("卸载完成。")
			return nil
		},
	}
}

func removePathLine(content, dir string) string {
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		if strings.Contains(line, "PATH") && strings.Contains(line, dir) {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
