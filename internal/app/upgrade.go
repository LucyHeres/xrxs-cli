package app

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const upgradeTimeout = 30 * time.Second

func newUpgradeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "更新到最新版本",
		Long:  "检查 GitHub Releases 获取最新版本，如有新版本则自动下载安装。",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade()
		},
	}
}

func runUpgrade() error {
	currentVersion := strings.TrimPrefix(version, "v")
	if currentVersion == "" || currentVersion == "dev" {
		currentVersion = "0.0.0"
	}

	fmt.Println("正在检查更新...")

	latest, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("获取最新版本失败: %w", err)
	}

	latestVersion := strings.TrimPrefix(latest.TagName, "v")

	if !needsUpgrade(currentVersion, latestVersion) {
		fmt.Printf("已是最新版本 v%s\n", currentVersion)
		return nil
	}

	fmt.Printf("发现新版本: v%s → v%s\n", currentVersion, latestVersion)
	fmt.Println("正在下载...")

	archiveName := fmt.Sprintf("xrxs_%s_%s-%s.tar.gz", latestVersion, runtime.GOOS, runtime.GOARCH)
	downloadURL := fmt.Sprintf("https://github.com/LucyHeres/xrxs-cli/releases/download/%s/%s",
		latest.TagName, archiveName)

	tmpDir, err := os.MkdirTemp("", "xrxs-upgrade")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, archiveName)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	binaryPath, err := extractBinary(archivePath, tmpDir)
	if err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}

	targetPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("找不到当前程序路径: %w", err)
	}

	// On macOS, the binary might be a symlink; resolve it
	if resolved, err := filepath.EvalSymlinks(targetPath); err == nil {
		targetPath = resolved
	}

	// Atomic replace: write to temp file, then rename
	tmpTarget := targetPath + ".new"
	if err := copyFile(binaryPath, tmpTarget); err != nil {
		return fmt.Errorf("安装失败: %w", err)
	}
	os.Chmod(tmpTarget, 0o755)

	if err := os.Rename(tmpTarget, targetPath); err != nil {
		// Fallback: copy+remove (cross-device on macOS)
		if err := copyFile(tmpTarget, targetPath); err != nil {
			os.Remove(tmpTarget)
			return fmt.Errorf("替换二进制失败: %w", err)
		}
		os.Remove(tmpTarget)
		os.Chmod(targetPath, 0o755)
	}

	fmt.Printf("已更新到 v%s\n", latestVersion)
	return nil
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func fetchLatestRelease() (*githubRelease, error) {
	apiURL := os.Getenv("XRXS_UPGRADE_API")
	if apiURL == "" {
		apiURL = "https://api.github.com/repos/LucyHeres/xrxs-cli/releases/latest"
	}

	client := &http.Client{Timeout: upgradeTimeout}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "xrxs-cli-upgrade")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接 GitHub 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 || resp.StatusCode == 429 {
		return nil, fmt.Errorf("GitHub API 频率限制，请稍后重试")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API 返回 HTTP %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("解析版本信息失败: %w", err)
	}
	return &release, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, destDir string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if hdr.Name == "xrxs" || filepath.Base(hdr.Name) == "xrxs" {
			dest := filepath.Join(destDir, "xrxs")
			out, err := os.Create(dest)
			if err != nil {
				return "", err
			}
			defer out.Close()
			if _, err := io.Copy(out, tr); err != nil {
				return "", err
			}
			os.Chmod(dest, 0o755)
			return dest, nil
		}
	}
	return "", fmt.Errorf("archive 中未找到 xrxs 二进制文件")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func needsUpgrade(current, latest string) bool {
	cp := parseVersion(current)
	lp := parseVersion(latest)
	for i := 0; i < 3; i++ {
		if lp[i] > cp[i] {
			return true
		}
		if lp[i] < cp[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	var parts [3]int
	v = strings.TrimPrefix(v, "v")
	if idx := strings.IndexByte(v, '-'); idx >= 0 {
		v = v[:idx]
	}
	for i, s := range strings.SplitN(v, ".", 3) {
		if i < 3 {
			fmt.Sscanf(s, "%d", &parts[i])
		}
	}
	return parts
}
