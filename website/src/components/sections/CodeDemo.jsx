import CodeBlock from '../ui/CodeBlock'
import { cliExamples } from '../../data/content'

export default function CodeDemo() {
  return (
    <section className="mx-auto max-w-6xl px-6 py-24">
      <div className="mb-12 text-center">
        <h2 className="mb-4 text-3xl font-bold text-white">CLI Usage</h2>
        <p className="mx-auto max-w-2xl text-gray-400">
          Clean, intuitive commands for device control, stealth management, and AI agent integration.
        </p>
      </div>
      <div className="grid gap-6 lg:grid-cols-3">
        {cliExamples.map((example) => (
          <CodeBlock key={example.title} title={example.title} commands={example.commands} />
        ))}
      </div>
    </section>
  )
}
