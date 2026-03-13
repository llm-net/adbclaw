package input

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llm-net/adb-claw/pkg/adb"
)

// ClearField clears the text in the currently focused input field.
// On SDK >= 31, uses keycombination Ctrl+A then DEL.
// On older SDKs, falls back to MOVE_END + repeated DEL keys.
// Returns the method used: "keycombination" or "repeated_del".
func ClearField(cmd adb.Commander) (string, error) {
	sdk, err := GetSDKLevel(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to detect SDK level: %w", err)
	}

	if sdk >= 31 {
		// Ctrl+A to select all
		if err := KeyCombination(cmd, "KEYCODE_CTRL_LEFT", "KEYCODE_A"); err != nil {
			return "", fmt.Errorf("select all failed: %w", err)
		}
		// DEL to delete selection
		if err := KeyEvent(cmd, "DEL"); err != nil {
			return "", fmt.Errorf("delete failed: %w", err)
		}
		return "keycombination", nil
	}

	// Fallback: move to end, then send 200 DEL keys in a single call
	if err := KeyEvent(cmd, "MOVE_END"); err != nil {
		return "", fmt.Errorf("move to end failed: %w", err)
	}
	delArgs := make([]string, 200)
	for i := range delArgs {
		delArgs[i] = "KEYCODE_DEL"
	}
	result, err := cmd.Shell(append([]string{"input", "keyevent"}, delArgs...)...)
	if err != nil {
		return "", fmt.Errorf("repeated delete failed: %w", err)
	}
	if result.ExitCode != 0 {
		return "", fmt.Errorf("repeated delete failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}
	return "repeated_del", nil
}

// KeyCombination sends a key combination via "adb shell input keycombination".
// Requires Android 12+ (SDK 31+).
func KeyCombination(cmd adb.Commander, keys ...string) error {
	args := append([]string{"input", "keycombination"}, keys...)
	result, err := cmd.Shell(args...)
	if err != nil {
		return fmt.Errorf("keycombination failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("keycombination failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}
	return nil
}

// GetSDKLevel returns the device's SDK API level.
func GetSDKLevel(cmd adb.Commander) (int, error) {
	result, err := cmd.Shell("getprop", "ro.build.version.sdk")
	if err != nil {
		return 0, err
	}
	level, err := strconv.Atoi(strings.TrimSpace(result.Stdout))
	if err != nil {
		return 0, fmt.Errorf("invalid SDK level %q: %w", strings.TrimSpace(result.Stdout), err)
	}
	return level, nil
}
