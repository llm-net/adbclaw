import { useState, useEffect } from 'react'

const terminalLines = [
  { cmd: 'adbclaw stealth enable', delay: 0 },
  { out: '✓ Wireless ADB connected', delay: 800 },
  { out: '✓ sendevent input mode active', delay: 1200 },
  { out: '✓ Stealth mode enabled', delay: 1600 },
  { cmd: 'adbclaw tap 540 1200 --humanize', delay: 2400 },
  { out: '✓ Tap injected via /dev/input/event2 (pressure: 62, area: 5)', delay: 3200 },
]

function TerminalAnimation() {
  const [visibleLines, setVisibleLines] = useState(0)

  useEffect(() => {
    const timers = terminalLines.map((line, i) =>
      setTimeout(() => setVisibleLines(i + 1), line.delay)
    )
    return () => timers.forEach(clearTimeout)
  }, [])

  return (
    <div className="rounded-lg border border-gray-800 bg-surface-900 overflow-hidden shadow-2xl shadow-brand-500/5">
      <div className="flex items-center gap-2 px-4 py-2.5 border-b border-gray-800 bg-surface-850">
        <div className="flex gap-1.5">
          <span className="w-3 h-3 rounded-full bg-red-500/70" />
          <span className="w-3 h-3 rounded-full bg-yellow-500/70" />
          <span className="w-3 h-3 rounded-full bg-green-500/70" />
        </div>
        <span className="text-xs text-gray-500 ml-2 font-mono">terminal</span>
      </div>
      <div className="p-4 font-mono text-sm leading-relaxed min-h-[200px]">
        {terminalLines.slice(0, visibleLines).map((line, i) => (
          <div key={i} className="flex gap-2">
            {line.cmd ? (
              <>
                <span className="text-brand-500 select-none">$</span>
                <span className="text-gray-200">{line.cmd}</span>
              </>
            ) : (
              <span className="text-gray-400 ml-4">{line.out}</span>
            )}
          </div>
        ))}
        {visibleLines < terminalLines.length && (
          <div className="flex gap-2 mt-0">
            <span className="text-brand-500 select-none">$</span>
            <span className="w-2 h-5 bg-brand-400 animate-pulse" />
          </div>
        )}
      </div>
    </div>
  )
}

export default function Hero() {
  return (
    <section className="relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--color-brand-500)_0%,_transparent_50%)] opacity-[0.07]" />
      <div className="relative mx-auto max-w-6xl px-6 pt-24 pb-20">
        <div className="grid gap-12 lg:grid-cols-2 lg:gap-16 items-center">
          <div>
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-gray-800 bg-surface-900/80 px-4 py-1.5 text-sm text-gray-400">
              <span className="h-2 w-2 rounded-full bg-brand-500 animate-pulse" />
              Under active development
            </div>
            <h1 className="mb-6 text-4xl font-bold tracking-tight text-white sm:text-5xl lg:text-6xl">
              <span className="text-brand-400">ADB</span> Claw
            </h1>
            <p className="mb-8 max-w-lg text-lg leading-relaxed text-gray-400">
              Stealthy Android device control CLI for AI agents. Inject touch events through real hardware input paths — invisible to app-level detection.
            </p>
            <div className="flex flex-wrap gap-4">
              <a
                href="https://github.com/llm-net/adbclaw"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 rounded-lg bg-brand-600 px-5 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-500"
              >
                <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                  <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2Z" clipRule="evenodd" />
                </svg>
                View on GitHub
              </a>
              <a
                href="#features"
                className="inline-flex items-center gap-2 rounded-lg border border-gray-700 px-5 py-2.5 text-sm font-medium text-gray-300 transition-colors hover:border-gray-600 hover:text-white"
              >
                Learn more
              </a>
            </div>
          </div>
          <div>
            <TerminalAnimation />
          </div>
        </div>
      </div>
    </section>
  )
}
