package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/llm-net/adbclaw/pkg/audio"
	"github.com/spf13/cobra"
)

var (
	audioDuration int
	audioRate     int
	audioFile     string
	audioStream   bool
)

var audioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Audio capture commands",
}

var audioCaptureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture system audio from device (Android 11+)",
	Long: `Capture system audio via REMOTE_SUBMIX. Requires Android 11+ (API 30).

By default, streams WAV to stdout for piping to other tools.
Use --file to save to a local file instead.

WARNING: Device speakers are muted while capturing.

Examples:
  adbclaw audio capture                          # Stream WAV to stdout (10s)
  adbclaw audio capture --duration 30000         # Stream 30s
  adbclaw audio capture --duration 0             # Stream until killed
  adbclaw audio capture --file recording.wav     # Save to file
  adbclaw audio capture --stream | asr-claw transcribe --stream`,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		writer.Verbose("audio capture: rate=%dHz duration=%dms file=%q", audioRate, audioDuration, audioFile)

		// Pre-check Android version for a friendly error message
		if err := checkAndroidVersion(); err != nil {
			writer.Fail("audio capture", "UNSUPPORTED_DEVICE", err.Error(),
				"Audio capture requires Android 11+ (API 30) with REMOTE_SUBMIX support", start)
			return nil
		}

		// Push DEX to device
		if err := audio.EnsureDEX(client); err != nil {
			writer.Fail("audio capture", "DEX_PUSH_FAILED", err.Error(),
				"Check device connection and try again", start)
			return nil
		}

		// Set up context: signal-aware + optional timeout
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if audioDuration > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(audioDuration+5000)*time.Millisecond)
			defer cancel()
		}

		proc, err := audio.Start(ctx, client, audioRate, audioDuration)
		if err != nil {
			writer.Fail("audio capture", "CAPTURE_START_FAILED", err.Error(),
				"Ensure device is connected and running Android 11+", start)
			return nil
		}
		defer proc.Stop()

		if audioFile != "" {
			return runFileCaptureMode(proc, start)
		}
		return runStreamCaptureMode(proc)
	},
}

func runStreamCaptureMode(proc *audio.Process) error {
	io.Copy(os.Stdout, proc.Stdout())
	proc.Wait()
	return nil
}

func runFileCaptureMode(proc *audio.Process, start time.Time) error {
	f, err := os.Create(audioFile)
	if err != nil {
		writer.Fail("audio capture", "FILE_CREATE_FAILED", err.Error(),
			"Check file path and permissions", start)
		return nil
	}
	defer f.Close()

	n, _ := io.Copy(f, proc.Stdout())
	proc.Wait()

	writer.Success("audio capture", map[string]interface{}{
		"file":        audioFile,
		"bytes":       n,
		"rate":        audioRate,
		"duration_ms": audioDuration,
	}, start)
	return nil
}

func checkAndroidVersion() error {
	result, err := client.Shell("getprop", "ro.build.version.sdk")
	if err != nil {
		return nil // can't check, let the device-side DEX handle it
	}
	sdkStr := strings.TrimSpace(result.Stdout)
	sdk, err := strconv.Atoi(sdkStr)
	if err != nil {
		return nil // can't parse, let device-side handle it
	}
	if sdk < 30 {
		return fmt.Errorf("device is Android SDK %d, but audio capture requires SDK 30+ (Android 11+)", sdk)
	}
	return nil
}

func init() {
	audioCaptureCmd.Flags().IntVar(&audioDuration, "duration", 10000, "Capture duration in milliseconds (0 = unlimited)")
	audioCaptureCmd.Flags().IntVar(&audioRate, "rate", 16000, "Sample rate in Hz")
	audioCaptureCmd.Flags().StringVar(&audioFile, "file", "", "Save to file instead of streaming to stdout")
	audioCaptureCmd.Flags().BoolVar(&audioStream, "stream", false, "Streaming mode (default behavior, for explicitness in pipe chains)")

	audioCmd.AddCommand(audioCaptureCmd)
	rootCmd.AddCommand(audioCmd)
}
