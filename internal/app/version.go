package app

import "fmt"

var (
	// Injected via ldflags at build time.
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// Version returns the CLI version string.
func Version() string {
	return version
}

// FullVersion returns the detailed version info.
func FullVersion() string {
	return fmt.Sprintf("xrxs version %s (built %s, commit %s)", version, buildTime, gitCommit)
}
