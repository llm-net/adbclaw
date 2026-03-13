import { useLanguage } from '../../i18n/context'

export default function CommandTree() {
  const { t } = useLanguage()

  return (
    <section id="commands" className="relative py-28">
      <div className="absolute inset-0 bg-gradient-to-b from-surface-950 via-surface-900/20 to-surface-950" />
      <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-stone-800/50 to-transparent" />

      <div className="relative mx-auto max-w-6xl px-6">
        <div className="mb-16">
          <span className="inline-block mb-4 text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
            {t.commandTree.label}
          </span>
          <h2 className="text-3xl font-display font-bold tracking-tight text-stone-100 sm:text-4xl">
            {t.commandTree.title}
          </h2>
          <p className="mt-4 max-w-xl text-stone-500 leading-relaxed">
            {t.commandTree.description} <code className="text-stone-400">-s</code> (serial), <code className="text-stone-400">-o</code> (output format), <code className="text-stone-400">--timeout</code>, and <code className="text-stone-400">--verbose</code>.
          </p>
        </div>

        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {t.commandTree.commands.map((group) => (
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
