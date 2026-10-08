package journeys

import (
	"subnode-ui-cli/internal/uxgates"
)

func init() {
	Register(Journey{
		Name: "gate-pagestate-fires",
		Suite: "gates",
		Desc: "Proves pagestate gate fires on broken page",
		Run: func(ctx *Context) error {
			// Navigate directly to broken page (bypassing gotoAndGates)
			_, err := ctx.Page.Goto(ctx.BaseURL + "/broken")
			if err != nil {
				return err
			}

			// Replace ctx.Gates with a fresh recorder to avoid journey runner gate check
			ctx.Gates = uxgates.NewRecorder(uxgates.ModeOff)

			// Run pagestate gate
			uxgates.RunPageState(ctx.Page, ctx.Gates)

			// Verify gate failed
			results := ctx.Gates.Results()
			if len(results) == 0 {
				return errNoResults
			}

			last := results[len(results)-1]
			if last.Pass {
				return errGateShouldFail
			}

			return nil
		},
	})

	Register(Journey{
		Name: "gate-all-fires",
		Suite: "gates",
		Desc: "Proves all gates fire when page is broken",
		Run: func(ctx *Context) error {
			// Navigate directly to broken page (bypassing gotoAndGates)
			_, err := ctx.Page.Goto(ctx.BaseURL + "/broken")
			if err != nil {
				return err
			}

			// Wait for page to load and JS to execute
			ctx.Page.WaitForTimeout(1000)

			// Replace ctx.Gates with a fresh recorder to avoid journey runner gate check
			ctx.Gates = uxgates.NewRecorder(uxgates.ModeOff)

			// Run all gates
			uxgates.RunPageState(ctx.Page, ctx.Gates)
			uxgates.RunRender(ctx.Page, ctx.Gates)
			uxgates.RunA11y(ctx.Page, ctx.Gates)

			// Verify at least one gate failed
			results := ctx.Gates.Results()
			if len(results) == 0 {
				return errNoResults
			}

			failedCount := 0
			for _, r := range results {
				if !r.Pass {
					failedCount++
				}
			}

			if failedCount == 0 {
				return errGateShouldFail
			}

			return nil
		},
	})
}

var errNoResults = &journeyError{msg: "no gate results returned"}
var errGateShouldFail = &journeyError{msg: "gate should have failed but passed"}

type journeyError struct {
	msg string
}

func (e *journeyError) Error() string {
	return e.msg
}
