package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// CheckResult represents a single diagnostic check.
type CheckResult struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "ok", "warning", "error"
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment and device readiness",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var checks []CheckResult

		// 1. Check adb binary
		checks = append(checks, checkADB())

		// 2. Check device connection
		deviceCheck := checkDevice()
		checks = append(checks, deviceCheck)

		// Only run device-specific checks if a device is connected
		if deviceCheck.Status == "ok" {
			// 3. Check screencap
			checks = append(checks, checkScreencap())

			// 4. Check uiautomator
			checks = append(checks, checkUIAutomator())

			// 5. Check input command
			checks = append(checks, checkInput())
		}

		// Summary
		allOK := true
		for _, c := range checks {
			if c.Status == "error" {
				allOK = false
				break
			}
		}

		writer.Success("doctor", map[string]interface{}{
			"checks":  checks,
			"healthy": allOK,
		}, start)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func checkADB() CheckResult {
	adbPath := client.ADBPath
	path, err := exec.LookPath(adbPath)
	if err != nil {
		return CheckResult{
			Name:    "adb",
			Status:  "error",
			Message: fmt.Sprintf("adb not found (%s)", adbPath),
			Detail:  "Install Android platform-tools: https://developer.android.com/tools/releases/platform-tools",
		}
	}

	// Get version
	out, err := exec.Command(path, "version").Output()
	if err != nil {
		return CheckResult{
			Name:    "adb",
			Status:  "warning",
			Message: fmt.Sprintf("adb found at %s but version check failed", path),
		}
	}

	version := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	return CheckResult{
		Name:    "adb",
		Status:  "ok",
		Message: version,
		Detail:  path,
	}
}

func checkDevice() CheckResult {
	result, err := client.RawCommand("devices")
	if err != nil {
		return CheckResult{
			Name:    "device",
			Status:  "error",
			Message: "Failed to list devices: " + err.Error(),
		}
	}

	var deviceCount int
	for _, line := range strings.Split(result.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == "device" {
			deviceCount++
		}
	}

	if deviceCount == 0 {
		return CheckResult{
			Name:    "device",
			Status:  "error",
			Message: "No devices connected",
			Detail:  "Connect a device via USB and enable USB debugging",
		}
	}

	return CheckResult{
		Name:    "device",
		Status:  "ok",
		Message: fmt.Sprintf("%d device(s) connected", deviceCount),
	}
}

func checkScreencap() CheckResult {
	result, err := client.Shell("which", "screencap")
	if err != nil {
		return CheckResult{
			Name:   "screencap",
			Status: "error",
			Message: "screencap check failed: " + err.Error(),
		}
	}
	path := strings.TrimSpace(result.Stdout)
	if path == "" {
		return CheckResult{
			Name:    "screencap",
			Status:  "error",
			Message: "screencap not found on device",
		}
	}
	return CheckResult{
		Name:    "screencap",
		Status:  "ok",
		Message: "screencap available",
		Detail:  path,
	}
}

func checkUIAutomator() CheckResult {
	result, err := client.Shell("which", "uiautomator")
	if err != nil {
		return CheckResult{
			Name:    "uiautomator",
			Status:  "error",
			Message: "uiautomator check failed: " + err.Error(),
		}
	}
	path := strings.TrimSpace(result.Stdout)
	if path == "" {
		return CheckResult{
			Name:    "uiautomator",
			Status:  "error",
			Message: "uiautomator not found on device",
		}
	}
	return CheckResult{
		Name:    "uiautomator",
		Status:  "ok",
		Message: "uiautomator available",
		Detail:  path,
	}
}

func checkInput() CheckResult {
	result, err := client.Shell("which", "input")
	if err != nil {
		return CheckResult{
			Name:    "input",
			Status:  "error",
			Message: "input command check failed: " + err.Error(),
		}
	}
	path := strings.TrimSpace(result.Stdout)
	if path == "" {
		return CheckResult{
			Name:    "input",
			Status:  "error",
			Message: "input command not found on device",
		}
	}
	return CheckResult{
		Name:    "input",
		Status:  "ok",
		Message: "input command available",
		Detail:  path,
	}
}
