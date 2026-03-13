import { useLanguage } from '../../i18n/context'

export default function Footer() {
  const { t } = useLanguage()

  return (
    <footer className="relative border-t border-stone-800/30 py-16">
      <div className="mx-auto max-w-6xl px-6">
        <div className="flex flex-col gap-10 md:flex-row md:items-start md:justify-between">
          {/* Brand */}
          <div>
            <div className="flex items-center gap-3 mb-3">
              <span className="text-xl">🦀</span>
              <span className="font-display text-lg font-bold tracking-tight text-stone-200">
                adbclaw
              </span>
            </div>
            <p className="max-w-xs text-sm text-stone-600 leading-relaxed">
              {t.footer.description}
            </p>
            <span className="inline-block mt-3 text-[10px] font-mono text-stone-700 tracking-wider">v1.4.1 · darwin-arm64 · darwin-amd64 · linux-arm64 · linux-amd64</span>
          </div>

          {/* Links */}
          <div className="flex gap-16">
            <div>
              <h4 className="text-[11px] font-mono uppercase tracking-[0.2em] text-stone-600 mb-4">{t.footer.project}</h4>
              <ul className="space-y-2.5">
                <li>
                  <a href="https://github.com/llm-net/adbclaw" target="_blank" rel="noopener noreferrer" className="text-sm text-stone-500 hover:text-stone-300 transition-colors">
                    GitHub
                  </a>
                </li>
                <li>
                  <a href="https://github.com/llm-net/adbclaw/releases" target="_blank" rel="noopener noreferrer" className="text-sm text-stone-500 hover:text-stone-300 transition-colors">
                    {t.footer.releases}
                  </a>
                </li>
                <li>
                  <a href="https://github.com/llm-net/adbclaw/tree/main/docs" target="_blank" rel="noopener noreferrer" className="text-sm text-stone-500 hover:text-stone-300 transition-colors">
                    {t.footer.documentation}
                  </a>
                </li>
                <li>
                  <a href="https://github.com/llm-net/adbclaw/issues" target="_blank" rel="noopener noreferrer" className="text-sm text-stone-500 hover:text-stone-300 transition-colors">
                    {t.footer.issues}
                  </a>
                </li>
              </ul>
            </div>
            <div>
              <h4 className="text-[11px] font-mono uppercase tracking-[0.2em] text-stone-600 mb-4">{t.footer.availableOn}</h4>
              <ul className="space-y-2.5">
                <li><span className="text-sm text-stone-500">{t.footer.claudeCodePlugin}</span></li>
                <li><span className="text-sm text-stone-500">{t.footer.openClawClawHub}</span></li>
                <li><span className="text-sm text-stone-500">{t.footer.standaloneCli}</span></li>
              </ul>
            </div>
            <div>
              <h4 className="text-[11px] font-mono uppercase tracking-[0.2em] text-stone-600 mb-4">{t.footer.stack}</h4>
              <ul className="space-y-2.5">
                <li><span className="text-sm text-stone-600">Go 1.24</span></li>
                <li><span className="text-sm text-stone-600">Cobra CLI</span></li>
                <li><span className="text-sm text-stone-600">ADB Shell</span></li>
              </ul>
            </div>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="mt-12 pt-6 border-t border-stone-800/20 flex flex-col sm:flex-row items-center justify-between gap-3">
          <span className="text-xs text-stone-700 font-mono">MIT License</span>
          <div className="flex items-center gap-4">
            <a
              href="https://adbclaw.com"
              className="text-xs text-stone-700 hover:text-stone-500 transition-colors font-mono"
            >
              adbclaw.com
            </a>
            <span className="text-stone-800">·</span>
            <a
              href="https://github.com/llm-net"
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs text-stone-700 hover:text-stone-500 transition-colors font-mono"
            >
              LLM.net
            </a>
          </div>
        </div>
      </div>
    </footer>
  )
}
