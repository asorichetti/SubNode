package uxgates

import (
	"github.com/mxschmitt/playwright-go"
)

const renderGateName = "render"

// RunRender runs the render gate to check CSS is applied
func RunRender(page playwright.Page, recorder *Recorder) {
	duration := NewDuration()

	result := Result{
		Gate:       renderGateName,
		Target:     "body styles",
		Pass:       true,
		DurationMs: 0,
	}

	// Run render checks in one evaluation
	checks, err := page.Evaluate(`() => {
		const checks = [];
		
		// Check stylesheets exist
		const hasStyles = document.styleSheets.length > 0;
		checks.push({name: "hasStylesheets", pass: hasStyles});
		
		// Check body margin is 0 (reset landed)
		const bodyStyle = getComputedStyle(document.body);
		checks.push({name: "bodyMarginReset", pass: bodyStyle.marginTop === "0px"});
		
		// Check body has background
		checks.push({name: "bodyBackground", pass: bodyStyle.backgroundColor !== "rgba(0, 0, 0, 0)"});
		
		// Check font is not default serif
		checks.push({name: "fontFamily", pass: !/^(Times|serif|-webkit-standard)/i.test(bodyStyle.fontFamily)});
		
		// Check no horizontal overflow
		const html = document.documentElement;
		checks.push({name: "noHorizontalOverflow", pass: html.scrollWidth <= html.clientWidth + 1});
		
		// Check images loaded
		const images = document.querySelectorAll("img");
		let allImagesLoaded = true;
		images.forEach(img => {
			if (!img.complete || img.naturalWidth === 0) allImagesLoaded = false;
		});
		checks.push({name: "imagesLoaded", pass: images.length === 0 || allImagesLoaded});
		
		// Check fonts loaded (lenient - don't fail if no fonts loaded)
		if (document.fonts) {
			let hasLoaded = false;
			let hasError = false;
			document.fonts.forEach(font => {
				if (font.status === "loaded") hasLoaded = true;
				if (font.status === "error") hasError = true;
			});
			checks.push({name: "fontsLoaded", pass: !hasError});
		}
		
		return checks;
	}`)

	if err != nil {
		result.Pass = false
		result.Findings = append(result.Findings, Finding{
			Gate:     renderGateName,
			Severity: "error",
			Message:  "Failed to evaluate render checks",
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
					Gate:     renderGateName,
					Severity: "warning",
					Message:  "Render check failed: " + name,
					Evidence: map[string]any{"check": name},
				})
			}
		}
	}

	result.DurationMs = duration.Elapsed()
	recorder.Add(result)
}
