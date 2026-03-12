import { relatedProjects } from '../../data/content'

export default function RelatedProjects() {
  return (
    <section className="relative py-28">
      <div className="absolute inset-0 bg-gradient-to-b from-surface-950 via-surface-900/20 to-surface-950" />
      <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-stone-800/50 to-transparent" />

      <div className="relative mx-auto max-w-6xl px-6">
        <div className="mb-16">
          <span className="inline-block mb-4 text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
            Ecosystem
          </span>
          <h2 className="text-3xl font-display font-bold tracking-tight text-stone-100 sm:text-4xl">
            Related projects
          </h2>
          <p className="mt-4 max-w-xl text-stone-500 leading-relaxed">
            Other tools in the Android automation and AI agent space.
          </p>
        </div>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {relatedProjects.map((project) => (
            <a
              key={project.name}
              href={project.url}
              target="_blank"
              rel="noopener noreferrer"
              className={`group relative rounded-xl border p-5 transition-all duration-300 hover:bg-surface-900/50 ${project.highlight ? 'border-amber-500/30 bg-surface-900/40 hover:border-amber-500/40' : 'border-stone-800/60 bg-surface-900/30 hover:border-amber-500/20'}`}
            >
              <div className="mb-3 flex items-center justify-between">
                <h3 className="text-sm font-semibold text-stone-300 group-hover:text-stone-100 transition-colors font-display">
                  {project.name}
                </h3>
                {project.stars ? (
                  <span className="flex items-center gap-1 text-[11px] text-stone-600 font-mono">
                    <svg className="h-3 w-3 text-amber-500/40" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 .587l3.668 7.568L24 9.306l-6.064 5.828 1.48 8.279L12 19.446l-7.417 3.967 1.481-8.279L0 9.306l8.332-1.151z" />
                    </svg>
                    {project.stars}
                  </span>
                ) : project.highlight ? (
                  <span className="text-[10px] font-mono text-amber-500/60 uppercase tracking-wider">Skill Platform</span>
                ) : null}
              </div>
              <p className="text-xs leading-relaxed text-stone-600 group-hover:text-stone-500 transition-colors">{project.description}</p>

              {/* Arrow indicator */}
              <div className="absolute top-5 right-5 opacity-0 group-hover:opacity-100 transition-opacity">
                <svg className="w-3.5 h-3.5 text-amber-500/40" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
                  <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 19.5 15-15m0 0H8.25m11.25 0v11.25" />
                </svg>
              </div>
            </a>
          ))}
        </div>
      </div>
    </section>
  )
}
