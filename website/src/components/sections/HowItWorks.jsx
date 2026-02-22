import { architectureSteps } from '../../data/content'

export default function HowItWorks() {
  return (
    <section className="border-y border-gray-800/50 bg-surface-900/30 py-24">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mb-16 text-center">
          <h2 className="mb-4 text-3xl font-bold text-white">How It Works</h2>
          <p className="mx-auto max-w-2xl text-gray-400">
            AI agent commands flow through ADB Claw, which translates them into stealthy device operations using the real hardware input path.
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-4">
          {architectureSteps.map((step, i) => (
            <div key={step.label} className="relative flex flex-col items-center">
              {i < architectureSteps.length - 1 && (
                <div className="absolute top-10 left-[calc(50%+2rem)] right-[calc(-50%+2rem)] hidden h-px bg-gradient-to-r from-brand-500/50 to-brand-500/10 md:block" />
              )}
              <div className="mb-4 flex h-20 w-20 items-center justify-center rounded-2xl border border-gray-800 bg-surface-900 text-2xl font-bold text-brand-400">
                {i + 1}
              </div>
              <h3 className="mb-1 text-lg font-semibold text-white">{step.label}</h3>
              <span className="mb-3 text-xs font-medium text-brand-400">{step.sublabel}</span>
              <p className="text-center text-sm text-gray-400">{step.description}</p>
            </div>
          ))}
        </div>

        <div className="mt-16 rounded-xl border border-gray-800 bg-surface-900/50 p-6">
          <h3 className="mb-4 text-sm font-semibold uppercase tracking-wider text-gray-500">Input Injection Levels</h3>
          <div className="grid gap-4 md:grid-cols-3">
            {[
              {
                level: 'Level 1 — sendevent',
                stealth: 'High',
                color: 'text-green-400',
                desc: 'Real hardware input path. Genuine device ID, real touch source, custom pressure & area. No root required.',
              },
              {
                level: 'Level 2 — UHID',
                stealth: 'High',
                color: 'text-green-400',
                desc: 'scrcpy UHID mode for keyboard and mouse input with real device identity. No root required.',
              },
              {
                level: 'Level 3 — adb input',
                stealth: 'Low',
                color: 'text-yellow-400',
                desc: 'Standard adb shell input — fast but easily detectable. Fallback for simple scenarios.',
              },
            ].map((item) => (
              <div key={item.level} className="rounded-lg border border-gray-800 bg-surface-900 p-4">
                <div className="mb-2 flex items-center justify-between">
                  <span className="text-sm font-medium text-gray-200">{item.level}</span>
                  <span className={`text-xs font-medium ${item.color}`}>{item.stealth} stealth</span>
                </div>
                <p className="text-xs leading-relaxed text-gray-500">{item.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
