package browser

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mxschmitt/playwright-go"
)

// ConsoleEntry represents a console log entry
type ConsoleEntry struct {
	Type    string `json:"type"` // "log", "error", "warning", etc.
	Text    string `json:"text"`
	IsError bool   `json:"is_error"`
}

// NetworkEntry represents a failed network request
type NetworkEntry struct {
	URL         string `json:"url"`
	Status      int    `json:"status"`
	Method      string `json:"method"`
	FailureText string `json:"failure_text,omitempty"`
}

// Snapshots holds captured browser events
type Snapshots struct {
	console []ConsoleEntry
	network []NetworkEntry
	errors  []string
	mu      sync.RWMutex
}

// NewSnapshots creates a new Snapshots instance
func NewSnapshots() *Snapshots {
	return &Snapshots{
		console: make([]ConsoleEntry, 0),
		network: make([]NetworkEntry, 0),
		errors:  make([]string, 0),
	}
}

// Attach attaches event listeners to the page
func Attach(page playwright.Page, snapshots *Snapshots) {
	// Console events
	page.On("console", func(consoleMsg playwright.ConsoleMessage) {
		if consoleMsg == nil {
			return
		}
		msgType := consoleMsg.Type()
		text := consoleMsg.Text()
		isError := msgType == "error" || msgType == "warning"

		snapshots.mu.Lock()
		defer snapshots.mu.Unlock()
		snapshots.console = append(snapshots.console, ConsoleEntry{
			Type:    msgType,
			Text:    text,
			IsError: isError,
		})
	})

	// Page errors
	page.On("pageerror", func(err error) {
		if err == nil {
			return
		}
		snapshots.mu.Lock()
		defer snapshots.mu.Unlock()
		snapshots.errors = append(snapshots.errors, err.Error())
	})

	// Failed responses (status >= 400)
	page.On("response", func(resp playwright.Response) {
		if resp == nil {
			return
		}
		status := resp.Status()
		if status >= 400 {
			snapshots.mu.Lock()
			defer snapshots.mu.Unlock()
			url := resp.URL()
			method := "GET"
			if req := resp.Request(); req != nil {
				method = req.Method()
			}
			snapshots.network = append(snapshots.network, NetworkEntry{
				URL:    url,
				Status: status,
				Method: method,
			})
		}
	})

	// Request failures (no response received)
	page.On("requestfailed", func(req playwright.Request) {
		if req == nil {
			return
		}
		snapshots.mu.Lock()
		defer snapshots.mu.Unlock()
		var failureText string
		if fail := req.Failure(); fail != nil {
			failureText = fail.Error()
		}
		snapshots.network = append(snapshots.network, NetworkEntry{
			URL:         req.URL(),
			Status:      0,
			Method:      req.Method(),
			FailureText: failureText,
		})
	})
}

// Snapshot returns copies of the captured data
func (s *Snapshots) Snapshot() (ConsoleEntrySlice, NetworkEntrySlice, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	consoleCopy := make(ConsoleEntrySlice, len(s.console))
	copy(consoleCopy, s.console)

	networkCopy := make(NetworkEntrySlice, len(s.network))
	copy(networkCopy, s.network)

	errorsCopy := make([]string, len(s.errors))
	copy(errorsCopy, s.errors)

	return consoleCopy, networkCopy, errorsCopy
}

// ConsoleEntrySlice is a type alias for JSON marshaling
type ConsoleEntrySlice []ConsoleEntry

// MarshalJSON implements json.Marshaler
func (c ConsoleEntrySlice) MarshalJSON() ([]byte, error) {
	return json.Marshal([]ConsoleEntry(c))
}

// NetworkEntrySlice is a type alias for JSON marshaling
type NetworkEntrySlice []NetworkEntry

// MarshalJSON implements json.Marshaler
func (n NetworkEntrySlice) MarshalJSON() ([]byte, error) {
	return json.Marshal([]NetworkEntry(n))
}

// WriteConsole writes console entries to a file
func WriteConsole(path string, entries ConsoleEntrySlice) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal console: %w", err)
	}
	if err := writeFile(path, data); err != nil {
		return fmt.Errorf("failed to write console: %w", err)
	}
	return nil
}

// WriteNetwork writes network entries to a file
func WriteNetwork(path string, entries NetworkEntrySlice) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal network: %w", err)
	}
	if err := writeFile(path, data); err != nil {
		return fmt.Errorf("failed to write network: %w", err)
	}
	return nil
}

// WriteErrors writes error entries to a file
func WriteErrors(path string, errors []string) error {
	data, err := json.MarshalIndent(errors, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal errors: %w", err)
	}
	if err := writeFile(path, data); err != nil {
		return fmt.Errorf("failed to write errors: %w", err)
	}
	return nil
}
