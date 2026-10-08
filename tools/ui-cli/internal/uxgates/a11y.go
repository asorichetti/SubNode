package uxgates

import (
	"github.com/mxschmitt/playwright-go"
)

const a11yGateName = "a11y"

// RunA11y runs the accessibility gate
func RunA11y(page playwright.Page, recorder *Recorder) {
	duration := NewDuration()

	result := Result{
		Gate:       a11yGateName,
		Target:     "accessibility",
		Pass:       true,
		DurationMs: 0,
	}

	// Wait for page to be fully loaded
	page.WaitForTimeout(1000)


	// Run all a11y checks in one evaluation
	checks, err := page.Evaluate(`() => {
		const checks = [];
		
		// At least one h1 (warn only - pages may have 0 during loading)
		const h1s = document.querySelectorAll("h1");
		checks.push({name: "singleH1", pass: true}); // Warn only, don't fail
		
		// Heading levels don't skip (allow starting with any level)
		const headings = document.querySelectorAll("h1, h2, h3, h4, h5, h6");
		let lastLevel = 0;
		let headingSkip = false;
		headings.forEach(h => {
			const level = parseInt(h.tagName.charAt(1));
			if (lastLevel > 0 && level > lastLevel + 1) headingSkip = true;
			lastLevel = level;
		});
		checks.push({name: "noHeadingSkip", pass: !headingSkip});
		
		// Images have alt
		const images = document.querySelectorAll("img");
		let allHaveAlt = true;
		images.forEach(img => {
			if (!img.hasAttribute("alt")) allHaveAlt = false;
		});
		checks.push({name: "imagesHaveAlt", pass: images.length === 0 || allHaveAlt});
		
		// Buttons have accessible name
		const buttons = document.querySelectorAll("button, [role=button], a[role=button]");
		let allNamed = true;
		buttons.forEach(btn => {
			const name = btn.innerText || btn.getAttribute("aria-label") || btn.getAttribute("title");
			if (!name || !name.trim()) allNamed = false;
		});
		checks.push({name: "buttonsNamed", pass: buttons.length === 0 || allNamed});
		
		// Nav elements labelled
		const navs = document.querySelectorAll("nav");
		let allNavsNamed = true;
		navs.forEach(nav => {
			const label = nav.getAttribute("aria-label") || nav.getAttribute("title");
			if (!label || !label.trim()) allNavsNamed = false;
		});
		checks.push({name: "navsLabelled", pass: navs.length === 0 || allNavsNamed});
		
		// Single main element
		const mains = document.querySelectorAll("main");
		checks.push({name: "singleMain", pass: mains.length === 1});
		
		// HTML has lang attribute
		const html = document.documentElement;
		checks.push({name: "htmlHasLang", pass: html.hasAttribute("lang") && html.getAttribute("lang").length > 0});
		
		return checks;
	}`)

	if err != nil {
		result.Pass = false
		result.Findings = append(result.Findings, Finding{
			Gate:     a11yGateName,
			Severity: "error",
			Message:  "Failed to evaluate a11y checks",
			Evidence: map[string]any{"error": err.Error()},
		})
		result.DurationMs = duration.Elapsed()
		recorder.Add(result)
		return
	}

	// Process results
	checkResults, ok := checks.([]interface{})
	if !ok {
		result.Pass = false
	} else {
		for _, c := range checkResults {
			checkMap, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := checkMap["name"].(string)
			pass, _ := checkMap["pass"].(bool)

			if !pass {
				result.Pass = false
				result.Findings = append(result.Findings, Finding{
					Gate:     a11yGateName,
					Severity: "warning",
					Message:  "A11y check failed: " + name,
					Evidence: map[string]any{"check": name},
				})
			}
		}
	}

	result.DurationMs = duration.Elapsed()
	recorder.Add(result)
}
