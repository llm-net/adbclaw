export default function CodeBlock({ commands, title }) {
  return (
    <div className="rounded-lg border border-gray-800 bg-surface-900 overflow-hidden">
      {title && (
        <div className="flex items-center gap-2 px-4 py-2.5 border-b border-gray-800 bg-surface-850">
          <div className="flex gap-1.5">
            <span className="w-3 h-3 rounded-full bg-red-500/70" />
            <span className="w-3 h-3 rounded-full bg-yellow-500/70" />
            <span className="w-3 h-3 rounded-full bg-green-500/70" />
          </div>
          <span className="text-xs text-gray-400 ml-2 font-mono">{title}</span>
        </div>
      )}
      <div className="p-4 overflow-x-auto">
        <pre className="text-sm font-mono leading-relaxed">
          {commands.map((line, i) => (
            <div key={i} className="flex gap-2">
              <span className="text-brand-500 select-none shrink-0">$</span>
              <span className="text-gray-200">{line.cmd}</span>
              {line.comment && (
                <span className="text-gray-600 shrink-0"> # {line.comment}</span>
              )}
            </div>
          ))}
        </pre>
      </div>
    </div>
  )
}
