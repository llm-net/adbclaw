package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/llm-net/adb-claw/pkg/monitor"
	"github.com/spf13/cobra"
)

var (
	monitorDuration int
	monitorInterval int
	monitorStream   bool
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Continuously monitor UI text via accessibility framework",
	Long: `Monitor UI text on the device by connecting directly to the Android accessibility
framework. Unlike 'ui tree' which uses uiautomator dump, this command skips video
surface nodes and works reliably during live streams and video playback.

Two modes:
  Bounded (default): runs for --duration ms, returns all captured text in JSON envelope
  Streaming (--stream): outputs each new text as a JSON line in real time

Examples:
  adb-claw monitor                          # 10s bounded, returns JSON envelope
  adb-claw monitor --duration 30000         # 30s bounded
  adb-claw monitor --stream --duration 60000  # 60s streaming, JSON lines
  adb-claw monitor --stream                 # 10s streaming`,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		writer.Verbose("monitor: duration=%dms interval=%dms stream=%v", monitorDuration, monitorInterval, monitorStream)

		// Push DEX to device
		if err := monitor.EnsureDEX(client); err != nil {
			writer.Fail("monitor", "DEX_PUSH_FAILED", err.Error(),
				"Check device connection and try again", start)
			return nil
		}

		if monitorStream {
			return runStreamingMode(start)
		}
		return runBoundedMode(start)
	},
}

func runBoundedMode(start time.Time) error {
	count := monitorDuration / monitorInterval
	if count < 1 {
		count = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(monitorDuration+5000)*time.Millisecond)
	defer cancel()

	proc, err := monitor.Start(ctx, client, monitorInterval, count)
	if err != nil {
		writer.Fail("monitor", "MONITOR_START_FAILED", err.Error(),
			"Ensure device is connected and adb is working", start)
		return nil
	}
	defer proc.Stop()

	var texts []monitor.TextEntry
	for line := range proc.Lines() {
		entry, err := monitor.ParseLine(line)
		if err != nil {
			writer.Verbose("monitor: skip unparseable line: %s", line)
			continue
		}
		texts = append(texts, *entry)
	}

	proc.Wait()

	if texts == nil {
		texts = []monitor.TextEntry{}
	}

	writer.Success("monitor", map[string]interface{}{
		"texts":       texts,
		"count":       len(texts),
		"duration_ms": time.Since(start).Milliseconds(),
	}, start)
	return nil
}

func runStreamingMode(start time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(monitorDuration+5000)*time.Millisecond)
	defer cancel()

	count := monitorDuration / monitorInterval
	if count < 1 {
		count = 1
	}

	proc, err := monitor.Start(ctx, client, monitorInterval, count)
	if err != nil {
		writer.Fail("monitor", "MONITOR_START_FAILED", err.Error(),
			"Ensure device is connected and adb is working", start)
		return nil
	}
	defer proc.Stop()

	enc := json.NewEncoder(os.Stdout)
	for line := range proc.Lines() {
		entry, err := monitor.ParseLine(line)
		if err != nil {
			writer.Verbose("monitor: skip unparseable line: %s", line)
			continue
		}
		enc.Encode(entry)
	}

	proc.Wait()
	return nil
}

func init() {
	monitorCmd.Flags().IntVar(&monitorDuration, "duration", 10000, "Total monitoring duration in milliseconds")
	monitorCmd.Flags().IntVar(&monitorInterval, "interval", 2000, "Poll interval in milliseconds")
	monitorCmd.Flags().BoolVar(&monitorStream, "stream", false, "Streaming mode: output JSON lines instead of envelope")

	rootCmd.AddCommand(monitorCmd)
}
