package cmd

import (
	"fmt"
	"log"
	"os"

	"subnode-ui-cli/internal/browser"

	"github.com/mxschmitt/playwright-go"
	"github.com/spf13/cobra"
)

var navigateURL string

var navigateCmd = &cobra.Command{
	Use:   "navigate [flags] --url /path",
	Short: "Open a URL and capture evidence",
	Long: `Navigate to a URL and capture a full debug bundle:
- Screenshot
- DOM HTML
- Console logs
- Network failures
- UX gate results

Example:
  ui-cli navigate --url /login --viewport mobile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if navigateURL == "" {
			return fmt.Errorf("--url is required")
		}

		return runNavigate(navigateURL)
	},
}

func init() {
	navigateCmd.Flags().StringVar(&navigateURL, "url", "", "URL to navigate to (required)")
	navigateCmd.Flags().StringArrayVar(&viewportPresets, "viewports", []string{"desktop"}, "Viewports to capture")
	navigateCmd.Flags().StringArrayVar(&colorSchemes, "schemes", []string{"light"}, "Color schemes")
}

var viewportPresets []string
var colorSchemes []string

func runNavigate(url string) error {
	// Create browser
	browserInstance, err := browser.New(headless)
	if err != nil {
		return fmt.Errorf("failed to create browser: %w", err)
	}
	defer browserInstance.Close()

	baseURL := GetBaseURL()
	fullURL := baseURL + url

	fmt.Fprintf(os.Stderr, "Navigating to %s...\n", fullURL)

	// Create run bundle
	vp := browser.SizeFromPreset(viewport)
	page, err := browserInstance.NewPage(vp, colorScheme)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	// Attach listeners before navigating
	snapshots := browser.NewSnapshots()
	browser.Attach(page, snapshots)

	// Create run
	bundle, err := browser.NewRun("navigate", headless, screenshotDir, vp, colorScheme)
	if err != nil {
		return fmt.Errorf("failed to create run: %w", err)
	}
	defer bundle.Close()
	bundle.SetPage(page)

	// Navigate
	_, err = page.Goto(fullURL)
	if err != nil {
		return fmt.Errorf("navigation failed: %w", err)
	}

	// Wait for content
	timeout := float64(5000)
	locator := page.Locator("[data-testid]")
	_ = locator.WaitFor(playwright.LocatorWaitForOptions{Timeout: &timeout})
	if err != nil {
		log.Printf("Warning: no testid found within 5s (page may not have testids)\n")
	}

	// Take screenshot
	if err := bundle.WriteScreenshot(page); err != nil {
		log.Printf("Warning: screenshot failed: %v\n", err)
	}

	// Finish run
	artifacts := bundle.FinishOrLog()

	// Output JSON to stdout
	fmt.Fprintf(os.Stdout, "{\"run_dir\":\"%s\",\"screenshot\":\"%s\"}\n",
		artifacts.RunDir, artifacts.ScreenshotPath)

	fmt.Fprintf(os.Stderr, "✓ Bundle written to %s\n", artifacts.RunDir)
	return nil
}
