package uxgates

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Mode controls gate behavior
type Mode int

const (
	ModeOff Mode = iota
	ModeWarn
	ModeFail
)

// String parses a mode from string
func (m Mode) String() string {
	switch m {
	case ModeOff:
		return "off"
	case ModeWarn:
		return "warn"
	case ModeFail:
		return "fail"
	default:
		return "unknown"
	}
}

// ParseMode parses a mode from string
func ParseMode(s string) Mode {
	switch s {
	case "off":
		return ModeOff
	case "warn":
		return ModeWarn
	case "fail":
		return ModeFail
	default:
		return ModeWarn
	}
}

// Finding represents a single gate finding
type Finding struct {
	Gate     string         `json:"gate"`
	Severity string         `json:"severity"`
	Message  string         `json:"message"`
	Evidence map[string]any `json:"evidence,omitempty"`
}

// Result represents a gate check result
type Result struct {
	Gate       string    `json:"gate"`
	Target     string    `json:"target"`
	Pass       bool      `json:"pass"`
	DurationMs int64     `json:"duration_ms"`
	Findings   []Finding `json:"findings,omitempty"`
}

// Recorder tracks gate results
type Recorder struct {
	results []Result
	mode    Mode
}

// NewRecorder creates a new recorder with the given mode
func NewRecorder(mode Mode) *Recorder {
	return &Recorder{
		results: make([]Result, 0),
		mode:    mode,
	}
}

// Add records a result. Returns error only in ModeFail
func (r *Recorder) Add(result Result) error {
	r.results = append(r.results, result)
	if !result.Pass && r.mode == ModeFail {
		return fmt.Errorf("gate %s failed: %d findings", result.Gate, len(result.Findings))
	}
	return nil
}

// Failed returns true if any gate failed
func (r *Recorder) Failed() bool {
	for _, res := range r.results {
		if !res.Pass {
			return true
		}
	}
	return false
}

// Findings returns all findings
func (r *Recorder) Findings() []Finding {
	var all []Finding
	for _, res := range r.results {
		all = append(all, res.Findings...)
	}
	return all
}

// Results returns all results
func (r *Recorder) Results() []Result {
	return r.results
}

// WriteJSON writes results to a file
func (r *Recorder) WriteJSON(dir string) (string, error) {
	path := dir + "/gates.json"
	data, err := json.MarshalIndent(r.results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal gates: %w", err)
	}
	if err := writeFile(path, data); err != nil {
		return "", fmt.Errorf("failed to write gates: %w", err)
	}
	return path, nil
}

// writeFile writes data to a file
func writeFile(path string, data []byte) error {
	dir := fmt.Sprintf("%s", path[:len(path)-len(path[len(path)-1:])])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Now is used for timing
var Now = time.Now

// Duration measures time between two points
type Duration struct {
	start time.Time
}

// NewDuration creates a new timer
func NewDuration() Duration {
	return Duration{start: time.Now()}
}

// Elapsed returns milliseconds elapsed
func (d Duration) Elapsed() int64 {
	return time.Since(d.start).Milliseconds()
}
