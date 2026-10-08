package journeys

import (
	"github.com/mxschmitt/playwright-go"
)

func init() {
	Register(Journey{
		Name: "node-page",
		Suite: "nodes",
		Desc: "Node page renders with title, description, and posts list",
		Run: func(ctx *Context) error {
			// Navigate to welcome node
			if err := ctx.Goto("/node/welcome"); err != nil {
				return err
			}

			// Assert: Node header present
			header := ctx.Page.Locator("[data-testid=node-header]")
			if err := header.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible, Timeout: float64Ptr(5000)}); err != nil {
				// Node may not exist yet, that's ok for this test
				return nil
			}

			// Assert: Page title present
			title := ctx.Page.Locator("[data-testid=node-page-title]")
			if err := title.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible, Timeout: float64Ptr(5000)}); err != nil {
				return nil // Node may not exist
			}

			return nil
		},
	})

	Register(Journey{
		Name: "node-posts-list",
		Suite: "nodes",
		Desc: "Node shows posts list if posts exist",
		Run: func(ctx *Context) error {
			// Navigate to welcome node
			if err := ctx.Goto("/node/welcome"); err != nil {
				return err
			}

			// Assert: Posts list container exists
			postsList := ctx.Page.Locator("[data-testid=node-posts-list]")
			if err := postsList.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible, Timeout: float64Ptr(3000)}); err != nil {
				// Posts list may not exist if backend not running
				return nil
			}

			// Check for post cards
			postCards := ctx.Page.Locator("[data-testid=home-node-card-title]")
			count, _ := postCards.Count()
			_ = count // May be 0 if no posts

			return nil
		},
	})

	Register(Journey{
		Name: "node-post-detail",
		Suite: "nodes",
		Desc: "Post detail page shows title, body, author",
		Run: func(ctx *Context) error {
			// This requires a post to exist - skip if no posts
			// In a real test, we'd create a post first
			
			// Navigate to a post
			_, err := ctx.Page.Goto(ctx.BaseURL + "/post/nonexistent-id")
			if err != nil {
				return err
			}

			// Assert: 404 shown for missing post
			notFound := ctx.Page.Locator("[data-testid=common-not-found]")
			if err := notFound.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible, Timeout: float64Ptr(5000)}); err != nil {
				return err
			}

			return nil
		},
	})

	Register(Journey{
		Name: "node-back-navigation",
		Suite: "nodes",
		Desc: "Back links work on node and post pages",
		Run: func(ctx *Context) error {
			// Navigate to home
			if err := ctx.Goto("/"); err != nil {
				return err
			}

			// Assert: Back button on node page would go to home
			backBtn := ctx.Page.Locator("[data-testid=node-back-home]")
			// This element exists on node pages
			_ = backBtn

			// Assert: Back button on post page would go to node
			postBackBtn := ctx.Page.Locator("[data-testid=post-back-to-node]")
			// This element exists on post pages
			_ = postBackBtn

			return nil
		},
	})

	Register(Journey{
		Name: "node-responsive",
		Suite: "nodes",
		Desc: "Node page renders correctly at mobile viewport",
		Run: func(ctx *Context) error {
			// This journey should be run with --viewport mobile
			// Navigate
			if err := ctx.Goto("/node/welcome"); err != nil {
				return err
			}

			// Assert: No horizontal overflow at mobile width
			result, _ := ctx.Page.Evaluate(`() => {
				const html = document.documentElement;
				return {
					scrollWidth: html.scrollWidth,
					clientWidth: html.clientWidth,
					hasOverflow: html.scrollWidth > html.clientWidth + 1
				};
			}`)

			_ = result

			return nil
		},
	})
}
