import { useLanguage } from '../../i18n/context'

export default function Install() {
  const { t } = useLanguage()

  return (
    <section id="install" className="relative py-28">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-48 h-px bg-gradient-to-r from-transparent via-amber-500/20 to-transparent" />

      <div className="mx-auto max-w-6xl px-6">
        <div className="mb-16">
          <span className="inline-block mb-4 text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
            {t.install.label}
          </span>
          <h2 className="text-3xl font-display font-bold tracking-tight text-stone-100 sm:text-4xl">
            {t.install.title}
          </h2>
          <p className="mt-4 max-w-xl text-stone-500 leading-relaxed">
            {t.install.description}
          </p>
        </div>

        <div className="grid gap-5 lg:grid-cols-3">
          {/* Option A: curl install */}
          <div className="group relative rounded-xl border border-amber-500/30 bg-surface-900/40 p-6 transition-all duration-300 hover:border-amber-500/40">
            <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-amber-500/[0.03] to-transparent" />
            <div className="relative">
              <span className="inline-block mb-3 text-[10px] font-mono text-amber-500/80 uppercase tracking-wider">{t.install.recommended}</span>
              <h3 className="mb-2 text-base font-semibold text-stone-100 font-display">{t.install.oneLiner}</h3>
              <p className="mb-5 text-sm text-stone-500 leading-relaxed">
                {t.install.oneLinerDesc} <code className="text-stone-400">~/.local/bin</code>.
              </p>
              <div className="rounded-lg border border-stone-800/60 bg-surface-950/80 p-3 overflow-x-auto">
                <code className="text-[13px] font-mono text-stone-300 whitespace-nowrap">
                  <span className="text-amber-500/60 select-none">$ </span>
                  curl -fsSL https://adbclaw.com/install.sh | bash
                </code>
              </div>
            </div>
          </div>

          {/* Option B: Direct download */}
          <div className="group relative rounded-xl border border-stone-800/60 bg-surface-900/30 p-6 transition-all duration-300 hover:border-amber-500/20">
            <div className="relative">
              <span className="inline-block mb-3 text-[10px] font-mono text-stone-600 uppercase tracking-wider">{t.install.manual}</span>
              <h3 className="mb-2 text-base font-semibold text-stone-100 font-display">{t.install.downloadBinary}</h3>
              <p className="mb-5 text-sm text-stone-500 leading-relaxed">
                {t.install.downloadBinaryDesc}
              </p>
              <div className="space-y-2">
                {['darwin-arm64', 'darwin-amd64', 'linux-arm64', 'linux-amd64'].map((platform) => (
                  <a
                    key={platform}
                    href={`https://github.com/llm-net/adbclaw/releases/latest/download/adbclaw-${platform}`}
                    className="flex items-center gap-2 text-sm font-mono text-stone-500 hover:text-amber-500/80 transition-colors"
                  >
                    <svg className="w-3.5 h-3.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="1.5">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
                    </svg>
                    adbclaw-{platform}
                  </a>
                ))}
              </div>
            </div>
          </div>

          {/* Option C: Build from source */}
          <div className="group relative rounded-xl border border-stone-800/60 bg-surface-900/30 p-6 transition-all duration-300 hover:border-amber-500/20">
            <div className="relative">
              <span className="inline-block mb-3 text-[10px] font-mono text-stone-600 uppercase tracking-wider">{t.install.fromSource}</span>
              <h3 className="mb-2 text-base font-semibold text-stone-100 font-display">{t.install.buildWithGo}</h3>
              <p className="mb-5 text-sm text-stone-500 leading-relaxed">
                {t.install.buildWithGoDesc}
              </p>
              <div className="rounded-lg border border-stone-800/60 bg-surface-950/80 p-3 space-y-1">
                <div><code className="text-[13px] font-mono text-stone-400"><span className="text-amber-500/60 select-none">$ </span>git clone https://github.com/llm-net/adbclaw</code></div>
                <div><code className="text-[13px] font-mono text-stone-400"><span className="text-amber-500/60 select-none">$ </span>cd adbclaw/src && make build</code></div>
              </div>
            </div>
          </div>
        </div>

        {/* Prerequisites note */}
        <div className="mt-8 rounded-lg border border-stone-800/40 bg-surface-900/20 px-5 py-4">
          <p className="text-xs text-stone-600 leading-relaxed">
            <span className="text-stone-500 font-semibold">{t.install.prerequisite}</span> {t.install.prerequisiteText}
            {' '}macOS: <code className="text-stone-500">brew install android-platform-tools</code>.
            {' '}Linux: <code className="text-stone-500">apt install adb</code> / <code className="text-stone-500">pacman -S android-tools</code>.
          </p>
        </div>
      </div>
    </section>
  )
}
