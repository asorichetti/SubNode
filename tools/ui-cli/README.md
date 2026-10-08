# SubNode UI-CLI

A browser-automation CLI that proves UI changes actually work. Generates evidence bundles with screenshots, DOM snapshots, console logs, and network failures for every command run.

## Installation

```bash
# Install Playwright browsers
npx playwright install chromium

# Build
go build -o ui-cli .

# Or run directly
go run .
```

## Commands

### `navigate`
Open a URL and capture evidence bundle.

```bash
./ui-cli navigate --url /login --viewport mobile
```

Outputs:
- `screenshots/navigate-<timestamp>/` with:
  - `screenshot.png` - Full page screenshot
  - `dom.html` - Rendered DOM
  - `console.json` - All console output
  - `network.json` - Failed requests only
  - `gates.json` - UX gate results

### `journey [name]`
Run test journeys.

```bash
# List available journeys
./ui-cli journey --list

# Run specific journeys
./ui-cli journey home-auth login-form

# Run all journeys
./ui-cli journey --all

# Filter by suite
./ui-cli journey --suite auth

# Filter by name
./ui-cli journey --grep "login"
```

### `smoke`
Run core smoke tests in headless mode with strict gate validation.

```bash
./ui-cli smoke
```

### `screenshot`
Take screenshots across viewports and color schemes.

```bash
./ui-cli screenshot / --viewports desktop,mobile --schemes light,dark
```

### `interact`
Perform single browser interactions for poking.

```bash
./ui-cli interact click --selector "[data-testid=auth-submit]"
./ui-cli interact fill --selector "[data-testid=auth-username]" "testuser"
./ui-cli interact press --selector "[data-testid=auth-password]" "Enter"
./ui-cli interact screenshot --selector "[data-testid=node-page-title]"
```

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--env` | `local` | Environment (local, preview) |
| `--base-url` | - | Override base URL |
| `--headless` | computed | Run headless (auto-detected) |
| `--screenshot-dir` | `./screenshots` | Bundle output directory |
| `--ux-gates` | `warn` | Gate mode: off, warn, fail |
| `--viewport` | `desktop` | Viewport: desktop, tablet, mobile |
| `--color-scheme` | - | Color scheme: light, dark |

## UX Gates

Gates check that the page is working correctly, not just that assertions pass.

| Gate | Checks |
|------|--------|
| `pagestate` | Content marker present within timeout |
| `render` | CSS applied, fonts loaded, images loaded, no overflow |
| `console` | No console errors (except allowlisted) |
| `network` | No failed requests (except expected 404s) |
| `a11y` | Single h1, heading levels, alt text, accessible names |

## Adding Journeys

```go
package journeys

import "subnode-ui-cli/internal/journeys"

func init() {
    journeys.Register(journeys.Journey{
        Name: "home-auth",
        Suite: "auth",
        Desc: "User can navigate from home to login",
        Run: func(ctx *journeys.Context) error {
            // Navigate
            if err := ctx.Goto("/login"); err != nil {
                return err
            }
            
            // Assert
            ctx.Page.Locator("[data-testid=auth-login-form]").WaitFor()
            
            return nil
        },
    })
}
```

## Running Against SubNode

```bash
# Start SubNode locally
cd /Users/alexsorichetti/projects/SubNode
make dev-backend
make dev-frontend

# Run UI tests
cd tools/ui-cli
./ui-cli navigate --url /
./ui-cli journey --list
./ui-cli smoke
```

## Definition of Done

- [x] `go build ./...` passes
- [x] `go vet ./...` passes  
- [x] `gofmt -l .` is clean
- [x] `journey --list` prints all journeys
- [x] All commands work and write bundles
- [ ] All journeys pass with `--ux-gates fail`
- [ ] Gates proven to fire on broken page (§8)
- [ ] README documents setup and commands
