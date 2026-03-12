package monitor

import (
	"bufio"
	"context"
	"crypto/md5"
	_ "embed"
	"encoding/json"
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

const deviceDEXPath = "/data/local/tmp/adbclaw-monitor.dex"

// TextEntry represents a single text node captured from the UI.
type TextEntry struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

// Process wraps an adb shell subprocess running the monitor DEX.
type Process struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	lines  chan string
	done   chan struct{}
}

// EnsureDEX pushes the embedded DEX file to the device if not already present
// or if the on-device copy differs.
func EnsureDEX(client *adb.Client) error {
	// Check if DEX already exists on device with matching md5
	localMD5 := fmt.Sprintf("%x", md5.Sum(dexData))

	result, err := client.Shell("md5sum", deviceDEXPath)
	if err == nil && result.ExitCode == 0 {
		parts := strings.Fields(strings.TrimSpace(result.Stdout))
		if len(parts) > 0 && parts[0] == localMD5 {
			return nil // already up to date
		}
	}

	// Write DEX to temp file on host
	tmpFile := filepath.Join(os.TempDir(), "adbclaw-monitor.dex")
	if err := os.WriteFile(tmpFile, dexData, 0644); err != nil {
		return fmt.Errorf("write temp DEX: %w", err)
	}
	defer os.Remove(tmpFile)

	// Push to device
	result, err = client.RawCommand("push", tmpFile, deviceDEXPath)
	if err != nil {
		return fmt.Errorf("push DEX to device: %w", err)
	}
	if result.ExitCode != 0 {
		msg := strings.TrimSpace(result.Stderr)
		if msg == "" {
			msg = "adb push returned exit code " + fmt.Sprintf("%d", result.ExitCode)
		}
		return fmt.Errorf("push DEX to device: %s", msg)
	}
	return nil
}

// Start launches the monitor DEX on the device via app_process.
// interval is the poll interval in milliseconds. count is the number of polls (0 = unlimited).
func Start(ctx context.Context, client *adb.Client, interval int, count int) (*Process, error) {
	args := client.BaseArgs()
	args = append(args, "shell",
		fmt.Sprintf("CLASSPATH=%s", deviceDEXPath),
		"app_process", "/", "ADBClawMonitor",
		"--interval", fmt.Sprintf("%d", interval),
		"--count", fmt.Sprintf("%d", count),
	)

	cmd := exec.CommandContext(ctx, client.ADBPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr // forward device-side logs to host stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start monitor process: %w", err)
	}

	p := &Process{
		cmd:    cmd,
		stdout: stdout,
		lines:  make(chan string, 64),
		done:   make(chan struct{}),
	}

	go p.readLines()
	return p, nil
}

// Lines returns a channel that receives each stdout line from the monitor process.
func (p *Process) Lines() <-chan string {
	return p.lines
}

// Wait waits for the monitor process to exit and returns any error.
func (p *Process) Wait() error {
	<-p.done
	return p.cmd.Wait()
}

// Stop kills the monitor process.
func (p *Process) Stop() {
	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
}

func (p *Process) readLines() {
	defer close(p.lines)
	defer close(p.done)

	scanner := bufio.NewScanner(p.stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			p.lines <- line
		}
	}
}

// ParseLine parses a JSON line from the monitor output into a TextEntry.
func ParseLine(line string) (*TextEntry, error) {
	var entry TextEntry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil, fmt.Errorf("parse monitor line: %w", err)
	}
	return &entry, nil
}
