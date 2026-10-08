package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	env           string
	baseURL       string
	headless      bool
	screenshotDir string
	uxGates       string
	viewport      string
	colorScheme   string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "ui-cli",
	Short: "SubNode UI test CLI",
	Long: `A browser-automation CLI that proves UI changes actually work.

Generates evidence bundles with screenshots, DOM snapshots,
console logs, and network failures for every command run.`,
}

// Execute adds all child commands and runs root
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&env, "env", "local", "Environment (local, preview)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Override base URL")
	rootCmd.PersistentFlags().BoolVar(&headless, "headless", computeHeadless(), "Run in headless mode")
	rootCmd.PersistentFlags().StringVar(&screenshotDir, "screenshot-dir", "./screenshots", "Directory for screenshots and bundles")
	rootCmd.PersistentFlags().StringVar(&uxGates, "ux-gates", "warn", "UX gate mode: off, warn, fail")
	rootCmd.PersistentFlags().StringVar(&viewport, "viewport", "desktop", "Viewport preset: desktop, tablet, mobile")
	rootCmd.PersistentFlags().StringVar(&colorScheme, "color-scheme", "", "Color scheme: light, dark")

	rootCmd.AddCommand(navigateCmd)
	rootCmd.AddCommand(journeyCmd)
	rootCmd.AddCommand(smokeCmd)
	rootCmd.AddCommand(screenshotCmd)
	rootCmd.AddCommand(interactCmd)
}

// computeHeadless computes the default headless mode based on environment
func computeHeadless() bool {
	if os.Getenv("CI") != "" {
		return true
	}
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		return false
	}
	return runtime.GOOS != "darwin" // macOS draws via Quartz, sets neither
}

// getDefaultBaseURL returns the default base URL based on environment
func getDefaultBaseURL() string {
	switch env {
	case "preview":
		return os.Getenv("SUBNODE_PREVIEW_URL")
	case "local":
		return "http://localhost:3000"
	default:
		return os.Getenv("SUBNODE_BASE_URL")
	}
}

// GetBaseURL returns the effective base URL
func GetBaseURL() string {
	if baseURL != "" {
		return baseURL
	}
	return getDefaultBaseURL()
}

// CheckPlaywright ensures Playwright browsers are installed
func CheckPlaywright() error {
	cmd := exec.Command("npx", "--yes", "playwright", "install", "chromium")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
