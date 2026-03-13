package device

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llm-net/adb-claw/pkg/adb"
	"github.com/llm-net/adb-claw/pkg/input"
)

// ScreenStatus holds the current screen state.
type ScreenStatus struct {
	Display  string `json:"display"`  // "on" or "off"
	Locked   bool   `json:"locked"`
	Rotation int    `json:"rotation"` // 0-3
}

// GetScreenStatus returns the current display state.
func GetScreenStatus(cmd adb.Commander) (*ScreenStatus, error) {
	result, err := cmd.Shell("dumpsys", "power")
	if err != nil {
		return nil, fmt.Errorf("dumpsys power failed: %w", err)
	}

	status := &ScreenStatus{Display: "off"}

	for _, line := range strings.Split(result.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Display Power") {
			if strings.Contains(line, "state=ON") {
				status.Display = "on"
			}
		}
		if strings.Contains(line, "mWakefulness=") {
			if strings.Contains(line, "Awake") {
				status.Display = "on"
			}
		}
	}

	// Check lock status
	winResult, err := cmd.Shell("dumpsys", "window")
	if err == nil {
		for _, line := range strings.Split(winResult.Stdout, "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "mDreamingLockscreen") && strings.Contains(line, "true") {
				status.Locked = true
			}
			if strings.Contains(line, "isStatusBarKeyguard") && strings.Contains(line, "true") {
				status.Locked = true
			}
			if strings.Contains(line, "mShowingLockscreen") && strings.Contains(line, "true") {
				status.Locked = true
			}
		}
	}

	// Check rotation
	rotResult, err := cmd.Shell("dumpsys", "window", "displays")
	if err == nil {
		for _, line := range strings.Split(rotResult.Stdout, "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "mCurrentRotation") || strings.Contains(line, "cur=") {
				// Parse rotation value (0-3)
				for _, word := range strings.Fields(line) {
					if strings.HasPrefix(word, "mCurrentRotation=") {
						val := strings.TrimPrefix(word, "mCurrentRotation=")
						if r, err := strconv.Atoi(strings.TrimRight(val, ",")); err == nil {
							status.Rotation = r
						}
					}
				}
			}
		}
	}

	return status, nil
}

// ScreenOn wakes up the screen.
func ScreenOn(cmd adb.Commander) error {
	return input.KeyEvent(cmd, "WAKEUP")
}

// ScreenOff turns off the screen.
func ScreenOff(cmd adb.Commander) error {
	return input.KeyEvent(cmd, "SLEEP")
}

// ScreenUnlock wakes the screen and performs a simple swipe-up gesture.
// This works for devices without a password/PIN lock screen.
func ScreenUnlock(cmd adb.Commander) error {
	// Wake up first
	if err := input.KeyEvent(cmd, "WAKEUP"); err != nil {
		return fmt.Errorf("wakeup failed: %w", err)
	}

	// Get screen size for swipe coordinates
	w, h, err := input.GetScreenSize(cmd)
	if err != nil {
		return fmt.Errorf("get screen size failed: %w", err)
	}

	// Swipe up from bottom center
	centerX := w / 2
	bottomY := h * 4 / 5
	centerY := h / 2

	return input.Swipe(cmd, centerX, bottomY, centerX, centerY, 300)
}

// SetRotation sets the screen rotation.
// "auto" enables auto-rotation; 0-3 sets a fixed rotation.
func SetRotation(cmd adb.Commander, mode string) error {
	if strings.ToLower(mode) == "auto" {
		result, err := cmd.Shell("settings", "put", "system", "accelerometer_rotation", "1")
		if err != nil {
			return fmt.Errorf("enable auto rotation failed: %w", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("enable auto rotation failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
		}
		return nil
	}

	rotation, err := strconv.Atoi(mode)
	if err != nil || rotation < 0 || rotation > 3 {
		return fmt.Errorf("invalid rotation %q: use auto, 0, 1, 2, or 3", mode)
	}

	// Disable auto-rotation
	result, err := cmd.Shell("settings", "put", "system", "accelerometer_rotation", "0")
	if err != nil {
		return fmt.Errorf("disable auto rotation failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("disable auto rotation failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}

	// Set user rotation
	result, err = cmd.Shell("settings", "put", "system", "user_rotation", mode)
	if err != nil {
		return fmt.Errorf("set rotation failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("set rotation failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}

	return nil
}
