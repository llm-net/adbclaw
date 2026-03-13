import CodeBlock from '../ui/CodeBlock'
import { useLanguage } from '../../i18n/context'

export default function CodeDemo() {
  const { t } = useLanguage()

  return (
    <section id="usage" className="relative py-28">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-48 h-px bg-gradient-to-r from-transparent via-amber-500/20 to-transparent" />

      <div className="mx-auto max-w-6xl px-6">
        <div className="mb-16">
          <span className="inline-block mb-4 text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
            {t.codeDemo.label}
          </span>
          <h2 className="text-3xl font-display font-bold tracking-tight text-stone-100 sm:text-4xl">
            {t.codeDemo.title}
          </h2>
          <p className="mt-4 max-w-xl text-stone-500 leading-relaxed">
            {t.codeDemo.description}
          </p>
        </div>
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
          {t.codeDemo.examples.map((example) => (
            <CodeBlock key={example.title} title={example.title} commands={example.commands} />
          ))}
        </div>

        {/* JSON output example */}
        <div className="mt-8 rounded-xl border border-stone-800/60 bg-surface-900/40 overflow-hidden">
          <div className="flex items-center justify-between px-5 py-3 border-b border-stone-800/40">
            <span className="text-xs text-stone-400 font-mono tracking-wide">{t.codeDemo.jsonEnvelope}</span>
            <span className="text-[10px] text-stone-600 font-mono uppercase tracking-wider">{t.codeDemo.everyCommand}</span>
          </div>
          <div className="p-5 overflow-x-auto scanline">
            <pre className="text-[13px] font-mono leading-[1.9]">
              <span className="text-stone-600">{'{'}</span>{'\n'}
              <span className="text-stone-500">{'  '}"ok"</span><span className="text-stone-600">: </span><span className="text-amber-500/80">true</span><span className="text-stone-600">,</span>{'\n'}
              <span className="text-stone-500">{'  '}"command"</span><span className="text-stone-600">: </span><span className="text-amber-600/60">"tap"</span><span className="text-stone-600">,</span>{'\n'}
              <span className="text-stone-500">{'  '}"data"</span><span className="text-stone-600">: {'{'} </span><span className="text-stone-500">"x"</span><span className="text-stone-600">: </span><span className="text-amber-500/80">540</span><span className="text-stone-600">, </span><span className="text-stone-500">"y"</span><span className="text-stone-600">: </span><span className="text-amber-500/80">960</span><span className="text-stone-600">{' }'}</span><span className="text-stone-600">,</span>{'\n'}
              <span className="text-stone-500">{'  '}"duration_ms"</span><span className="text-stone-600">: </span><span className="text-amber-500/80">38</span><span className="text-stone-600">,</span>{'\n'}
              <span className="text-stone-500">{'  '}"timestamp"</span><span className="text-stone-600">: </span><span className="text-amber-600/60">"2026-03-12T10:00:00.123Z"</span>{'\n'}
              <span className="text-stone-600">{'}'}</span>
            </pre>
          </div>
        </div>
      </div>
    </section>
  )
}
