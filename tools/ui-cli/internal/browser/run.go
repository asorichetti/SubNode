package browser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mxschmitt/playwright-go"
)

// Artifacts contains paths to all generated files for a run
type Artifacts struct {
	RunDir, ScreenshotPath, HTMLPath, ConsolePath, NetworkPath, UXGatesPath string
}

// RunBundle manages a single test run
type RunBundle struct {
	runDir         string
	screenshotPath string
	htmlPath       string
	consolePath    string
	networkPath    string
	uxgatesPath    string
	snapshots      *Snapshots
	page           playwright.Page
	finished       bool
	cmdName        string
}

// ArtifactsHook is the interface for writing gate results
type ArtifactsHook interface {
	WriteJSON(dir string) (string, error)
}

// NewRun creates a new run directory and bundle
func NewRun(cmdName string, headless bool, dir string, vp playwright.Size, scheme string) (*RunBundle, error) {
	// Create timestamped run directory
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	runDir := filepath.Join(dir, cmdName+"-"+timestamp)
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create run directory: %w", err)
	}

	bundle := &RunBundle{
		runDir:    runDir,
		snapshots: NewSnapshots(),
		cmdName:   cmdName,
	}

	// Set default paths
	bundle.screenshotPath = filepath.Join(runDir, "screenshot.png")
	bundle.htmlPath = filepath.Join(runDir, "dom.html")
	bundle.consolePath = filepath.Join(runDir, "console.json")
	bundle.networkPath = filepath.Join(runDir, "network.json")
	bundle.uxgatesPath = filepath.Join(runDir, "gates.json")

	return bundle, nil
}

// SetPage associates a page with this bundle for DOM capture
func (r *RunBundle) SetPage(page playwright.Page) {
	r.page = page
}

// Finish writes all artifacts and marks the run as complete
func (r *RunBundle) Finish() (Artifacts, error) {
	if r.finished {
		return Artifacts{}, fmt.Errorf("run already finished")
	}

	var artifacts Artifacts
	var firstErr error

	// Write DOM snapshot
	if r.page != nil {
		html, err := r.page.Content()
		if err != nil {
			firstErr = fmt.Errorf("failed to get DOM: %w", err)
		} else {
			if err := writeFile(r.htmlPath, []byte(html)); err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("failed to write DOM: %w", err)
				}
			}
		}
	}

	// Write console
	console, _, _ := r.snapshots.Snapshot()
	if err := WriteConsole(r.consolePath, console); err != nil {
		if firstErr == nil {
			firstErr = fmt.Errorf("failed to write console: %w", err)
		}
	}

	// Write network
	_, network, _ := r.snapshots.Snapshot()
	if err := WriteNetwork(r.networkPath, network); err != nil {
		if firstErr == nil {
			firstErr = fmt.Errorf("failed to write network: %w", err)
		}
	}

	// Write gate results if hook is set
	// Note: Gates are written separately via the hook

	r.finished = true
	artifacts.RunDir = r.runDir
	artifacts.ScreenshotPath = r.screenshotPath
	artifacts.HTMLPath = r.htmlPath
	artifacts.ConsolePath = r.consolePath
	artifacts.NetworkPath = r.networkPath
	artifacts.UXGatesPath = r.uxgatesPath

	return artifacts, firstErr
}

// FinishOrLog is like Finish but logs errors instead of returning them
func (r *RunBundle) FinishOrLog() Artifacts {
	artifacts, err := r.Finish()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: run finish errors: %v\n", err)
	}
	return artifacts
}

// Close writes the bundle as a last resort if nothing else did
func (r *RunBundle) Close() {
	if r.finished {
		return
	}

	// Write what we can
	r.Finish()

	// Log that this was a fallback
	fmt.Fprintf(os.Stderr, "Run bundle written via Close() fallback\n")
}

// WriteGateResults writes UX gate results to the run directory
func (r *RunBundle) WriteGateResults(hook ArtifactsHook) (string, error) {
	path, err := hook.WriteJSON(r.runDir)
	if err != nil {
		return "", fmt.Errorf("failed to write gate results: %w", err)
	}
	r.uxgatesPath = path
	return path, nil
}

// WriteScreenshot captures and saves a screenshot
func (r *RunBundle) WriteScreenshot(page playwright.Page) error {
	_, err := Screenshot(page, r.screenshotPath)
	return err
}

// writeFile writes data to a file
func writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// ReadArtifacts reads the artifacts metadata file
func ReadArtifacts(path string) (Artifacts, error) {
	var artifacts Artifacts
	data, err := os.ReadFile(path)
	if err != nil {
		return artifacts, fmt.Errorf("failed to read artifacts: %w", err)
	}
	if err := json.Unmarshal(data, &artifacts); err != nil {
		return artifacts, fmt.Errorf("failed to parse artifacts: %w", err)
	}
	return artifacts, nil
}

// ListRuns lists all run directories in a folder
func ListRuns(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to list runs: %w", err)
	}

	var runs []string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 20 {
			runs = append(runs, entry.Name())
		}
	}
	return runs, nil
}
