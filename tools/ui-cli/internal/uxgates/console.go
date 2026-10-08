package uxgates

import (
	"strings"

	"subnode-ui-cli/internal/browser"

	"github.com/mxschmitt/playwright-go"
)

const consoleGateName = "console"

// RunConsole runs the console gate
func RunConsole(page playwright.Page, snapshots *browser.Snapshots, recorder *Recorder) {
	duration := NewDuration()

	result := Result{
		Gate:       consoleGateName,
		Target:     "console",
		Pass:       true,
		DurationMs: 0,
	}

	_, _, errors := snapshots.Snapshot()

	for _, errMsg := range errors {
		// Skip expected 500 errors (e.g., from auth transitions)
		if strings.Contains(errMsg, "500") {
			continue
		}
		result.Pass = false
		result.Findings = append(result.Findings, Finding{
			Gate:     consoleGateName,
			Severity: "error",
			Message:  "Page error: " + errMsg,
		})
	}

	result.DurationMs = duration.Elapsed()
	recorder.Add(result)
}
