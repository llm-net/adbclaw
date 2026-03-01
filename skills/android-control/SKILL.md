# Android Device Control

Control Android devices via adbclaw — tap, swipe, type, launch apps, take screenshots, and inspect UI elements.

## Triggers

- User asks to control, interact with, or automate an Android device
- User asks to test a mobile app or UI
- User mentions tapping, swiping, screenshots, or app launching on Android
- User wants to automate Android device operations

## Binary

The adbclaw binary is located at `${CLAUDE_PLUGIN_ROOT}/bin/adbclaw`.

If the binary is not found, the SessionStart hook will automatically download it. You can also run manually:

```bash
bash "${CLAUDE_PLUGIN_ROOT}/scripts/setup.sh"
```

## Global Flags

All commands accept these flags:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--serial` | `-s` | Target device serial (when multiple devices connected) | auto-detect |
| `--output` | `-o` | Output format: `json`, `text`, `quiet` | `json` |
| `--timeout` | | Command timeout in milliseconds | `30000` |
| `--verbose` | | Enable debug output to stderr | `false` |

## Commands

### observe — Screenshot + UI Tree (Primary Command)

Captures screenshot and UI element tree in one call. **Always use this before and after actions.**

```bash
adbclaw observe
```

Returns: base64 PNG screenshot, indexed UI elements with text/id/bounds/center coordinates, and device state.

### screenshot — Capture Screen

```bash
adbclaw screenshot              # Returns base64 PNG in JSON
adbclaw screenshot -f output.png  # Save to file
```

### tap — Tap UI Element

Tap by element index (preferred), resource ID, text, or coordinates:

```bash
adbclaw tap --index 5            # Tap element #5 from observe output
adbclaw tap --id "com.app:id/btn" # Tap by resource ID
adbclaw tap --text "Submit"       # Tap by visible text
adbclaw tap 540 960              # Tap coordinates (x y)
```

**Always prefer `--index` over coordinates.** Index values come from `observe` output.

### long-press — Long Press

```bash
adbclaw long-press 540 960              # Default duration
adbclaw long-press 540 960 --duration 2000  # 2 seconds
```

### swipe — Swipe Gesture

```bash
adbclaw swipe 540 1800 540 600           # Swipe up (scroll down)
adbclaw swipe 540 600 540 1800           # Swipe down (scroll up)
adbclaw swipe 900 960 100 960            # Swipe left
adbclaw swipe 100 960 900 960            # Swipe right
adbclaw swipe 540 1800 540 600 --duration 500  # Slow swipe
```

### type — Input Text

Type text into the currently focused input field. **Tap the field first to focus it.**

```bash
adbclaw type "Hello world"
```

Note: Only ASCII text is supported. CJK characters and emoji are not supported via adb input.

### key — Press System Key

```bash
adbclaw key HOME        # Go to home screen
adbclaw key BACK        # Navigate back
adbclaw key ENTER       # Confirm / submit
adbclaw key TAB         # Next field
adbclaw key DEL         # Delete character
adbclaw key POWER       # Power button
adbclaw key VOLUME_UP   # Volume up
adbclaw key VOLUME_DOWN # Volume down
```

### app launch — Launch App

```bash
adbclaw app launch com.android.settings       # Launch by package name
adbclaw app launch com.app/.MainActivity       # Launch specific activity
```

Use `app list` first to find package names.

### app stop — Force Stop App

```bash
adbclaw app stop com.example.app
```

### app current — Current Foreground App

```bash
adbclaw app current
```

Returns the package name and activity of the currently focused app.

### app list — List Installed Apps

```bash
adbclaw app list         # Third-party apps only
adbclaw app list --all   # Include system apps
```

### device list — List Connected Devices

```bash
adbclaw device list
```

### device info — Device Details

```bash
adbclaw device info
```

Returns: model, brand, Android version, SDK level, screen size, density, supported ABIs.

### ui tree — Full UI Element Tree

```bash
adbclaw ui tree
```

Returns all significant UI elements with index, class, text, resource-id, bounds, and interaction states.

### ui find — Search UI Elements

```bash
adbclaw ui find --text "Settings"
adbclaw ui find --id "com.app:id/title"
adbclaw ui find --index 3
```

## Workflow Patterns

### Always Observe First

Before taking any action, run `observe` to see the current screen state and available UI elements. After every action, run `observe` again to verify the result.

```
1. adbclaw observe          → See what's on screen
2. adbclaw tap --index 3    → Perform action
3. adbclaw observe          → Verify result
```

### Prefer Index-Based Targeting

When tapping elements, prefer `--index N` over coordinates. Index values from `observe` are reliable across different screen sizes and densities.

### Type After Focus

Always tap an input field first, then use `type`:

```
1. adbclaw tap --index 7       → Focus the text field
2. adbclaw type "search query" → Enter text
3. adbclaw key ENTER           → Submit
```

### Scroll Pattern

To scroll through content, use swipe gestures. On a typical 1080x1920 screen:

```
Scroll down:  adbclaw swipe 540 1500 540 500
Scroll up:    adbclaw swipe 540 500 540 1500
```

After scrolling, always `observe` to see newly visible content.

### Error Recovery

If an action fails or produces unexpected results:
1. Run `observe` to see the current state
2. Check if the screen changed unexpectedly (dialog, permission prompt)
3. Adjust and retry

## Output Format

All commands return JSON with a consistent envelope:

```json
{
  "ok": true,
  "command": "tap",
  "data": { ... },
  "duration_ms": 150,
  "timestamp": "2025-03-01T10:00:00Z"
}
```

On error:

```json
{
  "ok": false,
  "command": "tap",
  "error": {
    "code": "DEVICE_NOT_FOUND",
    "message": "No device connected",
    "suggestion": "Connect a device via USB or start an emulator"
  }
}
```

## Troubleshooting

| Problem | Solution |
|---------|----------|
| "no devices/emulators found" | Connect a device via USB with USB debugging enabled, or start an Android emulator |
| "adb not found" | Install Android SDK Platform-Tools: `brew install android-platform-tools` (macOS) |
| Tap hits wrong element | Use `--index` instead of coordinates; re-run `observe` to get fresh indices |
| `type` doesn't work | Ensure an input field is focused first (tap it); only ASCII text is supported |
| Command timeout | Increase with `--timeout 60000`; device may be slow or unresponsive |
| Permission dialog blocks | Use `observe` to see the dialog, then tap the appropriate button |
