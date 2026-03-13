import { useLanguage } from '../../i18n/context'

export default function HowItWorks() {
  const { t } = useLanguage()

  return (
    <section className="relative py-28">
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-b from-surface-950 via-surface-900/20 to-surface-950" />
      <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-stone-800/50 to-transparent" />
      <div className="absolute bottom-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-stone-800/50 to-transparent" />

      <div className="relative mx-auto max-w-6xl px-6">
        <div className="mb-16">
          <span className="inline-block mb-4 text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
            {t.howItWorks.label}
          </span>
          <h2 className="text-3xl font-display font-bold tracking-tight text-stone-100 sm:text-4xl">
            {t.howItWorks.title}
          </h2>
          <p className="mt-4 max-w-xl text-stone-500 leading-relaxed">
            {t.howItWorks.description}
          </p>
        </div>

        {/* Architecture pipeline */}
        <div className="grid gap-3 md:grid-cols-4 mb-20">
          {t.howItWorks.architectureSteps.map((step, i) => (
            <div key={step.label} className="relative group">
              {/* Connector arrow */}
              {i < t.howItWorks.architectureSteps.length - 1 && (
                <div className="absolute top-1/2 -right-3 z-10 hidden md:block">
                  <svg className="w-6 h-6 text-amber-500/30" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
                  </svg>
                </div>
              )}

              <div className="h-full rounded-xl border border-stone-800/60 bg-surface-900/40 p-5 transition-all duration-300 hover:border-amber-500/20">
                <div className="flex items-center gap-3 mb-3">
                  <span className="flex items-center justify-center w-8 h-8 rounded-lg bg-amber-500/10 text-sm font-mono font-bold text-amber-500/80">
                    {i + 1}
                  </span>
                  <div>
                    <h3 className="text-sm font-semibold text-stone-200 font-display">{step.label}</h3>
                    <span className="text-[11px] font-mono text-stone-600">{step.sublabel}</span>
                  </div>
                </div>
                <p className="text-xs leading-relaxed text-stone-500">{step.description}</p>
              </div>
            </div>
          ))}
        </div>

        {/* Agent workflow */}
        <div className="rounded-xl border border-stone-800/60 bg-surface-900/30 overflow-hidden">
          <div className="px-6 py-4 border-b border-stone-800/40">
            <h3 className="text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
              {t.howItWorks.agentLoop}
            </h3>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4">
            {t.howItWorks.agentWorkflow.map((item, i) => (
              <div
                key={item.step}
                className={`p-6 ${i < t.howItWorks.agentWorkflow.length - 1 ? 'lg:border-r border-b lg:border-b-0 border-stone-800/30' : ''}`}
              >
                <span className="block mb-3 text-2xl font-display font-bold text-amber-500/20">{item.step}</span>
                <h4 className="mb-1.5 text-sm font-semibold text-stone-200 font-display">{item.action}</h4>
                <p className="text-xs leading-relaxed text-stone-500">{item.detail}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
