package uxgates

import (
	"github.com/mxschmitt/playwright-go"
)

const pagestateGateName = "pagestate"

// RunPageState runs the pagestate gate
func RunPageState(page playwright.Page, recorder *Recorder) {
	duration := NewDuration()

	result := Result{
		Gate:       pagestateGateName,
		Target:     "page-content",
		Pass:       true,
		DurationMs: 0,
	}

	// Try to find content marker within timeout
	timeout := float64(3000)

	locator := page.Locator("[data-testid=\"page-content\"]")
	err := locator.WaitFor(playwright.LocatorWaitForOptions{Timeout: &timeout})
	if err != nil {
		// Try main element as fallback
		mainLocator := page.Locator("main")
		err = mainLocator.WaitFor(playwright.LocatorWaitForOptions{Timeout: &timeout})
		if err != nil {
			result.Pass = false
			result.Findings = append(result.Findings, Finding{
				Gate:     pagestateGateName,
				Severity: "error",
				Message:  "No content marker found within timeout",
				Evidence: map[string]any{
					"timeout_ms": timeout,
				},
			})
		}
	}

	result.DurationMs = duration.Elapsed()

	if err := recorder.Add(result); err != nil {
		// Gate failed
	}
}
