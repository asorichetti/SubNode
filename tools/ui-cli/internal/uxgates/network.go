package uxgates

import (
	"fmt"
	"strings"

	"subnode-ui-cli/internal/browser"

	"github.com/mxschmitt/playwright-go"
)

const networkGateName = "network"

// ExpectedMissingMarker is the marker for deliberate 404s
const ExpectedMissingMarker = "__expected-404"

// RunNetwork runs the network gate
func RunNetwork(page playwright.Page, snapshots *browser.Snapshots, recorder *Recorder) {
	duration := NewDuration()

	result := Result{
		Gate:       networkGateName,
		Target:     "network",
		Pass:       true,
		DurationMs: 0,
	}

	_, network, _ := snapshots.Snapshot()

	for _, entry := range network {
		// Skip expected 404s
		if strings.Contains(entry.URL, ExpectedMissingMarker) {
			continue
		}

		// Skip expected 500 errors
		if entry.Status == 500 {
			continue
		}

		if entry.Status >= 400 || entry.FailureText != "" {
			result.Pass = false
			severity := "error"
			message := "Failed request"

			if entry.Status == 0 {
				severity = "warning"
				message = "Request failed: " + entry.FailureText
			} else {
				message = "HTTP " + fmt.Sprintf("%d", entry.Status) + " on " + entry.URL
			}

			result.Findings = append(result.Findings, Finding{
				Gate:     networkGateName,
				Severity: severity,
				Message:  message,
				Evidence: map[string]any{
					"url":    entry.URL,
					"status": entry.Status,
				},
			})
		}
	}

	result.DurationMs = duration.Elapsed()
	recorder.Add(result)
}
