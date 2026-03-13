package audio

import (
	"context"
	"crypto/md5"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/llm-net/adbclaw/pkg/adb"
)

//go:embed classes.dex
var dexData []byte

const deviceDEXPath = "/data/local/tmp/adbclaw-audio.dex"

// Process wraps an adb exec-out subprocess streaming audio from the device.
type Process struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
}

// EnsureDEX pushes the embedded audio DEX file to the device if not already present
// or if the on-device copy differs.
func EnsureDEX(client *adb.Client) error {
	localMD5 := fmt.Sprintf("%x", md5.Sum(dexData))

	result, err := client.Shell("md5sum", deviceDEXPath)
	if err == nil && result.ExitCode == 0 {
		parts := strings.Fields(strings.TrimSpace(result.Stdout))
		if len(parts) > 0 && parts[0] == localMD5 {
			return nil // already up to date
		}
	}

	tmpFile := filepath.Join(os.TempDir(), "adbclaw-audio.dex")
	if err := os.WriteFile(tmpFile, dexData, 0644); err != nil {
		return fmt.Errorf("write temp DEX: %w", err)
	}
	defer os.Remove(tmpFile)

	result, err = client.RawCommand("push", tmpFile, deviceDEXPath)
	if err != nil {
		return fmt.Errorf("push DEX to device: %w", err)
	}
	if result.ExitCode != 0 {
		msg := strings.TrimSpace(result.Stderr)
		if msg == "" {
			msg = fmt.Sprintf("adb push returned exit code %d", result.ExitCode)
		}
		return fmt.Errorf("push DEX to device: %s", msg)
	}
	return nil
}

// Start launches the audio capture DEX on the device via app_process.
// Uses exec-out (not shell) for binary-safe stdout streaming.
// sampleRate is in Hz (e.g. 16000). durationMs is capture duration (0 = unlimited).
func Start(ctx context.Context, client *adb.Client, sampleRate, durationMs int) (*Process, error) {
	args := client.BaseArgs()
	args = append(args, "exec-out",
		fmt.Sprintf("CLASSPATH=%s", deviceDEXPath),
		"app_process", "/", "ADBClawAudio",
		"--rate", fmt.Sprintf("%d", sampleRate),
		"--duration", fmt.Sprintf("%d", durationMs),
	)

	cmd := exec.CommandContext(ctx, client.ADBPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr // forward device-side logs to host stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start audio capture: %w", err)
	}

	return &Process{cmd: cmd, stdout: stdout}, nil
}

// Stdout returns the raw audio stream (WAV header + PCM data).
func (p *Process) Stdout() io.Reader {
	return p.stdout
}

// Wait waits for the audio capture process to exit.
func (p *Process) Wait() error {
	return p.cmd.Wait()
}

// Stop kills the audio capture process.
func (p *Process) Stop() {
	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
}
