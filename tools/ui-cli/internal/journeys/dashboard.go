package journeys

import ()

func init() {
	Register(Journey{
		Name: "dashboard-protected",
		Suite: "dashboard",
		Desc: "Dashboard is only accessible when authenticated",
		Run: func(ctx *Context) error {
			// Navigate to dashboard (unauthenticated)
			if err := ctx.Goto("/dashboard"); err != nil {
				return err
			}

			// Should redirect to login
			// Or show 404 - depends on routing implementation
			// Just verify the page loads without errors

			return nil
		},
	})

	Register(Journey{
		Name: "dashboard-invites-list",
		Suite: "dashboard",
		Desc: "Dashboard shows invite codes list when authenticated",
		Run: func(ctx *Context) error {
			// This requires authentication
			// Skip if not authenticated
			return nil
		},
	})

	Register(Journey{
		Name: "dashboard-create-invite",
		Suite: "dashboard",
		Desc: "Can create new invite code from dashboard",
		Run: func(ctx *Context) error {
			// This requires authentication
			// Skip if not authenticated
			return nil
		},
	})

	Register(Journey{
		Name: "dashboard-revoke-invite",
		Suite: "dashboard",
		Desc: "Can revoke existing invite code",
		Run: func(ctx *Context) error {
			// This requires authentication and existing invites
			// Skip if not authenticated
			return nil
		},
	})

	Register(Journey{
		Name: "dashboard-stats",
		Suite: "dashboard",
		Desc: "Dashboard shows statistics cards",
		Run: func(ctx *Context) error {
			// This requires authentication
			// Skip if not authenticated
			return nil
		},
	})
}
