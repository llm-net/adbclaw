import FeatureCard from '../ui/FeatureCard'
import { features } from '../../data/content'

export default function Features() {
  return (
    <section id="features" className="mx-auto max-w-6xl px-6 py-24">
      <div className="mb-12 text-center">
        <h2 className="mb-4 text-3xl font-bold text-white">Core Features</h2>
        <p className="mx-auto max-w-2xl text-gray-400">
          Built from the ground up for stealth and automation. Every command is designed to be undetectable by the target app.
        </p>
      </div>
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {features.map((feature) => (
          <FeatureCard key={feature.title} {...feature} />
        ))}
      </div>
    </section>
  )
}
