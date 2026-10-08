package cmd

import (
	"fmt"
	"os"

	"subnode-ui-cli/internal/browser"

	"github.com/spf13/cobra"
)

var screenshotViewports []string
var screenshotSchemes []string

var screenshotCmd = &cobra.Command{
	Use:   "screenshot --viewports a,b --schemes light,dark",
	Short: "Take screenshots in multiple viewports and color schemes",
	Long:  `Capture screenshots across a matrix of viewports and color schemes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("URL required")
		}
		url := args[0]

		if len(screenshotViewports) == 0 {
			screenshotViewports = []string{"desktop"}
		}
		if len(screenshotSchemes) == 0 {
			screenshotSchemes = []string{"light"}
		}

		baseURL := GetBaseURL()
		fullURL := baseURL + url

		fmt.Fprintf(os.Stderr, "Screenshots: %s\n", fullURL)

		for _, vpName := range screenshotViewports {
			for _, scheme := range screenshotSchemes {
				vp := browser.SizeFromPreset(vpName)
				browserInstance, err := browser.New(headless)
				if err != nil {
					return err
				}

				page, err := browserInstance.NewPage(vp, scheme)
				if err != nil {
					browserInstance.Close()
					return err
				}

				_, err = page.Goto(fullURL)
				if err != nil {
					browserInstance.Close()
					return err
				}

				// Wait for load
				if err := page.WaitForLoadState(); err != nil {
					browserInstance.Close()
					return err
				}

				// Create output filename
				outFile := fmt.Sprintf("%s/%s-%s-%s.png", screenshotDir, "screenshot", vpName, scheme)
				if _, err := browser.Screenshot(page, outFile); err != nil {
					fmt.Fprintf(os.Stderr, "  ✗ %s/%s: %v\n", vpName, scheme, err)
				} else {
					fmt.Fprintf(os.Stderr, "  ✓ %s/%s\n", vpName, scheme)
				}

				browserInstance.Close()
			}
		}

		return nil
	},
}

func init() {
	screenshotCmd.Flags().StringArrayVar(&screenshotViewports, "viewports", []string{"desktop"}, "Viewports to capture")
	screenshotCmd.Flags().StringArrayVar(&screenshotSchemes, "schemes", []string{"light"}, "Color schemes")
}
