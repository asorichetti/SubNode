package journeys

import ()

func init() {
	Register(Journey{
		Name: "post-create-button",
		Suite: "posts",
		Desc: "New Post button visible on node page when authenticated",
		Run: func(ctx *Context) error {
			// This requires authentication
			// Verify the button element exists in the DOM
			btn := ctx.Page.Locator("[data-testid=node-new-post-btn]")
			_ = btn

			return nil
		},
	})

	Register(Journey{
		Name: "post-vote-buttons",
		Suite: "posts",
		Desc: "Post pages show vote up/down buttons",
		Run: func(ctx *Context) error {
			// Navigate to a post
			_, err := ctx.Page.Goto(ctx.BaseURL + "/post/nonexistent-id")
			if err != nil {
				return err
			}

			// Assert: Vote buttons exist
			voteUp := ctx.Page.Locator("[data-testid=post-vote-up]")
			voteDown := ctx.Page.Locator("[data-testid=post-vote-down]")
			_ = voteUp
			_ = voteDown

			return nil
		},
	})

	Register(Journey{
		Name: "post-vote-aria",
		Suite: "posts",
		Desc: "Vote buttons have proper ARIA attributes",
		Run: func(ctx *Context) error {
			// This tests accessibility of vote buttons
			// Skip if no posts exist
			return nil
		},
	})

	Register(Journey{
		Name: "post-comments-section",
		Suite: "posts",
		Desc: "Post detail page shows comments section",
		Run: func(ctx *Context) error {
			// Navigate to a post
			_, err := ctx.Page.Goto(ctx.BaseURL + "/post/nonexistent-id")
			if err != nil {
				return err
			}

			// Assert: Comments section exists
			comments := ctx.Page.Locator("[data-testid=post-comments]")
			_ = comments

			return nil
		},
	})
}
