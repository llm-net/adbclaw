package adb

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Result holds the output of an ADB command.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Commander is the interface for executing ADB commands.
// All pkg-level code uses this interface, enabling mock-based unit tests.
type Commander interface {
	// Shell runs "adb shell <args...>" and returns text output.
	Shell(args ...string) (*Result, error)
	// ExecOut runs "adb exec-out <args...>" and returns raw bytes (binary-safe).
	ExecOut(args ...string) ([]byte, error)
	// RawCommand runs "adb <args...>" (no shell prefix).
	RawCommand(args ...string) (*Result, error)
}

// Client is the default Commander implementation backed by the adb binary.
type Client struct {
	ADBPath string        // path to adb binary, default "adb"
	Serial  string        // device serial (-s flag), empty = default device
	Timeout time.Duration // command timeout, default 30s
}

// NewClient creates a Client with sensible defaults.
func NewClient(serial string, timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		ADBPath: "adb",
		Serial:  serial,
		Timeout: timeout,
	}
}

// BaseArgs returns the base ADB arguments (e.g., -s serial) for this client.
func (c *Client) BaseArgs() []string {
	if c.Serial != "" {
		return []string{"-s", c.Serial}
	}
	return nil
}

// Shell runs "adb [-s serial] shell <args...>".
func (c *Client) Shell(args ...string) (*Result, error) {
	full := append(c.BaseArgs(), "shell")
	full = append(full, args...)
	return c.run(full...)
}

// ExecOut runs "adb [-s serial] exec-out <args...>" returning raw bytes.
func (c *Client) ExecOut(args ...string) ([]byte, error) {
	full := append(c.BaseArgs(), "exec-out")
	full = append(full, args...)
	return c.runRaw(full...)
}

// RawCommand runs "adb [-s serial] <args...>".
func (c *Client) RawCommand(args ...string) (*Result, error) {
	full := append(c.BaseArgs(), args...)
	return c.run(full...)
}

func (c *Client) run(args ...string) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.ADBPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if ctx.Err() == context.DeadlineExceeded {
		return result, fmt.Errorf("adb command timed out after %v", c.Timeout)
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
		// Many adb shell commands return non-zero for valid reasons;
		// let callers decide whether to treat as error.
		return result, nil
	}
	if err != nil {
		return result, fmt.Errorf("adb exec error: %w", err)
	}
	return result, nil
}

func (c *Client) runRaw(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.ADBPath, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("adb command timed out after %v", c.Timeout)
	}
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return nil, fmt.Errorf("adb exec-out error: %s", errMsg)
		}
		return nil, fmt.Errorf("adb exec-out error: %w", err)
	}
	return stdout.Bytes(), nil
}
