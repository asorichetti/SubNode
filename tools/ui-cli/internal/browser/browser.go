package browser

import (
	"fmt"
	"log"
	"sync"

	"github.com/mxschmitt/playwright-go"
)

// ViewportPresets defines standard viewport sizes
var ViewportPresets = map[string]playwright.Size{
	"desktop": {Width: 1280, Height: 900},
	"tablet":  {Width: 820, Height: 1180},
	"mobile":  {Width: 390, Height: 844},
}

// Browser wraps playwright.Browser with additional lifecycle methods
type Browser struct {
	pw       *playwright.Playwright
	browser  playwright.Browser
	headless bool
}

// New creates a new browser instance
func New(headless bool) (*Browser, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to start playwright: %w", err)
	}

	opts := []playwright.BrowserTypeLaunchOptions{
		{
			Headless: &headless,
			Args:     []string{"--disable-dev-shm-usage"},
		},
	}

	browser, err := pw.Chromium.Launch(opts...)
	if err != nil {
		_ = pw.Stop()
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	return &Browser{
		pw:       pw,
		browser:  browser,
		headless: headless,
	}, nil
}

// NewPage creates a new page with the given viewport and color scheme
func (b *Browser) NewPage(viewport playwright.Size, colorScheme string) (playwright.Page, error) {
	ctxOpts := []playwright.BrowserNewContextOptions{
		{Viewport: &viewport},
	}

	if colorScheme != "" {
		cs := playwright.ColorScheme(colorScheme)
		ctxOpts = append(ctxOpts, playwright.BrowserNewContextOptions{ColorScheme: &cs})
	}

	context, err := b.browser.NewContext(ctxOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	return page, nil
}

// Screenshot captures a screenshot of the page
func Screenshot(page playwright.Page, outPath string) (string, error) {
	opts := []playwright.PageScreenshotOptions{
		{
			Path:     &outPath,
			FullPage: boolPtr(true),
		},
	}

	_, err := page.Screenshot(opts...)
	if err != nil {
		return "", fmt.Errorf("failed to screenshot: %w", err)
	}
	return outPath, nil
}

func boolPtr(b bool) *bool {
	return &b
}

// Close shuts down the browser
func (b *Browser) Close() {
	if err := SafeClose(b); err != nil && len(err.Errors) > 0 {
		log.Printf("browser close errors: %v", err)
	}
}

// CloseError is the return type for Close operations that might fail
type CloseError struct {
	Errors []error
}

func (ce *CloseError) Error() string {
	msg := "close errors:"
	for _, err := range ce.Errors {
		msg += fmt.Sprintf("\n  - %v", err)
	}
	return msg
}

// SafeClose wraps Close and collects errors safely for defer
func SafeClose(b *Browser) *CloseError {
	var errs CloseError
	if err := b.browser.Close(); err != nil {
		errs.Errors = append(errs.Errors, fmt.Errorf("browser: %w", err))
	}
	if err := b.pw.Stop(); err != nil {
		errs.Errors = append(errs.Errors, fmt.Errorf("playwright: %w", err))
	}
	if len(errs.Errors) > 0 {
		log.Printf("browser close warnings: %v", &errs)
	}
	return &errs
}

// Size represents viewport dimensions
type Size struct {
	Width  int
	Height int
}

// SizeFromPreset returns a Size from a named preset
func SizeFromPreset(name string) playwright.Size {
	preset, ok := ViewportPresets[name]
	if !ok {
		return ViewportPresets["desktop"]
	}
	return preset
}

// Mutex provides thread-safe access to shared resources
var Mutex = &sync.Mutex{}
