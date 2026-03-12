---
name: adb-claw
version: 0.2.0
description: "Your eyes and hands on Android. See the screen (screenshot + indexed UI tree), interact (tap, swipe, scroll, type, clear-field), navigate via deep links (bypass CJK text input limits), wait for UI state changes instead of polling, manage full app lifecycle (install/uninstall/clear), control screen (on/off/unlock/rotation), run shell commands, and transfer files. Agent-optimized: structured JSON output, indexed element targeting, and App Profiles with pre-built deep links and layouts for popular apps. Zero device-side install — pure ADB."
homepage: https://github.com/llm-net/adbclaw
metadata:
  {
    "openclaw":
      {
        "emoji": "📱",
        "version": "0.2.0",
        "os": ["darwin", "linux"],
        "tags": ["android", "adb", "mobile", "automation", "ui-testing", "device-control", "deep-link", "screenshot"],
        "requires": { "bins": ["adbclaw", "adb"] },
        "install":
          [
            {
              "id": "adbclaw-curl",
              "kind": "script",
              "script": "curl -fsSL https://github.com/llm-net/adbclaw/releases/latest/download/install.sh | bash",
              "bins": ["adbclaw"],
              "label": "Install adbclaw (curl)",
            },
            {
              "id": "adb-brew",
              "kind": "brew",
              "formula": "android-platform-tools",
              "bins": ["adb"],
              "label": "Install ADB (brew)",
            },
          ],
      },
  }
---

# ADB Claw — Android Device Control

Your eyes and hands on Android. See what's on screen, tap any element, scroll through content, open deep links, wait for UI changes, manage apps, and more — all through a single CLI with structured JSON output.

## Why ADB Claw

- **Observe → Act → Verify loop** — `observe` returns screenshot + indexed UI tree in one call; use element indices to target precisely across any screen size
- **Deep links bypass CJK limits** — `adb input text` can't type Chinese/Japanese/Korean; `adbclaw open 'app://search?keyword=中文'` can
- **Wait, don't poll** — `wait --text "Done"` blocks until the UI element appears, replacing fragile sleep/observe loops
- **Smart scroll** — auto-calculates swipe coordinates from screen size; supports direction, page count, and scrolling within specific elements
- **App Profiles** — pre-built knowledge (deep links, layouts, known issues) for popular apps like Douyin; load once, skip trial-and-error
- **Full app lifecycle** — install, launch, stop, uninstall, clear data — no raw `adb` needed
- **Zero device-side install** — pure ADB commands, nothing installed or running on the Android device
- **Agent-optimized JSON** — every command returns `{ok, command, data, error, duration_ms}` with actionable `suggestion` on errors

## Triggers

- User asks to control, interact with, or automate an Android device
- User asks to test a mobile app or UI on Android
- User mentions tapping, swiping, scrolling, screenshots, or app management on Android
- User wants to open a URL, deep link, or specific app screen on a connected device
- User wants to wait for UI elements to appear/disappear on Android
- User wants to manage screen state (on/off/unlock/rotation) on Android
- User wants to push/pull files to/from an Android device
- User wants to run shell commands on an Android device

## Binary

The adbclaw binary is located at `${CLAUDE_PLUGIN_ROOT}/bin/adbclaw`.

If the binary is not found, the SessionStart hook will automatically download it. You can also run manually:

```bash
bash "${CLAUDE_PLUGIN_ROOT}/scripts/setup.sh"
```

## Setup

Requires two binaries:

1. **adbclaw** — the control CLI
2. **adb** — Android Debug Bridge (from Android SDK Platform-Tools)

### Install adbclaw

Pre-built binaries are available for **macOS** and **Linux** (amd64 / arm64). No Go toolchain needed.

```bash
# One-line install (auto-detects OS and architecture)
curl -fsSL https://github.com/llm-net/adbclaw/releases/latest/download/install.sh | bash
```

Or download a specific binary directly from [GitHub Releases](https://github.com/llm-net/adbclaw/releases):

| Platform | Binary |
|----------|--------|
| macOS Apple Silicon (M1/M2/M3/M4) | `adbclaw-darwin-arm64` |
| macOS Intel | `adbclaw-darwin-amd64` |
| Linux x86_64 | `adbclaw-linux-amd64` |
| Linux ARM64 | `adbclaw-linux-arm64` |

Each release includes `checksums.txt` for SHA256 verification.

### Install adb

```bash
# macOS
brew install android-platform-tools

# Linux (Debian/Ubuntu)
sudo apt install android-tools-adb
```

### Connect device

The Android device must have **USB debugging enabled** and be connected via USB.

```bash
# Verify everything is working
adbclaw doctor
```

## Quick Start

The core loop is **observe → decide → act → observe**:

```bash
# 1. See what's on screen
adbclaw observe --width 540

# 2. Act on what you see (use element index from observe output)
adbclaw tap --index 3

# 3. Verify the result
adbclaw observe --width 540
```

For CJK apps, use deep links to bypass text input limits:

```bash
# Search in Douyin (Chinese TikTok) — no manual typing needed
adbclaw open 'snssdk1128://search/result?keyword=猫咪'

# Wait for results to load
adbclaw wait --text "综合" --timeout 5000
```

## App Profiles

App Profiles are pre-built knowledge bases for specific apps — deep links, UI layouts, device-specific behavior, and known issues. They dramatically reduce the trial-and-error needed to automate an app.

**Available Profiles**: `skills/apps/` directory

| App | File | Key Content |
|-----|------|-------------|
| Douyin (抖音) | `douyin.md` | Search/user/live deep links, feed/search/profile layouts, Phone vs Pad differences, UI dump workarounds |

**Usage**:
1. `adbclaw app current` → get foreground app package name
2. Check `skills/apps/` for a matching Profile
3. Has Profile → use deep links and known layouts (fast path)
4. No Profile → `observe` + explore (slow path)
5. Check device form factor: `adbclaw device info` → short edge < 1200px = Phone, >= 1200px = Pad

Profiles are plain Markdown files. New app support = drop a `.md` file into `skills/apps/`.

## Global Flags

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
adbclaw observe              # Default
adbclaw observe --width 540  # Scale screenshot width
```

Returns: base64 PNG screenshot, indexed UI elements with text/id/bounds/center coordinates.

### screenshot — Capture Screen

```bash
adbclaw screenshot                      # Returns base64 PNG in JSON
adbclaw screenshot -f output.png        # Save to file
adbclaw screenshot --width 540          # Scale down
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
adbclaw swipe 540 1800 540 600 --duration 500  # Slow swipe
```

### type — Input Text (ASCII only)

```bash
adbclaw type "Hello world"
```

**Important**: Only ASCII text is supported. For CJK/emoji input, use `open` with deep links (e.g., `adbclaw open 'myapp://search?keyword=中文'`).

### key — Press System Key

```bash
adbclaw key HOME        # Home screen
adbclaw key BACK        # Navigate back
adbclaw key ENTER       # Confirm / submit
adbclaw key TAB         # Next field
adbclaw key DEL         # Delete character
adbclaw key POWER       # Power button
adbclaw key VOLUME_UP   # Volume up
adbclaw key VOLUME_DOWN # Volume down
adbclaw key PASTE       # Paste from clipboard
adbclaw key COPY        # Copy selection
adbclaw key CUT         # Cut selection
adbclaw key WAKEUP      # Wake screen
adbclaw key SLEEP       # Sleep screen
```

### clear-field — Clear Input Field

Clear text in the currently focused input field. Optionally tap an element first to focus it.

```bash
adbclaw clear-field                   # Clear focused field
adbclaw clear-field --index 5         # Focus element #5 then clear
adbclaw clear-field --id "input_name" # Focus by resource ID then clear
adbclaw clear-field --text "Username" # Focus by text then clear
```

Uses Ctrl+A+DEL on SDK 31+, falls back to repeated DEL on older devices.

### open — Open URI (Deep Link)

Open any URI using Android's ACTION_VIEW intent. The key to CJK text input — pass Chinese/Japanese/Korean text as URL parameters in deep links.

```bash
adbclaw open https://www.google.com
adbclaw open myapp://path/to/screen
adbclaw open "market://details?id=com.example"
adbclaw open "snssdk1128://search/result?keyword=猫咪"   # Douyin search in Chinese
```

### scroll — Smart Scroll

Scroll the screen or a specific scrollable element. Auto-calculates swipe coordinates from screen size — no manual coordinate math needed.

```bash
adbclaw scroll down                  # Scroll down one page
adbclaw scroll up                    # Scroll up one page
adbclaw scroll down --pages 3        # Scroll down 3 pages
adbclaw scroll down --index 5        # Scroll within element #5
adbclaw scroll left --distance 500   # Scroll left 500 pixels
```

**Always prefer `scroll` over manual `swipe` for page navigation.**

### wait — Wait for UI Condition

Wait for a UI element or activity to appear or disappear. Replaces fragile sleep+observe polling loops with a single blocking call.

```bash
adbclaw wait --text "Login"                 # Wait for text to appear
adbclaw wait --id "btn_submit"              # Wait for element by ID
adbclaw wait --text "Loading" --gone        # Wait for text to disappear
adbclaw wait --activity ".MainActivity"     # Wait for activity
adbclaw wait --text "Done" --timeout 20000  # Custom timeout (20s)
```

Default timeout: 10s. Default poll interval: 800ms.

### screen — Screen Management

```bash
adbclaw screen status               # Display on/off, locked, rotation
adbclaw screen on                   # Wake up screen
adbclaw screen off                  # Turn off screen
adbclaw screen unlock               # Wake + swipe unlock (no password)
adbclaw screen rotation auto        # Enable auto-rotation
adbclaw screen rotation 0           # Portrait
adbclaw screen rotation 1           # Landscape
```

### app — App Management

```bash
adbclaw app list         # Third-party apps
adbclaw app list --all   # Include system apps
adbclaw app current      # Current foreground app
adbclaw app launch <pkg> # Launch app by package name
adbclaw app stop <pkg>   # Force stop app
adbclaw app install <apk> [--replace]  # Install APK
adbclaw app uninstall <pkg>            # Uninstall app
adbclaw app clear <pkg>               # Clear app data/cache
```

### shell — Run Raw Shell Command

Escape hatch for anything `adbclaw` doesn't have a dedicated command for.

```bash
adbclaw shell "ls /sdcard/"
adbclaw shell "getprop ro.build.version.release"
adbclaw shell "settings put system screen_brightness 128"
```

Returns stdout, stderr, and exit_code in JSON envelope.

### file — File Transfer

```bash
adbclaw file push ./local.apk /sdcard/      # Push to device
adbclaw file pull /sdcard/photo.jpg ./       # Pull from device
```

### device — Device Info

```bash
adbclaw device list      # List connected devices
adbclaw device info      # Model, Android version, screen size, density
```

### ui — UI Element Inspection

```bash
adbclaw ui tree                    # Full UI element tree
adbclaw ui find --text "Settings"  # Find by text
adbclaw ui find --id "com.app:id/title"  # Find by resource ID
adbclaw ui find --index 3          # Find by index
```

## Workflow Patterns

### Always Observe First

Before any action, run `observe` to see the screen. After every action, `observe` again to verify.

```
1. adbclaw observe          → See what's on screen
2. adbclaw tap --index 3    → Perform action
3. adbclaw observe          → Verify result
```

### Prefer Index-Based Targeting

Use `--index N` over coordinates. Indices from `observe` are stable across screen sizes.

### Type After Focus

Always tap an input field first, then type:

```
1. adbclaw tap --index 7       → Focus the text field
2. adbclaw type "search query" → Enter text
3. adbclaw key ENTER           → Submit
```

### Scroll Pattern

```
adbclaw scroll down             # Scroll down one page
adbclaw scroll up --pages 3     # Scroll up multiple pages
adbclaw scroll down --index 5   # Scroll within a specific element
```

**Always prefer `scroll` over manual `swipe`.** After scrolling, always `observe`.

### CJK Text Input

`adbclaw type` only supports ASCII. For Chinese/Japanese/Korean input:

1. **Preferred**: Use `adbclaw open` with deep links (e.g., `adbclaw open 'myapp://search?keyword=中文'`)
2. **Clear before input**: `adbclaw clear-field --index 7` to clear existing text first

### Wait for Page Load

Instead of repeated `observe` polling, use `wait`:

```
1. adbclaw tap --index 3                      → Trigger navigation
2. adbclaw wait --text "Welcome" --timeout 15000  → Wait for new page
3. adbclaw observe                             → See the loaded page
```

### Device Form Factor Detection

Use `adbclaw device info` to get screen size, then determine form factor:
- Short edge < 1200px → **Phone** (portrait-first)
- Short edge >= 1200px → **Pad/Fold** (landscape-first)

Swipe coordinates and UI layouts differ between Phone and Pad. App Profiles document these differences.

### Error Recovery

If an action fails or produces unexpected results:
1. Run `observe` to see the current state
2. Check if the screen changed unexpectedly (dialog, permission prompt)
3. Adjust and retry

## Output Format

All commands return JSON:

```json
{
  "ok": true,
  "command": "tap",
  "data": { ... },
  "duration_ms": 150,
  "timestamp": "2026-03-12T10:00:00Z"
}
```

On error:

```json
{
  "ok": false,
  "error": {
    "code": "DEVICE_NOT_FOUND",
    "message": "No device connected",
    "suggestion": "Connect a device via USB and enable USB debugging"
  }
}
```

## Troubleshooting

| Problem | Solution |
|---------|----------|
| No devices found | Connect device via USB with USB debugging enabled |
| adb not found | `brew install android-platform-tools` (macOS) |
| Tap hits wrong element | Use `--index` instead of coordinates; re-run `observe` |
| `type` doesn't work | Tap input field first to focus; ASCII only |
| CJK text needed | Use `adbclaw open` with deep link containing the text as URL parameter |
| UI dump fails | Pause animations (tap to pause video), wait 1s, retry |
| UI dump fails on search pages | Search results may auto-play video previews; use `screenshot` instead |
| Command timeout | Increase with `--timeout 60000` |
| Permission dialog | Use `observe` to see it, tap the allow/skip button |
| Screen is off | `adbclaw screen on` or `adbclaw screen unlock` |
