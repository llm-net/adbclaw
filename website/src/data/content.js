export const features = [
  {
    title: 'Anti-Detection',
    description:
      'Uses sendevent to inject touch events through the real hardware input path. Events carry genuine device IDs and SOURCE_TOUCHSCREEN — invisible to app-level detection.',
    icon: 'shield',
  },
  {
    title: 'Human-Like Input',
    description:
      'Gaussian coordinate jitter, realistic pressure curves, natural touch area variation, and Bézier curve swipe trajectories that pass behavioral analysis.',
    icon: 'fingerprint',
  },
  {
    title: 'Observe Commands',
    description:
      'Screenshot, screen streaming, UI hierarchy dumps, element search by text or ID — everything an AI agent needs to see and understand the device.',
    icon: 'eye',
  },
  {
    title: 'No APK Required',
    description:
      'Pure ADB-based control. No app installation on the target device. No accessibility services. No permissions dialogs. Just connect and go.',
    icon: 'zap',
  },
  {
    title: 'Single Binary',
    description:
      'Distributed as a single compiled binary. No Python, no Node.js, no Java runtime. Download, connect a device, and start controlling.',
    icon: 'package',
  },
  {
    title: 'Agent-First Design',
    description:
      'JSON stdout mode for piping, MCP server for Claude and AI agents, structured outputs designed for LLM consumption. Built for machines, usable by humans.',
    icon: 'bot',
  },
]

export const cliExamples = [
  {
    title: 'Device Control',
    commands: [
      { cmd: 'adbclaw devices', comment: 'List connected devices' },
      { cmd: 'adbclaw screenshot -o screen.png', comment: 'Capture screen' },
      { cmd: 'adbclaw tap 500 800 --humanize', comment: 'Human-like tap' },
      { cmd: 'adbclaw swipe 100 500 400 500 --duration 300ms', comment: 'Swipe gesture' },
      { cmd: 'adbclaw type "hello world"', comment: 'Type text' },
    ],
  },
  {
    title: 'Stealth Mode',
    commands: [
      { cmd: 'adbclaw stealth status', comment: 'Check detection status' },
      { cmd: 'adbclaw stealth enable', comment: 'Enable anti-detection' },
      { cmd: 'adbclaw tap 500 800 --method sendevent', comment: 'Force sendevent injection' },
    ],
  },
  {
    title: 'AI Agent Integration',
    commands: [
      { cmd: 'adbclaw mcp serve', comment: 'Start MCP server for AI agents' },
      { cmd: 'adbclaw ui tree --json', comment: 'UI hierarchy as JSON' },
      { cmd: 'adbclaw app current --json', comment: 'Current app info as JSON' },
    ],
  },
]

export const architectureSteps = [
  {
    label: 'AI Agent',
    sublabel: 'Claude / LLM',
    description: 'Sends commands via MCP protocol or JSON stdout pipe',
  },
  {
    label: 'ADB Claw',
    sublabel: 'CLI Tool',
    description: 'Translates commands to stealthy device operations',
  },
  {
    label: 'ADB',
    sublabel: 'USB / WiFi',
    description: 'Transports commands to the Android device',
  },
  {
    label: 'Device',
    sublabel: 'sendevent',
    description: 'Executes via real hardware input path — undetectable',
  },
]

export const relatedProjects = [
  {
    name: 'DroidRun',
    url: 'https://github.com/droidrun/droidrun',
    description: 'LLM-powered Android device control framework with multi-model support.',
    stars: '7.7k',
  },
  {
    name: 'mobile-use',
    url: 'https://github.com/anthropics/mobile-use',
    description: 'AI agent for controlling real mobile devices, top performer on AndroidWorld benchmark.',
    stars: '2.2k',
  },
  {
    name: 'DroidClaw',
    url: 'https://github.com/unitedbyai/droidclaw',
    description: 'Natural language to ADB operations bridge, TypeScript implementation.',
    stars: '875',
  },
  {
    name: 'scrcpy',
    url: 'https://github.com/Genymobile/scrcpy',
    description: 'Display and control Android devices — the gold standard for screen mirroring.',
    stars: '136k',
  },
]

export const stealthLevels = [
  {
    level: 'Level 1',
    name: 'sendevent',
    stealth: 'High',
    description: 'Real hardware input path. Genuine device ID, real touch source, custom pressure & area.',
    rootRequired: false,
  },
  {
    level: 'Level 2',
    name: 'UHID',
    stealth: 'High',
    description: 'scrcpy UHID mode for keyboard and mouse input with real device identity.',
    rootRequired: false,
  },
  {
    level: 'Level 3',
    name: 'adb input',
    stealth: 'Low',
    description: 'Standard adb shell input — fast but easily detectable. Fallback for simple scenarios.',
    rootRequired: false,
  },
]
