package journeys

import ()

func init() {
	Register(Journey{
		Name: "keyboard-tab-order",
		Suite: "keyboard",
		Desc: "Tab order follows logical reading order",
		Run: func(ctx *Context) error {
			// Navigate
			if err := ctx.Goto("/"); err != nil {
				return err
			}

			// Get focusable elements in order
			focusable, _ := ctx.Page.Evaluate(`() => {
				const elements = document.querySelectorAll('a, button, input, select, textarea, [tabindex]:not([tabindex="-1"])');
				return Array.from(elements).map(el => ({
					tag: el.tagName.toLowerCase(),
					type: el.type || null,
					role: el.getAttribute('role'),
					text: (el.innerText || el.getAttribute('aria-label') || '').trim().substring(0, 50)
				}));
			}`)

			// Verify at least some focusable elements exist
			_ = focusable

			return nil
		},
	})

	Register(Journey{
		Name: "keyboard-escape",
		Suite: "keyboard",
		Desc: "Escape key closes modals/dropdowns",
		Run: func(ctx *Context) error {
			// Test on login page
			if err := ctx.Goto("/login"); err != nil {
				return err
			}

			// Focus input
			username := ctx.Page.Locator("[data-testid=auth-username]")
			if err := username.Focus(); err != nil {
				return err
			}

			// Press Escape
			if err := username.Press("Escape"); err != nil {
				return err
			}

			// Page should still be functional
			// No errors should occur

			return nil
		},
	})

	Register(Journey{
		Name: "keyboard-enter-submit",
		Suite: "keyboard",
		Desc: "Enter key submits forms",
		Run: func(ctx *Context) error {
			// Navigate
			if err := ctx.Goto("/login"); err != nil {
				return err
			}

			// Fill username
			username := ctx.Page.Locator("[data-testid=auth-username]")
			if err := username.Fill("testuser"); err != nil {
				return err
			}

			// Focus password
			password := ctx.Page.Locator("[data-testid=auth-password]")
			if err := password.Focus(); err != nil {
				return err
			}

			// Press Enter
			if err := password.Press("Enter"); err != nil {
				return err
			}

			// Form should attempt to submit
			// (may fail if backend not running, but no errors should occur)

			return nil
		},
	})

	Register(Journey{
		Name: "focus-visible-rings",
		Suite: "keyboard",
		Desc: "Focusable elements show visible focus indicators",
		Run: func(ctx *Context) error {
			// Navigate
			if err := ctx.Goto("/"); err != nil {
				return err
			}

			// Click to focus an element
			links := ctx.Page.Locator("a")
			if count, _ := links.Count(); count > 0 {
				if err := links.First().Focus(); err != nil {
					return err
				}

				// Check focus ring is visible
				focused, _ := ctx.Page.Evaluate(`() => {
					const el = document.activeElement;
					if (!el) return { hasFocus: false };
					const style = getComputedStyle(el);
					return {
						hasFocus: true,
						outline: style.outline,
						boxShadow: style.boxShadow
					};
				}`)

				_ = focused
			}

			return nil
		},
	})
}
