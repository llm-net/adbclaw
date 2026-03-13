# adb-claw

Android device control CLI for AI agent automation. Pure tool layer — no LLM/Agent logic included.

**Website: [adb-claw.llm.net](https://adb-claw.llm.net)**

adb-claw wraps standard `adb` commands into a structured JSON API that AI agents can reliably call. It handles screen observation, UI element indexing, input injection, navigation, state management, app lifecycle, and file transfer — with zero dependencies on the Android device side.

## Features

- **Unified JSON output** — Every command returns `{ok, command, data, error, duration_ms, timestamp}` with actionable `suggestion` on errors
- **UI element indexing** — Parses UI tree into indexed elements, supports tap/long-press/clear-field by index/id/text
- **Parallel observation** — Screenshot + UI tree captured concurrently with partial failure tolerance
- **Deep link navigation** — `open` command bypasses CJK text input limits via URI intents
- **Smart scroll** — Auto-calculates coordinates from screen size; supports direction, page count, element-scoped scrolling
- **Wait for UI** — Block until an element appears/disappears, replacing fragile sleep/observe loops
- **Screen management** — Status, on/off, unlock, rotation control
- **Full app lifecycle** — List, launch, stop, install, uninstall, clear data
- **Live stream monitoring** — `monitor` connects to Android accessibility framework, reads UI text even during video playback where `uiautomator dump` fails
- **Shell & file transfer** — Escape hatch for raw commands + push/pull files
- **App Profiles** — Pre-built knowledge (deep links, layouts, known issues) for popular apps
- **Multi-device support** — Target specific devices via `-s <serial>`
- **Minimal device footprint** — Nearly all operations use pure `adb` commands; only `monitor` pushes a temporary ~7KB DEX helper that exits when done

## Install

### Pre-built binaries (recommended)

Pre-built binaries are available for macOS and Linux (amd64 / arm64). No Go toolchain needed.

```bash
curl -fsSL https://github.com/llm-net/adb-claw/releases/latest/download/install.sh | bash
```

Or download directly from [GitHub Releases](https://github.com/llm-net/adb-claw/releases):

| Platform | Binary |
|----------|--------|
| macOS Apple Silicon | `adb-claw-darwin-arm64` |
| macOS Intel | `adb-claw-darwin-amd64` |
| Linux x86_64 | `adb-claw-linux-amd64` |
| Linux ARM64 | `adb-claw-linux-arm64` |

### Build from source

```bash
cd src
make build    # outputs to bin/adb-claw (project root)
```

Requires Go 1.24+.

### Prerequisites

- **adb** (Android Debug Bridge) installed and in PATH
  - macOS: `brew install android-platform-tools`
  - Linux: `sudo apt install android-tools-adb`
- Android device connected via USB with **USB debugging enabled**

```bash
adb-claw doctor    # verify setup
```

## Use as AI Skill

adb-claw is published as an AI Skill on two platforms, sharing the same Skill definition (`skills/adb-claw/SKILL.md`).

### Claude Code

Install the plugin, then just talk to Claude — no slash commands needed.

```bash
# Install from Plugin Marketplace
claude plugin add llm-net/adb-claw
```

After installation, Claude will automatically use adb-claw when you ask about Android devices. Examples:

```
"Take a screenshot of my Android device"
"Open Douyin and search for 猫咪"
"Tap the Login button on screen"
"Monitor the live stream chat for 30 seconds"
"Install this APK on my phone"
```

The plugin's SessionStart hook downloads the binary on first use. As long as adb is installed and a device is connected, everything works out of the box.

### OpenClaw

Install from ClawHub, then use naturally in conversation with your OpenClaw agent.

```bash
# Install from ClawHub
claw install adb-claw
```

The same natural-language triggers apply — ask your agent to control an Android device and it will invoke adb-claw commands automatically.

### How Triggering Works

Both platforms use the **Triggers** list in `SKILL.md` to decide when to activate the skill. When your message matches any trigger (e.g., mentions tapping, screenshots, Android automation, live stream monitoring), the agent loads the skill and gains access to all adb-claw commands. No explicit invocation is needed — just describe what you want to do with the Android device.

## Commands

```
adb-claw
├── observe [--width px]                        # Screenshot + UI tree (primary command)
├── screenshot [--file path] [--width px]       # Screenshot only
├── ui tree                                     # UI element tree
├── ui find --text/--id/--index                 # Find UI elements
├── tap <x> <y> | --index/--id/--text           # Tap
├── long-press <x> <y> [--duration ms]          # Long press
├── swipe <x1> <y1> <x2> <y2> [--duration ms]  # Swipe
├── key <HOME|BACK|ENTER|...>                   # Key event (30+ aliases)
├── type <text>                                 # Input text (ASCII only)
├── clear-field [--index/--id/--text]           # Clear input field
├── open <uri>                                  # Open URI / deep link
├── scroll <up|down|left|right>                 # Smart scroll
│   [--index N] [--pages N] [--distance px]
├── wait --text/--id/--activity                 # Wait for UI state
│   [--gone] [--timeout ms] [--interval ms]
├── monitor [--duration ms] [--interval ms]     # Continuous UI text monitoring
│   [--stream]                                  #   (accessibility framework)
├── screen status|on|off|unlock|rotation        # Screen management
├── app list|current|launch|stop                # App management
├── app install|uninstall|clear                 # App lifecycle
├── shell <command>                             # Raw shell command
├── file push|pull                              # File transfer
├── device list|info                            # Device info
├── doctor                                      # Environment check
└── skill                                       # Output skill.json
```

## Usage

### Observe & Interact

```bash
# Screenshot + UI tree (always start here)
adb-claw observe --width 540

# Tap by element index (preferred) or coordinates
adb-claw tap --index 5
adb-claw tap --text "Login"
adb-claw tap 540 960

# Type text, press keys
adb-claw type "hello world"
adb-claw key ENTER
adb-claw key BACK

# Clear an input field then retype
adb-claw clear-field --index 7
adb-claw type "new text"
```

### Navigate

```bash
# Smart scroll (auto-calculates coordinates)
adb-claw scroll down
adb-claw scroll up --pages 3
adb-claw scroll down --index 5    # within a specific scrollable element

# Open deep links (key for CJK text — bypasses input text limits)
adb-claw open "snssdk1128://search/result?keyword=猫咪"
adb-claw open "https://www.google.com"

# Wait for UI state instead of sleep+observe polling
adb-claw wait --text "Welcome" --timeout 10000
adb-claw wait --text "Loading" --gone
adb-claw wait --activity ".MainActivity"
```

### Monitor (Live Streams & Video)

```bash
# Read UI text via accessibility framework (works during video playback)
adb-claw monitor                            # 10s bounded, JSON envelope
adb-claw monitor --duration 30000           # 30s bounded
adb-claw monitor --stream --duration 60000  # 60s streaming, JSON lines
adb-claw monitor --interval 1000            # Faster polling (1s)
```

### Screen & App Management

```bash
# Screen control
adb-claw screen status
adb-claw screen on
adb-claw screen unlock
adb-claw screen rotation auto

# App lifecycle
adb-claw app current
adb-claw app launch com.example.app
adb-claw app stop com.example.app
adb-claw app install ./app.apk --replace
adb-claw app uninstall com.example.app
adb-claw app clear com.example.app

# Shell & file transfer
adb-claw shell "settings get system screen_brightness"
adb-claw file push ./test.txt /sdcard/Download/
adb-claw file pull /sdcard/photo.jpg ./
```

## Output Format

All commands return a JSON envelope:

```json
{
  "ok": true,
  "command": "tap",
  "data": { ... },
  "duration_ms": 45,
  "timestamp": "2026-03-12T10:00:00Z"
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
    "suggestion": "Try 'adb-claw ui tree' to see available elements"
  }
}
```

Output format can be changed with `-o`:

```bash
adb-claw observe -o text       # Human-readable
adb-claw tap 100 200 -o quiet  # Errors only
```

## Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-s, --serial` | Target device serial | Auto-detect |
| `-o, --output` | Output format: `json`, `text`, `quiet` | `json` |
| `--timeout` | Command timeout in ms | `30000` |
| `--verbose` | Debug output to stderr | `false` |

## App Profiles

App Profiles are pre-built knowledge bases for specific apps — deep links, UI layouts, device-specific behavior, and known issues. They let agents skip the trial-and-error exploration phase.

Available profiles in `skills/apps/`:

| App | File | Content |
|-----|------|---------|
| Douyin (抖音) | `douyin.md` | Search/user/live deep links, feed/search/profile layouts, Phone vs Pad differences, live stream chat monitoring |
| Meituan (美团) | `meituan.md` | Search/waimai deep links, homepage/menu/search layouts, WebView workarounds, popup chain handling |

Usage:
1. `adb-claw app current` → get package name
2. Check `skills/apps/` for a matching profile
3. Has profile → use deep links and known layouts
4. No profile → `observe` + explore

Contributions welcome — see `skills/apps/README.md` for the profile spec.

## Agent Workflow

1. **Observe first** — Always `observe` before deciding an action
2. **Prefer index** — Use `--index` over coordinates for cross-device reliability
3. **Scroll, don't swipe** — `scroll down` over manual `swipe` coordinates
4. **Wait, don't poll** — `wait --text "Done"` over sleep+observe loops
5. **Deep link for CJK** — `open 'app://search?keyword=中文'` instead of `type`
6. **Clear before type** — `clear-field` then `type` for input fields
7. **Monitor for video/live** — Use `monitor` instead of `ui tree` when video is playing
8. **Check App Profiles** — Load profile before exploring unfamiliar apps
9. **Error recovery** — If action fails, re-observe, handle dialogs/permissions, retry

## Architecture

```
src/
├── cmd/                  # CLI commands (Cobra)
│   ├── root.go           # Root + global flags
│   ├── observe.go        # observe / screenshot
│   ├── ui.go             # ui tree / find
│   ├── input.go          # tap / long-press / swipe / key / type
│   ├── clearfield.go     # clear-field
│   ├── scroll.go         # scroll
│   ├── open.go           # open (deep links)
│   ├── wait.go           # wait (UI conditions)
│   ├── monitor.go        # monitor (accessibility-based text monitoring)
│   ├── screen.go         # screen management
│   ├── app.go            # app lifecycle
│   ├── shell.go          # shell command
│   ├── file.go           # file push/pull
│   └── device.go         # device list/info
├── pkg/
│   ├── adb/shell.go      # Commander interface (all ADB calls go through this)
│   ├── input/             # Input injection + scroll + clear-field
│   ├── monitor/           # DEX push + process management + text parsing
│   ├── device/            # Screen status/control
│   ├── observe/           # Screenshot + UI tree parsing
│   └── output/            # JSON envelope formatting
```

Key design decisions:
- **Commander interface** — All packages call ADB through an interface, enabling mock-based testing
- **Input as top-level commands** — `adb-claw tap` instead of `adb-claw input tap`
- **UI tree filtering** — Only indexes elements with text/resource-id/content-desc or clickable/scrollable attributes
- **Partial failure tolerance** — `observe` succeeds if either screenshot or UI tree succeeds
- **Accessibility fallback** — `monitor` uses a temporary DEX helper to read UI text via accessibility framework when `uiautomator dump` fails (video playback, live streams)
- **Minimal device footprint** — Nearly all operations use pure `adb` commands; only `monitor` pushes a ~7KB temporary helper

## License

MIT
