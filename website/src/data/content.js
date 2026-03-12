export const features = [
  {
    title: 'Unified JSON Output',
    description:
      'Every command returns a structured envelope: {ok, command, data, error, duration_ms, timestamp}. Errors include actionable suggestions. Three output modes: json, text, quiet.',
    icon: 'json',
  },
  {
    title: 'UI Element Indexing',
    description:
      'Parses the Android UI tree into indexed elements with bounds and center coordinates. Tap by index, resource ID, or text — no coordinate guessing needed.',
    icon: 'grid',
  },
  {
    title: 'Parallel Observation',
    description:
      'Screenshot and UI tree captured concurrently. Partial failure tolerance — if one fails, the other still returns. One command to see everything.',
    icon: 'eye',
  },
  {
    title: 'No APK Required',
    description:
      'Pure ADB-based control. No app installation, no accessibility services, no permissions dialogs. Just connect a device and start automating.',
    icon: 'zap',
  },
  {
    title: 'Single Go Binary',
    description:
      'Distributed as a single compiled binary. No Python, no Node.js, no Java runtime. Built with Go 1.24. Download, connect, control.',
    icon: 'package',
  },
  {
    title: 'Agent-First Design',
    description:
      'The skill command outputs machine-readable capability descriptions for AI agents. Structured JSON output designed for LLM consumption. Built for machines, usable by humans.',
    icon: 'bot',
  },
]

export const cliExamples = [
  {
    title: 'Observe & Inspect',
    commands: [
      { cmd: 'adbclaw observe', comment: 'Screenshot + UI tree' },
      { cmd: 'adbclaw screenshot --width 720', comment: 'Downscaled capture' },
      { cmd: 'adbclaw ui tree', comment: 'Interactive elements' },
      { cmd: 'adbclaw ui find --text "Login"', comment: 'Find by text' },
    ],
  },
  {
    title: 'Input & Control',
    commands: [
      { cmd: 'adbclaw tap --index 5', comment: 'Tap by element index' },
      { cmd: 'adbclaw tap --text "Login"', comment: 'Tap by text' },
      { cmd: 'adbclaw swipe 540 1800 540 600', comment: 'Scroll down' },
      { cmd: 'adbclaw key BACK', comment: 'Press back key' },
      { cmd: 'adbclaw type "hello world"', comment: 'Type text' },
    ],
  },
  {
    title: 'Device & Apps',
    commands: [
      { cmd: 'adbclaw device list', comment: 'Connected devices' },
      { cmd: 'adbclaw device info', comment: 'Model, screen, version' },
      { cmd: 'adbclaw app current', comment: 'Foreground app' },
      { cmd: 'adbclaw app launch com.example', comment: 'Launch app' },
      { cmd: 'adbclaw doctor', comment: 'Environment check' },
    ],
  },
]

export const architectureSteps = [
  {
    label: 'AI Agent',
    sublabel: 'Claude / LLM',
    description: 'Sends structured commands via JSON stdout or reads skill descriptions',
  },
  {
    label: 'adbclaw',
    sublabel: 'Go CLI',
    description: 'Translates commands to ADB operations with structured JSON responses',
  },
  {
    label: 'ADB',
    sublabel: 'USB / WiFi',
    description: 'Transports shell commands to the connected Android device',
  },
  {
    label: 'Device',
    sublabel: 'Android',
    description: 'Executes operations — screenshots, UI dumps, input events, app management',
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
  { step: '01', action: 'Observe', detail: 'Run observe to get screenshot + UI tree' },
  { step: '02', action: 'Decide', detail: 'AI agent analyzes screen state and picks action' },
  { step: '03', action: 'Act', detail: 'Execute tap/swipe/type using element index' },
  { step: '04', action: 'Verify', detail: 'Re-observe to confirm the action succeeded' },
]
