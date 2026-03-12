import Hero from './components/sections/Hero'
import Features from './components/sections/Features'
import Install from './components/sections/Install'
import HowItWorks from './components/sections/HowItWorks'
import CodeDemo from './components/sections/CodeDemo'
import CommandTree from './components/sections/CommandTree'
import RelatedProjects from './components/sections/RelatedProjects'
import Footer from './components/sections/Footer'

export default function App() {
  return (
    <div className="min-h-screen bg-surface-950 noise">
      <Hero />
      <Features />
      <Install />
      <HowItWorks />
      <CodeDemo />
      <CommandTree />
      <RelatedProjects />
      <Footer />
    </div>
  )
}
