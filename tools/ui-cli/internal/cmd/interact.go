package cmd

import (
	"fmt"
	"os"

	"subnode-ui-cli/internal/browser"

	"github.com/spf13/cobra"
)

var interactAction string
var interactSelector string

var interactCmd = &cobra.Command{
	Use:   "interact [action] --selector [selector]",
	Short: "Perform a single browser interaction",
	Long: `One action for poking at the UI.
Actions: click, fill, press, focus, screenshot

Example:
  ui-cli interact click --selector "[data-testid=auth-submit]"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if interactAction == "" || interactSelector == "" {
			return fmt.Errorf("--action and --selector are required")
		}

		browserInstance, err := browser.New(headless)
		if err != nil {
			return err
		}
		defer browserInstance.Close()

		page, err := browserInstance.NewPage(browser.SizeFromPreset(viewport), colorScheme)
		if err != nil {
			return err
		}

		baseURL := GetBaseURL()
		_, err = page.Goto(baseURL)
		if err != nil {
			return err
		}

		locator := page.Locator(interactSelector)

		switch interactAction {
		case "click":
			if err := locator.Click(); err != nil {
				return fmt.Errorf("click failed: %w", err)
			}
			fmt.Fprintf(os.Stdout, `{"action":"click","selector":"%s","success":true}`+"\n", interactSelector)
		case "fill":
			if len(args) == 0 {
				return fmt.Errorf("fill requires a value argument")
			}
			if err := locator.Fill(args[0]); err != nil {
				return fmt.Errorf("fill failed: %w", err)
			}
			fmt.Fprintf(os.Stdout, `{"action":"fill","selector":"%s","value":"%s","success":true}`+"\n", interactSelector, args[0])
		case "press":
			if len(args) == 0 {
				return fmt.Errorf("press requires a key argument")
			}
			if err := locator.Press(args[0]); err != nil {
				return fmt.Errorf("press failed: %w", err)
			}
			fmt.Fprintf(os.Stdout, `{"action":"press","selector":"%s","key":"%s","success":true}`+"\n", interactSelector, args[0])
		case "focus":
			if err := locator.Focus(); err != nil {
				return fmt.Errorf("focus failed: %w", err)
			}
			fmt.Fprintf(os.Stdout, `{"action":"focus","selector":"%s","success":true}`+"\n", interactSelector)
		case "screenshot":
			outFile := fmt.Sprintf("%s/interact-%s.png", screenshotDir, interactAction)
			if _, err := browser.Screenshot(page, outFile); err != nil {
				return fmt.Errorf("screenshot failed: %w", err)
			}
			fmt.Fprintf(os.Stdout, `{"action":"screenshot","selector":"%s","file":"%s"}`+"\n", interactSelector, outFile)
		default:
			return fmt.Errorf("unknown action: %s", interactAction)
		}

		return nil
	},
}

func init() {
	interactCmd.Flags().StringVar(&interactAction, "action", "", "Action to perform (click, fill, press, focus, screenshot)")
	interactCmd.Flags().StringVar(&interactSelector, "selector", "", "CSS selector to interact with")
}
