# adbclaw

Android device control CLI for AI agent automation. Pure tool layer — no LLM/Agent logic included.

adbclaw wraps `adb shell` commands into a structured JSON API that AI agents can reliably call. It handles screen observation, UI element indexing, input injection, and app management.

## Features

- **Unified JSON output** — Every command returns `{ok, command, data, error, duration_ms, timestamp}`
- **UI element indexing** — Parses UI tree into indexed elements, supports tap-by-index/id/text
- **Parallel observation** — Screenshot + UI tree captured concurrently with partial failure tolerance
- **Multi-device support** — Target specific devices via `-s <serial>`
- **AI-native design** — `skill` command outputs machine-readable capability description

## Quick Start

### Prerequisites

- Go 1.24+
- ADB (Android Debug Bridge) installed and in PATH
- Android device connected via USB or WiFi with USB debugging enabled

### Build

```bash
cd src
make build
```

Binary outputs to `src/bin/adbclaw`.

### Verify Setup

```bash
adbclaw doctor
```

Checks ADB installation, device connection, and device capabilities.

## Usage

### Observe Screen

```bash
# Screenshot + UI tree (recommended first step)
adbclaw observe

# Screenshot only, downscaled to 720px width
adbclaw screenshot --width 720

# Save screenshot to file
adbclaw screenshot --file screen.png
```

### UI Inspection

```bash
# Get all interactive UI elements with index numbers
adbclaw ui tree

# Find elements by text or resource ID
adbclaw ui find --text "Login"
adbclaw ui find --id "com.example:id/btn_login"

# Get specific element by index
adbclaw ui find --index 5
```

### Input

```bash
# Tap by coordinates
adbclaw tap 540 960

# Tap by UI element (recommended)
adbclaw tap --index 5
adbclaw tap --text "Login"
adbclaw tap --id "com.example:id/btn_login"

# Long press
adbclaw long-press 540 960 --duration 1500

# Swipe (scroll down)
adbclaw swipe 540 1800 540 600

# Key events
adbclaw key HOME
adbclaw key BACK

# Type text (ASCII only)
adbclaw type "hello world"
```

### App Management

```bash
# List third-party apps
adbclaw app list

# List all apps (including system)
adbclaw app list --all

# Current foreground app
adbclaw app current

# Launch / stop app
adbclaw app launch com.example.app
adbclaw app stop com.example.app
```

### Device Info

```bash
# List connected devices
adbclaw device list

# Device details (model, Android version, screen size, etc.)
adbclaw device info
```

## Output Format

All commands return a JSON envelope:

```json
{
  "ok": true,
  "command": "tap",
  "data": { ... },
  "duration_ms": 45,
  "timestamp": "2025-03-01T10:00:00.123Z"
}
```

Error responses include actionable suggestions:

```json
{
  "ok": false,
  "command": "tap",
  "error": {
    "code": "ELEMENT_NOT_FOUND",
    "message": "No element found with text 'Login'",
    "suggestion": "Try 'adbclaw ui tree' to see available elements"
  }
}
```

Output format can be changed with `-o`:

```bash
adbclaw observe -o text    # Human-readable
adbclaw tap 100 200 -o quiet  # Errors only
```

## Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-s, --serial` | Target device serial | Auto-detect |
| `-o, --output` | Output format: `json`, `text`, `quiet` | `json` |
| `--timeout` | Command timeout in ms | `30000` |
| `--verbose` | Debug output to stderr | `false` |

## Agent Integration

adbclaw provides a `skill` command that outputs a machine-readable JSON describing all available tools, their parameters, and workflow hints:

```bash
adbclaw skill
```

### Recommended Agent Workflow

1. **Observe first** — Always run `observe` before deciding an action
2. **Prefer index** — Use `--index` over coordinates for reliability
3. **Type after focus** — Tap an input field first, then use `type`
4. **Scroll pattern** — `swipe 540 1800 540 600` to scroll down
5. **Error recovery** — If action fails, re-observe before retrying

## Architecture

```
cmd/          CLI commands (Cobra)
pkg/adb/      Commander interface — all ADB calls go through this
pkg/input/    Input injection (tap, swipe, key, type)
pkg/observe/  Screen capture + UI tree parsing
pkg/output/   JSON envelope formatting
```

Key design decisions:
- **Commander interface** — All packages call ADB through an interface, enabling mock-based testing
- **Input as top-level commands** — `adbclaw tap` instead of `adbclaw input tap`
- **UI tree filtering** — Only indexes elements with text/resource-id/content-desc or clickable/scrollable attributes
- **Partial failure tolerance** — `observe` succeeds if either screenshot or UI tree succeeds

## Roadmap

Currently at **Phase 1 MVP** (pure `adb shell` implementation).

See [`docs/adbclaw-technical-plan.md`](docs/adbclaw-technical-plan.md) for the full roadmap including:
- Phase 2: Device-side service (adbclawd), sendevent input, JPEG screenshots
- Phase 3: Input humanization, gesture support, MCP server
- Phase 4: Deep stealth mode, advanced observation

## License

MIT
