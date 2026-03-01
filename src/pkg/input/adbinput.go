package input

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/llm-net/adbclaw/pkg/adb"
)

// Tap sends a tap at the given coordinates.
func Tap(cmd adb.Commander, x, y int) error {
	result, err := cmd.Shell("input", "tap", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y))
	if err != nil {
		return fmt.Errorf("tap failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("tap failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}
	return nil
}

// LongPress sends a long press at the given coordinates for the given duration (ms).
func LongPress(cmd adb.Commander, x, y, durationMs int) error {
	// Long press is implemented as a swipe from (x,y) to (x,y) with duration
	result, err := cmd.Shell("input", "swipe",
		fmt.Sprintf("%d", x), fmt.Sprintf("%d", y),
		fmt.Sprintf("%d", x), fmt.Sprintf("%d", y),
		fmt.Sprintf("%d", durationMs))
	if err != nil {
		return fmt.Errorf("long-press failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("long-press failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}
	return nil
}

// Swipe sends a swipe gesture.
func Swipe(cmd adb.Commander, x1, y1, x2, y2, durationMs int) error {
	result, err := cmd.Shell("input", "swipe",
		fmt.Sprintf("%d", x1), fmt.Sprintf("%d", y1),
		fmt.Sprintf("%d", x2), fmt.Sprintf("%d", y2),
		fmt.Sprintf("%d", durationMs))
	if err != nil {
		return fmt.Errorf("swipe failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("swipe failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}
	return nil
}

// KeyEvent sends a key event by name (e.g., HOME, BACK, ENTER).
func KeyEvent(cmd adb.Commander, key string) error {
	// Normalize key: if it doesn't start with KEYCODE_, check common aliases
	keycode := resolveKeycode(key)
	result, err := cmd.Shell("input", "keyevent", keycode)
	if err != nil {
		return fmt.Errorf("key event failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("key event failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}
	return nil
}

// TypeText inputs text via "adb shell input text".
// Special characters are escaped for shell safety.
// Returns an error if the text contains non-ASCII characters (CJK, emoji, etc.)
// because "adb shell input text" does not support them.
func TypeText(cmd adb.Commander, text string) error {
	if HasNonASCII(text) {
		return fmt.Errorf("text contains non-ASCII characters (CJK/emoji/etc.) which are not supported by 'adb shell input text'; consider using clipboard-based input instead")
	}
	escaped := escapeForInput(text)
	result, err := cmd.Shell("input", "text", escaped)
	if err != nil {
		return fmt.Errorf("type text failed: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("type text failed: %s", strings.TrimSpace(result.Stderr+result.Stdout))
	}
	return nil
}

// HasNonASCII returns true if the string contains any non-ASCII character.
func HasNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

// escapeForInput escapes text for "adb shell input text".
// The input command requires spaces and special chars to be escaped.
func escapeForInput(text string) string {
	// Characters that need escaping for adb shell input text.
	// Backslash must be first to avoid double-escaping other replacements.
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		" ", "%s",
		"'", "\\'",
		"\"", "\\\"",
		"(", "\\(",
		")", "\\)",
		"&", "\\&",
		"<", "\\<",
		">", "\\>",
		";", "\\;",
		"|", "\\|",
		"$", "\\$",
		"`", "\\`",
		"~", "\\~",
		"!", "\\!",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
	)
	return replacer.Replace(text)
}

// resolveKeycode maps short key names to Android KEYCODE_ constants.
func resolveKeycode(key string) string {
	upper := strings.ToUpper(key)
	aliases := map[string]string{
		"HOME":        "KEYCODE_HOME",
		"BACK":        "KEYCODE_BACK",
		"ENTER":       "KEYCODE_ENTER",
		"TAB":         "KEYCODE_TAB",
		"DEL":         "KEYCODE_DEL",
		"DELETE":      "KEYCODE_DEL",
		"POWER":       "KEYCODE_POWER",
		"VOLUME_UP":   "KEYCODE_VOLUME_UP",
		"VOLUME_DOWN": "KEYCODE_VOLUME_DOWN",
		"MENU":        "KEYCODE_MENU",
		"SEARCH":      "KEYCODE_SEARCH",
		"DPAD_UP":     "KEYCODE_DPAD_UP",
		"DPAD_DOWN":   "KEYCODE_DPAD_DOWN",
		"DPAD_LEFT":   "KEYCODE_DPAD_LEFT",
		"DPAD_RIGHT":  "KEYCODE_DPAD_RIGHT",
		"DPAD_CENTER": "KEYCODE_DPAD_CENTER",
		"APP_SWITCH":  "KEYCODE_APP_SWITCH",
		"CAMERA":      "KEYCODE_CAMERA",
		"SPACE":       "KEYCODE_SPACE",
		"ESCAPE":      "KEYCODE_ESCAPE",
		"RECENTS":     "KEYCODE_APP_SWITCH",
	}
	if code, ok := aliases[upper]; ok {
		return code
	}
	// If already has KEYCODE_ prefix, use as-is
	if strings.HasPrefix(upper, "KEYCODE_") {
		return upper
	}
	// Otherwise, prepend KEYCODE_
	return "KEYCODE_" + upper
}
