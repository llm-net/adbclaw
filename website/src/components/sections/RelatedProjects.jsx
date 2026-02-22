import { relatedProjects } from '../../data/content'

export default function RelatedProjects() {
  return (
    <section className="border-t border-gray-800/50 bg-surface-900/30 py-24">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mb-12 text-center">
          <h2 className="mb-4 text-3xl font-bold text-white">Related Projects</h2>
          <p className="mx-auto max-w-2xl text-gray-400">
            Other great tools in the Android automation and AI agent ecosystem.
          </p>
        </div>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {relatedProjects.map((project) => (
            <a
              key={project.name}
              href={project.url}
              target="_blank"
              rel="noopener noreferrer"
              className="group rounded-xl border border-gray-800 bg-surface-900/50 p-5 transition-all hover:border-gray-700 hover:bg-surface-900"
            >
              <div className="mb-3 flex items-center justify-between">
                <h3 className="font-semibold text-gray-200 group-hover:text-white transition-colors">{project.name}</h3>
                <span className="flex items-center gap-1 text-xs text-gray-500">
                  <svg className="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 .587l3.668 7.568L24 9.306l-6.064 5.828 1.48 8.279L12 19.446l-7.417 3.967 1.481-8.279L0 9.306l8.332-1.151z" />
                  </svg>
                  {project.stars}
                </span>
              </div>
              <p className="text-sm leading-relaxed text-gray-500">{project.description}</p>
            </a>
          ))}
        </div>
      </div>
    </section>
  )
}
