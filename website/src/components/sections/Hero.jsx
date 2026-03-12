import { useState, useEffect } from 'react'

const terminalLines = [
  { cmd: 'adbclaw observe --width 720', delay: 0 },
  { out: '{"ok":true,"command":"observe","data":{...}}', delay: 800 },
  { cmd: 'adbclaw tap --text "Search"', delay: 1600 },
  { out: '{"ok":true,"command":"tap","duration_ms":38}', delay: 2200 },
  { cmd: 'adbclaw type "cat videos"', delay: 3000 },
  { out: '{"ok":true,"command":"type","duration_ms":52}', delay: 3600 },
  { cmd: 'adbclaw scroll down --pages 2', delay: 4400 },
  { out: '{"ok":true,"command":"scroll","duration_ms":680}', delay: 5000 },
  { cmd: 'adbclaw wait --text "Results" --timeout 5000', delay: 5800 },
  { out: '{"ok":true,"command":"wait","data":{"found":true}}', delay: 6600 },
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
    <div className="relative">
      <div className="absolute -inset-px rounded-xl bg-gradient-to-b from-amber-500/20 via-amber-500/5 to-transparent" />
      <div className="relative rounded-xl border border-stone-800 bg-surface-900 overflow-hidden glow-amber">
        <div className="flex items-center gap-2 px-4 py-3 border-b border-stone-800/80 bg-surface-850">
          <div className="flex gap-2">
            <span className="w-3 h-3 rounded-full bg-stone-700 hover:bg-red-500/80 transition-colors" />
            <span className="w-3 h-3 rounded-full bg-stone-700 hover:bg-yellow-500/80 transition-colors" />
            <span className="w-3 h-3 rounded-full bg-stone-700 hover:bg-green-500/80 transition-colors" />
          </div>
          <span className="text-[11px] text-stone-600 ml-2 font-mono tracking-wider uppercase">adbclaw</span>
        </div>
        <div className="p-5 font-mono text-[13px] leading-[1.8] min-h-[260px] scanline">
          {terminalLines.slice(0, visibleLines).map((line, i) => (
            <div key={i} className="flex gap-2">
              {line.cmd ? (
                <>
                  <span className="text-amber-500 select-none shrink-0">$</span>
                  <span className="text-stone-200">{line.cmd}</span>
                </>
              ) : (
                <span className="text-stone-500 ml-4">{line.out}</span>
              )}
            </div>
          ))}
          {visibleLines < terminalLines.length && (
            <div className="flex gap-2">
              <span className="text-amber-500 select-none">$</span>
              <span className="w-2.5 h-5 bg-amber-500/80 animate-blink" />
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default function Hero() {
  return (
    <section className="relative overflow-hidden">
      {/* Background effects */}
      <div className="absolute inset-0">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-amber-500/[0.03] rounded-full blur-[100px]" />
        <div className="absolute bottom-0 right-1/4 w-64 h-64 bg-amber-600/[0.02] rounded-full blur-[80px]" />
      </div>

      {/* Grid overlay */}
      <div
        className="absolute inset-0 opacity-[0.03]"
        style={{
          backgroundImage: 'linear-gradient(rgba(245,158,11,0.3) 1px, transparent 1px), linear-gradient(90deg, rgba(245,158,11,0.3) 1px, transparent 1px)',
          backgroundSize: '64px 64px',
        }}
      />

      <div className="relative mx-auto max-w-6xl px-6 pt-28 pb-24">
        {/* Nav-like top bar */}
        <div className="mb-20 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className="text-2xl">🦀</span>
            <span className="font-display text-lg font-bold tracking-tight text-stone-100">
              adbclaw
            </span>
          </div>
          <div className="flex items-center gap-6">
            <a href="#features" className="text-sm text-stone-500 hover:text-stone-300 transition-colors font-mono">
              features
            </a>
            <a href="#install" className="text-sm text-stone-500 hover:text-stone-300 transition-colors font-mono">
              install
            </a>
            <a href="#commands" className="text-sm text-stone-500 hover:text-stone-300 transition-colors font-mono">
              commands
            </a>
            <a href="#usage" className="text-sm text-stone-500 hover:text-stone-300 transition-colors font-mono">
              usage
            </a>
            <a
              href="https://github.com/llm-net/adbclaw"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 rounded-lg border border-stone-800 bg-surface-900/80 px-4 py-2 text-sm text-stone-400 hover:border-amber-500/30 hover:text-stone-200 transition-all font-mono"
            >
              <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2Z" clipRule="evenodd" />
              </svg>
              GitHub
            </a>
          </div>
        </div>

        <div className="grid gap-16 lg:grid-cols-[1.1fr_1fr] items-center">
          <div>
            <div className="mb-5 flex flex-wrap items-center gap-2.5">
              <span className="inline-flex items-center gap-2 rounded-full border border-amber-500/20 bg-amber-500/5 px-4 py-1.5">
                <span className="h-1.5 w-1.5 rounded-full bg-amber-500 animate-pulse" />
                <span className="text-xs font-mono text-amber-500/80 tracking-wide">Claude Code Plugin</span>
              </span>
              <span className="inline-flex items-center gap-2 rounded-full border border-amber-500/20 bg-amber-500/5 px-4 py-1.5">
                <span className="h-1.5 w-1.5 rounded-full bg-amber-500 animate-pulse" />
                <span className="text-xs font-mono text-amber-500/80 tracking-wide">OpenClaw Skill</span>
              </span>
            </div>

            <h1 className="mb-6 font-display text-5xl font-bold tracking-tight text-stone-50 sm:text-6xl lg:text-7xl leading-[1.05]">
              Android control<br />
              <span className="text-gradient">for AI agents</span>
            </h1>

            <p className="mb-10 max-w-lg text-lg leading-relaxed text-stone-400 font-body">
              30+ commands over ADB — observe screens, tap by element index, scroll smartly, open deep links, wait for UI state, manage apps, transfer files. Structured JSON in, structured JSON out. Available as a Claude Code plugin and OpenClaw skill.
            </p>

            <div className="flex flex-wrap gap-4">
              <a
                href="https://github.com/llm-net/adbclaw"
                target="_blank"
                rel="noopener noreferrer"
                className="group inline-flex items-center gap-2.5 rounded-lg bg-amber-500 px-6 py-3 text-sm font-semibold text-surface-950 transition-all hover:bg-amber-400 hover:shadow-lg hover:shadow-amber-500/20"
              >
                <span>Get Started</span>
                <svg className="w-4 h-4 transition-transform group-hover:translate-x-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
                </svg>
              </a>
              <a
                href="#usage"
                className="inline-flex items-center gap-2 rounded-lg border border-stone-800 px-6 py-3 text-sm font-medium text-stone-400 transition-all hover:border-stone-700 hover:text-stone-200"
              >
                See examples
              </a>
            </div>

            {/* Quick install */}
            <div className="mt-10 flex items-center gap-3 rounded-lg border border-stone-800/60 bg-surface-900/50 px-4 py-2.5 max-w-lg">
              <span className="text-amber-500/60 font-mono text-sm select-none">$</span>
              <code className="text-sm font-mono text-stone-400 truncate">curl -fsSL https://adbclaw.com/install.sh | bash</code>
              <span className="ml-auto text-[10px] text-stone-600 font-mono uppercase tracking-wider shrink-0">v1.3.0</span>
            </div>
          </div>

          <div className="animate-slide-right" style={{ animationDelay: '0.3s' }}>
            <TerminalAnimation />
          </div>
        </div>
      </div>
    </section>
  )
}
