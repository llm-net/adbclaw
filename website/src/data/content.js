export const features = [
  {
    title: 'Structured JSON API',
    description:
      'Every command returns {ok, command, data, error, duration_ms, timestamp}. Parse responses reliably. Errors include codes and actionable suggestions. Three output modes: json, text, quiet.',
    icon: 'json',
  },
  {
    title: 'Smart Element Targeting',
    description:
      'Tap, long-press, or scroll by element index, resource ID, or text content. No coordinate guessing. The UI tree is indexed with bounds and center points for pixel-accurate actions.',
    icon: 'grid',
  },
  {
    title: 'Deep Link Navigation',
    description:
      'Jump directly to any app screen via URI. Open web pages, trigger WeChat scan, search Taobao — skip multi-step navigation entirely. One command, instant arrival.',
    icon: 'link',
  },
  {
    title: 'Wait for State Changes',
    description:
      'Block until a UI element appears or disappears, or an Activity loads. No polling loops in your agent code. Configurable timeout and interval. Returns the matched element on success.',
    icon: 'clock',
  },
  {
    title: 'Zero Device Setup',
    description:
      'Pure ADB-based control. No APK to install, no accessibility service, no permissions dialogs. Connect a device over USB or WiFi and start automating immediately.',
    icon: 'zap',
  },
  {
    title: 'Pre-built Binaries',
    description:
      'Download a single compiled binary for your platform — darwin-arm64, darwin-amd64, linux-arm64, or linux-amd64. No Go toolchain required. curl | bash one-liner install available.',
    icon: 'package',
  },
  {
    title: 'Full Device Control',
    description:
      '30+ commands covering screen observation, input injection, smart scrolling, app lifecycle, screen management, shell access, and file transfer. Everything over standard ADB.',
    icon: 'layers',
  },
  {
    title: 'App Knowledge Profiles',
    description:
      'Pre-built profiles for popular apps (Douyin, WeChat, etc.) with deep links, UI layouts, and known issues. Your agent gets expert-level app knowledge out of the box.',
    icon: 'book',
  },
  {
    title: 'Agent-First Design',
    description:
      'Built as a Skill for AI agents — available on Claude Code and OpenClaw/ClawHub. Machine-readable skill descriptions. Designed for LLM consumption, usable by humans.',
    icon: 'bot',
  },
]

export const cliExamples = [
  {
    title: 'Observe & Inspect',
    commands: [
      { cmd: 'adbclaw observe', comment: 'Screenshot + UI tree' },
      { cmd: 'adbclaw screenshot --width 720', comment: 'Downscaled capture' },
      { cmd: 'adbclaw ui tree', comment: 'Indexed elements' },
      { cmd: 'adbclaw ui find --text "Login"', comment: 'Find by text' },
    ],
  },
  {
    title: 'Input & Navigate',
    commands: [
      { cmd: 'adbclaw tap --index 5', comment: 'Tap by element index' },
      { cmd: 'adbclaw type "hello world"', comment: 'Input text' },
      { cmd: 'adbclaw scroll down --pages 3', comment: 'Smart scroll' },
      { cmd: 'adbclaw open "weixin://dl/scan"', comment: 'Deep link' },
      { cmd: 'adbclaw clear-field --index 2', comment: 'Clear input' },
    ],
  },
  {
    title: 'Wait & Screen',
    commands: [
      { cmd: 'adbclaw wait --text "Done"', comment: 'Wait for element' },
      { cmd: 'adbclaw wait --text "Loading" --gone', comment: 'Wait until gone' },
      { cmd: 'adbclaw screen status', comment: 'On/off/lock/rotation' },
      { cmd: 'adbclaw screen unlock', comment: 'Wake + swipe unlock' },
    ],
  },
  {
    title: 'Apps & System',
    commands: [
      { cmd: 'adbclaw app launch com.example', comment: 'Launch app' },
      { cmd: 'adbclaw app install ./app.apk', comment: 'Install APK' },
      { cmd: 'adbclaw shell "pm list packages"', comment: 'Raw shell' },
      { cmd: 'adbclaw file pull /sdcard/log.txt .', comment: 'Pull file' },
    ],
  },
]

export const architectureSteps = [
  {
    label: 'AI Agent',
    sublabel: 'Claude / OpenClaw / LLM',
    description: 'Reads skill description, sends structured commands, parses JSON responses to decide next actions',
  },
  {
    label: 'adbclaw',
    sublabel: 'Go CLI · v1.3.0',
    description: 'Translates 30+ commands to ADB operations. Returns structured JSON with error codes and suggestions',
  },
  {
    label: 'ADB',
    sublabel: 'USB / WiFi',
    description: 'Transports shell commands, screenshots, and file transfers to the Android device',
  },
  {
    label: 'Device',
    sublabel: 'Android',
    description: 'Executes operations — UI dumps, screenshots, input events, app management, file I/O',
  },
]

export const relatedProjects = [
  {
    name: 'OpenClaw',
    url: 'https://github.com/openclaw/openclaw',
    description: 'Local-first personal AI assistant platform. adbclaw is published as an OpenClaw Skill on ClawHub.',
    stars: '',
    highlight: true,
  },
  {
    name: 'mobile-use',
    url: 'https://github.com/anthropics/mobile-use',
    description: 'Anthropic\'s AI agent for controlling real mobile devices, top performer on AndroidWorld benchmark.',
    stars: '2.2k',
  },
  {
    name: 'DroidRun',
    url: 'https://github.com/droidrun/droidrun',
    description: 'LLM-powered Android device control framework with multi-model support.',
    stars: '7.7k',
  },
  {
    name: 'scrcpy',
    url: 'https://github.com/Genymobile/scrcpy',
    description: 'The gold standard for Android screen mirroring and control.',
    stars: '136k',
  },
]

export const agentWorkflow = [
  { step: '01', action: 'Observe', detail: 'Run observe to capture screenshot + UI tree in one call' },
  { step: '02', action: 'Decide', detail: 'AI agent analyzes screen state and plans the next action' },
  { step: '03', action: 'Act', detail: 'Tap, scroll, open deep link, or type — by element index' },
  { step: '04', action: 'Wait', detail: 'Use wait to block until UI state changes, then re-observe' },
]

export const commandStats = {
  totalCommands: '30+',
  platforms: '4',
  outputFormats: '3',
  zeroDeviceApps: true,
}
