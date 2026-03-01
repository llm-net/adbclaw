package cmd

import (
	"encoding/base64"
	"os"
	"time"

	"github.com/llm-net/adbclaw/pkg/observe"
	"github.com/spf13/cobra"
)

var (
	observeMaxWidth int
)

var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Capture screenshot + UI tree + device state",
	Long:  "Captures a screenshot and UI element tree in parallel, returning both in a single response.",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		writer.Verbose("starting observe (screenshot + ui tree in parallel)")

		result := observe.Observe(client, observeMaxWidth)

		// Check if we got at least something
		if result.Screenshot == nil && result.UI == nil {
			writer.Fail("observe", "OBSERVE_FAILED",
				"Both screenshot and UI tree failed",
				"Check device connection: adbclaw device list", start)
			return nil
		}

		writer.Success("observe", result, start)
		return nil
	},
}

var (
	screenshotOutput   string
	screenshotMaxWidth int
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Capture a screenshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		writer.Verbose("capturing screenshot")

		data, err := observe.TakeScreenshot(client, screenshotMaxWidth)
		if err != nil {
			writer.Fail("screenshot", "SCREENSHOT_FAILED", err.Error(),
				"Ensure the device screen is on and unlocked", start)
			return nil
		}

		// If output file specified, write raw PNG
		if screenshotOutput != "" {
			if err := os.WriteFile(screenshotOutput, data, 0644); err != nil {
				writer.Fail("screenshot", "FILE_WRITE_ERROR", err.Error(), "", start)
				return nil
			}
			writer.Success("screenshot", map[string]interface{}{
				"format":     "png",
				"path":       screenshotOutput,
				"size_bytes": len(data),
			}, start)
			return nil
		}

		// Otherwise return base64 in JSON envelope
		writer.Success("screenshot", map[string]interface{}{
			"format":     "png",
			"base64":     base64.StdEncoding.EncodeToString(data),
			"size_bytes": len(data),
		}, start)
		return nil
	},
}

func init() {
	screenshotCmd.Flags().StringVarP(&screenshotOutput, "file", "f", "", "Save screenshot to file instead of base64 output")
	screenshotCmd.Flags().IntVar(&screenshotMaxWidth, "width", 0, "Max image width in pixels (0 = original size)")

	observeCmd.Flags().IntVar(&observeMaxWidth, "width", 0, "Max screenshot width in pixels (0 = original size)")

	rootCmd.AddCommand(observeCmd)
	rootCmd.AddCommand(screenshotCmd)
}
