const commands = [
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
]

export default function CommandTree() {
  return (
    <section id="commands" className="relative py-28">
      <div className="absolute inset-0 bg-gradient-to-b from-surface-950 via-surface-900/20 to-surface-950" />
      <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-stone-800/50 to-transparent" />

      <div className="relative mx-auto max-w-6xl px-6">
        <div className="mb-16">
          <span className="inline-block mb-4 text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
            Reference
          </span>
          <h2 className="text-3xl font-display font-bold tracking-tight text-stone-100 sm:text-4xl">
            Complete command reference
          </h2>
          <p className="mt-4 max-w-xl text-stone-500 leading-relaxed">
            Every command returns structured JSON. All commands support <code className="text-stone-400">-s</code> (serial), <code className="text-stone-400">-o</code> (output format), <code className="text-stone-400">--timeout</code>, and <code className="text-stone-400">--verbose</code>.
          </p>
        </div>

        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {commands.map((group) => (
            <div key={group.category} className="rounded-xl border border-stone-800/60 bg-surface-900/30 overflow-hidden">
              <div className="px-5 py-3 border-b border-stone-800/40">
                <h3 className="text-sm font-semibold text-stone-300 font-display">{group.category}</h3>
              </div>
              <div className="p-4 space-y-2.5">
                {group.items.map((item) => (
                  <div key={item.cmd} className="flex items-start gap-3">
                    <code className="text-[12px] font-mono text-amber-500/70 shrink-0 pt-0.5 min-w-[120px]">{item.cmd}</code>
                    <div className="min-w-0">
                      <span className="text-xs text-stone-500">{item.desc}</span>
                      {item.flags && (
                        <span className="block text-[10px] text-stone-700 font-mono mt-0.5">{item.flags}</span>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
