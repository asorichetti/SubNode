package journeys

import (
	"fmt"
	"strings"

	"subnode-ui-cli/internal/browser"
	"subnode-ui-cli/internal/uxgates"

	"github.com/mxschmitt/playwright-go"
)

// Journey represents a single UI test journey
type Journey struct {
	Name  string
	Suite string
	Desc  string
	Run   func(*Context) error
}

// Context carries the page and utilities for a journey
type Context struct {
	Page    playwright.Page
	BaseURL string
	Gates   *uxgates.Recorder
	Logf    func(format string, args ...any)
	Goto    func(url string) error
}

var journeys = make(map[string]Journey)

// Register adds a journey to the registry
func Register(j Journey) {
	if _, exists := journeys[j.Name]; exists {
		panic("duplicate journey name: " + j.Name)
	}
	journeys[j.Name] = j
}

// All returns all journeys sorted by name
func All() []Journey {
	result := make([]Journey, 0, len(journeys))
	for _, j := range journeys {
		result = append(result, j)
	}
	// Sort by name (simple bubble sort)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Name > result[j].Name {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

// GetAllJourneys returns all journeys (exported for cmd package)
var GetAllJourneys = All

// Expect is a helper for assertions
func Expect(actual string, expected string, format string, args ...any) {
	if actual != expected {
		msg := format
		if len(args) > 0 {
			msg = fmt.Sprintf(format, args...)
		}
		_ = msg
		panic(fmt.Sprintf("expected %s, found %s", expected, actual))
	}
}

// BrowserInstance creates a new browser
var BrowserInstance = func() (*browser.Browser, error) {
	return browser.New(false)
}

// NewJourneyContext creates a new context for a journey
var NewJourneyContext = func(b *browser.Browser, j Journey) (*Context, error) {
	vp := browser.SizeFromPreset("desktop")
	page, err := b.NewPage(vp, "light")
	if err != nil {
		return nil, err
	}

	recorder := uxgates.NewRecorder(uxgates.ModeWarn)

	return &Context{
		Page:  page,
		Gates: recorder,
		Logf:  func(format string, args ...any) {},
	}, nil
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
