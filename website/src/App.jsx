import Hero from './components/sections/Hero'
import Features from './components/sections/Features'
import HowItWorks from './components/sections/HowItWorks'
import CodeDemo from './components/sections/CodeDemo'
import RelatedProjects from './components/sections/RelatedProjects'
import Footer from './components/sections/Footer'

export default function App() {
  return (
    <div className="min-h-screen bg-gray-950">
      <Hero />
      <Features />
      <HowItWorks />
      <CodeDemo />
      <RelatedProjects />
      <Footer />
    </div>
  )
}
