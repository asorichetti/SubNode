package cmd

import (
	"fmt"
	"os"
	"strings"

	"subnode-ui-cli/internal/browser"
	"subnode-ui-cli/internal/journeys"
	"subnode-ui-cli/internal/uxgates"

	"github.com/mxschmitt/playwright-go"
	"github.com/spf13/cobra"
)

var journeyList bool
var journeyGrep string
var journeySuite string

var journeyCmd = &cobra.Command{
	Use:   "journey [name] | --list | --suite X | --grep X",
	Short: "Run test journeys",
	Long:  `Run UI test journeys against the running application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if journeyList {
			return listJourneys()
		}
		return runJourneys(args)
	},
}

func init() {
	journeyCmd.Flags().BoolVar(&journeyList, "list", false, "List available journeys")
	journeyCmd.Flags().StringVar(&journeyGrep, "grep", "", "Filter journeys by name")
	journeyCmd.Flags().StringVar(&journeySuite, "suite", "", "Run only journeys in this suite")
}

func listJourneys() error {
	allJourneys := journeys.All()
	fmt.Fprintf(os.Stdout, "Available journeys:\n")
	for _, j := range allJourneys {
		fmt.Fprintf(os.Stdout, "  %-30s [%s] %s\n", j.Name, j.Suite, j.Desc)
	}
	return nil
}

func runJourneys(args []string) error {
	allJourneys := journeys.All()
	if len(args) > 0 {
		// Run specific journeys
		for _, name := range args {
			for _, j := range allJourneys {
				if j.Name == name {
					fmt.Fprintf(os.Stderr, "Running journey: %s\n", name)
					if err := runJourney(j); err != nil {
						return err
					}
				}
			}
		}
	} else {
		// Run all or filtered
		for _, j := range allJourneys {
			if journeyGrep != "" && !strings.Contains(strings.ToLower(j.Name), strings.ToLower(journeyGrep)) {
				continue
			}
			if journeySuite != "" && j.Suite != journeySuite {
				continue
			}
			fmt.Fprintf(os.Stderr, "Running journey: %s\n", j.Name)
			if err := runJourney(j); err != nil {
				return err
			}
		}
	}
	return nil
}

func runJourney(j journeys.Journey) error {
	// Create browser
	browserInstance, err := browser.New(headless)
	if err != nil {
		return fmt.Errorf("failed to create browser: %w", err)
	}
	defer browserInstance.Close()

	// Create page
	vp := browser.SizeFromPreset(viewport)
	page, err := browserInstance.NewPage(vp, colorScheme)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	// Attach listeners
	snapshots := browser.NewSnapshots()
	browser.Attach(page, snapshots)

	// Create recorder
	recorder := uxgates.NewRecorder(uxgates.ParseMode(uxGates))

	// Create context
	ctx := &journeys.Context{
		Page:    page,
		BaseURL: GetBaseURL(),
		Gates:   recorder,
		Logf:    func(format string, args ...any) {},
	}

	// Goto helper
	ctx.Goto = func(url string) error {
		return gotoAndGates(page, GetBaseURL()+url, snapshots, recorder)
	}

	// Run journey
	if err := j.Run(ctx); err != nil {
		return fmt.Errorf("journey %s failed: %w", j.Name, err)
	}

	// Skip gate checking for gate test journeys
	if j.Suite == "gates" {
		fmt.Fprintf(os.Stderr, "✓ Journey %s passed\n", j.Name)
		return nil
	}

	// Run console and network gates at the end
	uxgates.RunConsole(page, snapshots, recorder)
	uxgates.RunNetwork(page, snapshots, recorder)

	if recorder.Failed() && uxGates != "off" {
		fmt.Fprintf(os.Stderr, "\n⚠️ Gate failures:\n")
		for _, finding := range recorder.Findings() {
			fmt.Fprintf(os.Stderr, "  - [%s] %s: %s\n", finding.Gate, finding.Severity, finding.Message)
		}
		return fmt.Errorf("journey %s failed: gate failures detected", j.Name)
	}

	fmt.Fprintf(os.Stderr, "✓ Journey %s passed\n", j.Name)
	return nil
}

func gotoAndGates(page playwright.Page, url string, snapshots *browser.Snapshots, recorder *uxgates.Recorder) error {
	_, err := page.Goto(url)
	if err != nil {
		return err
	}

	// Run gates
	uxgates.RunPageState(page, recorder)
	uxgates.RunRender(page, recorder)
	uxgates.RunA11y(page, recorder)

	return nil
}
