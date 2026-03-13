export default {
  nav: {
    features: 'features',
    install: 'install',
    commands: 'commands',
    usage: 'usage',
  },
  hero: {
    title: 'Eyes, hands & ears',
    titleHighlight: 'on Android',
    description:
      'Built for agents, claws, bots, and LLMs. 30+ structured commands — observe screens, tap by element index, read live chat during video playback, capture system audio for speech recognition, manage apps. JSON in, JSON out. Your bridge to the physical world.',
    getStarted: 'Install Skill',
    seeExamples: 'See examples',
    versionNote: 'audio capture + monitor + actively shipping',
  },
  features: {
    label: 'Capabilities',
    title: 'Superpowers for agents',
    description:
      'Pure tool layer. No LLM logic. No agent framework. Just reliable, structured commands that give you sensory and motor capabilities on Android devices.',
    items: [
      {
        title: 'Live Stream Intelligence',
        description:
          'monitor connects to Android\'s accessibility framework, reading all UI text in real-time — even during video playback where uiautomator dump hangs. Chat messages, captions, dynamic overlays. See what no other tool can expose to you.',
        icon: 'zap',
      },
      {
        title: 'System Audio Capture',
        description:
          'audio capture records device audio via REMOTE_SUBMIX (Android 11+), streaming WAV to stdout. Pipe to ASR for speech-to-text. Combined with monitor, you get full sensory coverage — visual text AND audio from any live stream.',
        icon: 'layers',
      },
      {
        title: 'Structured JSON API',
        description:
          'Every command returns {ok, command, data, error, duration_ms, timestamp}. Parse with confidence. Errors include codes and actionable suggestions. No guesswork, no regex scraping.',
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
          'Jump directly to any app screen via URI. Open web pages, trigger WeChat scan, search Taobao — skip multi-step navigation entirely. One command, instant arrival. The key to CJK text input.',
        icon: 'link',
      },
      {
        title: 'Wait for State Changes',
        description:
          'Block until a UI element appears or disappears, or an Activity loads. No polling loops in your agent code. Configurable timeout and interval. Returns the matched element on success.',
        icon: 'clock',
      },
      {
        title: 'Full Device Control',
        description:
          '30+ commands covering screen observation, input injection, smart scrolling, app lifecycle, screen management, shell access, and file transfer. Everything over standard ADB. No APK needed.',
        icon: 'package',
      },
      {
        title: 'App Knowledge Profiles',
        description:
          'Pre-built profiles for popular apps (Douyin, Meituan, etc.) with deep links, UI layouts, and known issues. Load once, skip trial-and-error. New profiles ship with every release.',
        icon: 'book',
      },
      {
        title: 'Actively Evolving',
        description:
          'New capabilities ship regularly. monitor, audio capture, App Profiles — each release expands what you can perceive and control. Install once, gain new abilities as they land. Built by agents, for agents.',
        icon: 'bot',
      },
    ],
  },
  install: {
    label: 'Install',
    title: 'Get started in seconds',
    description: 'Pre-built binaries for macOS and Linux. No Go toolchain required. Just download and run.',
    recommended: 'Recommended',
    oneLiner: 'One-liner Install',
    oneLinerDesc: 'Auto-detects your OS and architecture. Downloads the latest binary to',
    manual: 'Manual',
    downloadBinary: 'Download Binary',
    downloadBinaryDesc: 'Grab the pre-built binary for your platform from GitHub Releases.',
    fromSource: 'From source',
    buildWithGo: 'Build with Go',
    buildWithGoDesc: 'Clone the repo and build. Requires Go 1.24+.',
    prerequisite: 'Prerequisite:',
    prerequisiteText: 'ADB (Android Debug Bridge) must be installed and in your PATH.',
  },
  howItWorks: {
    label: 'Architecture',
    title: 'How it works',
    description: 'Commands flow from AI agent through adb-claw to the device. Every response is structured JSON.',
    agentLoop: 'Recommended Agent Loop',
    architectureSteps: [
      {
        label: 'AI Agent',
        sublabel: 'Claude / OpenClaw / LLM',
        description: 'Reads skill description, sends structured commands, parses JSON responses to decide next actions',
      },
      {
        label: 'adb-claw',
        sublabel: 'Go CLI · v1.5.3',
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
    ],
    agentWorkflow: [
      { step: '01', action: 'Observe', detail: 'Run observe to capture screenshot + UI tree in one call' },
      { step: '02', action: 'Decide', detail: 'AI agent analyzes screen state and plans the next action' },
      { step: '03', action: 'Act', detail: 'Tap, scroll, open deep link, or type — by element index' },
      { step: '04', action: 'Wait', detail: 'Use wait to block until UI state changes, then re-observe' },
    ],
  },
  codeDemo: {
    label: 'Usage',
    title: '30+ commands, one binary',
    description:
      'Observe, navigate, wait, and manage — all as top-level commands with structured JSON output. Prefer element index over coordinates.',
    jsonEnvelope: 'JSON Envelope',
    everyCommand: 'Every command returns this',
    examples: [
      {
        title: 'Observe & Inspect',
        commands: [
          { cmd: 'adb-claw observe', comment: 'Screenshot + UI tree' },
          { cmd: 'adb-claw screenshot --width 720', comment: 'Downscaled capture' },
          { cmd: 'adb-claw ui tree', comment: 'Indexed elements' },
          { cmd: 'adb-claw ui find --text "Login"', comment: 'Find by text' },
        ],
      },
      {
        title: 'Input & Navigate',
        commands: [
          { cmd: 'adb-claw tap --index 5', comment: 'Tap by element index' },
          { cmd: 'adb-claw type "hello world"', comment: 'Input text' },
          { cmd: 'adb-claw scroll down --pages 3', comment: 'Smart scroll' },
          { cmd: 'adb-claw open "weixin://dl/scan"', comment: 'Deep link' },
          { cmd: 'adb-claw clear-field --index 2', comment: 'Clear input' },
        ],
      },
      {
        title: 'Wait & Screen',
        commands: [
          { cmd: 'adb-claw wait --text "Done"', comment: 'Wait for element' },
          { cmd: 'adb-claw wait --text "Loading" --gone', comment: 'Wait until gone' },
          { cmd: 'adb-claw screen status', comment: 'On/off/lock/rotation' },
          { cmd: 'adb-claw screen unlock', comment: 'Wake + swipe unlock' },
        ],
      },
      {
        title: 'Monitor & Audio',
        commands: [
          { cmd: 'adb-claw monitor --stream', comment: 'Live UI text (video-safe)' },
          { cmd: 'adb-claw monitor --duration 30000', comment: '30s bounded capture' },
          { cmd: 'adb-claw audio capture --file out.wav', comment: 'Record system audio' },
          { cmd: 'adb-claw audio capture --stream | asrclaw transcribe', comment: 'Pipe to ASR' },
        ],
      },
      {
        title: 'Apps & System',
        commands: [
          { cmd: 'adb-claw app launch com.example', comment: 'Launch app' },
          { cmd: 'adb-claw app install ./app.apk', comment: 'Install APK' },
          { cmd: 'adb-claw shell "pm list packages"', comment: 'Raw shell' },
          { cmd: 'adb-claw file pull /sdcard/log.txt .', comment: 'Pull file' },
        ],
      },
    ],
  },
  commandTree: {
    label: 'Reference',
    title: 'Complete command reference',
    description: 'Every command returns structured JSON. All commands support',
    commands: [
      {
        category: 'Observation',
        items: [
          { cmd: 'observe', desc: 'Screenshot + UI tree in parallel', flags: '--width' },
          { cmd: 'screenshot', desc: 'Capture screen (base64 or file)', flags: '--file, --width' },
          { cmd: 'ui tree', desc: 'Indexed UI element tree' },
          { cmd: 'ui find', desc: 'Find elements by text/id/index', flags: '--text, --id, --index' },
        ],
      },
      {
        category: 'Input',
        items: [
          { cmd: 'tap', desc: 'Tap by coordinates or element', flags: '--index, --id, --text' },
          { cmd: 'long-press', desc: 'Long press with duration', flags: '--duration' },
          { cmd: 'swipe', desc: 'Swipe between coordinates', flags: '--duration' },
          { cmd: 'key', desc: 'Press key (30+ aliases)', flags: 'HOME, BACK, ENTER...' },
          { cmd: 'type', desc: 'Input ASCII text' },
          { cmd: 'clear-field', desc: 'Clear focused input', flags: '--index, --id, --text' },
        ],
      },
      {
        category: 'Navigation',
        items: [
          { cmd: 'scroll', desc: 'Smart scroll in any direction', flags: '--pages, --distance, --index' },
          { cmd: 'open', desc: 'Open URI / deep link' },
        ],
      },
      {
        category: 'State',
        items: [
          { cmd: 'wait', desc: 'Wait for UI element or Activity', flags: '--text, --id, --gone, --timeout' },
          { cmd: 'screen status', desc: 'Display on/off, lock, rotation' },
          { cmd: 'screen on/off', desc: 'Wake or sleep screen' },
          { cmd: 'screen unlock', desc: 'Wake + swipe unlock' },
          { cmd: 'screen rotation', desc: 'Set rotation mode', flags: 'auto, 0-3' },
        ],
      },
      {
        category: 'Apps',
        items: [
          { cmd: 'app list', desc: 'Installed apps', flags: '--all' },
          { cmd: 'app current', desc: 'Foreground package/activity' },
          { cmd: 'app launch', desc: 'Start an app by package' },
          { cmd: 'app stop', desc: 'Force stop an app' },
          { cmd: 'app install', desc: 'Install APK', flags: '--replace' },
          { cmd: 'app uninstall', desc: 'Remove app' },
          { cmd: 'app clear', desc: 'Clear app data' },
        ],
      },
      {
        category: 'Sensing',
        items: [
          { cmd: 'monitor', desc: 'Live UI text via accessibility (video-safe)', flags: '--duration, --interval, --stream' },
          { cmd: 'audio capture', desc: 'System audio → WAV stream (Android 11+)', flags: '--file, --duration, --rate, --stream' },
        ],
      },
      {
        category: 'System',
        items: [
          { cmd: 'device list', desc: 'Connected devices' },
          { cmd: 'device info', desc: 'Model, screen, SDK version' },
          { cmd: 'shell', desc: 'Execute raw ADB shell command' },
          { cmd: 'file push', desc: 'Send file to device' },
          { cmd: 'file pull', desc: 'Retrieve file from device' },
          { cmd: 'doctor', desc: 'Environment health check' },
          { cmd: 'skill', desc: 'Output skill.json for agents' },
        ],
      },
    ],
  },
  relatedProjects: {
    label: 'Ecosystem',
    title: 'Related projects',
    description: 'Other tools in the Android automation and AI agent space.',
    skillPlatform: 'Skill Platform',
    items: [
      {
        name: 'OpenClaw',
        url: 'https://github.com/openclaw/openclaw',
        description: 'Local-first personal AI assistant platform. adb-claw is published as an OpenClaw Skill on ClawHub.',
        stars: '',
        highlight: true,
      },
      {
        name: 'mobile-use',
        url: 'https://github.com/anthropics/mobile-use',
        description: "Anthropic's AI agent for controlling real mobile devices, top performer on AndroidWorld benchmark.",
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
    ],
  },
  footer: {
    description:
      'Android device control CLI for AI agent automation. 30+ commands over ADB. Pure tool layer — no LLM/Agent logic included.',
    project: 'Project',
    documentation: 'Documentation',
    issues: 'Issues',
    releases: 'Releases',
    availableOn: 'Available on',
    claudeCodePlugin: 'Claude Code Plugin',
    openClawClawHub: 'OpenClaw / ClawHub',
    standaloneCli: 'Standalone CLI',
    stack: 'Stack',
  },
}
