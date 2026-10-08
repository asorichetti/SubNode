package cmd

import (
	"github.com/spf13/cobra"
)

var smokeCmd = &cobra.Command{
	Use:   "smoke",
	Short: "Run core smoke tests",
	Long:  `Run core journeys with --ux-gates fail in headless mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Override flags for smoke tests
		uxGates = "fail"
		headless = true

		return runJourneys([]string{})
	},
}
