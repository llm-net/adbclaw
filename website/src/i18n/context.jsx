import { createContext, useContext, useState } from 'react'
import en from './en'
import zh from './zh'

const translations = { en, zh }
const LanguageContext = createContext()

export function LanguageProvider({ children }) {
  const [lang, setLang] = useState(() => {
    try {
      return localStorage.getItem('lang') || 'en'
    } catch {
      return 'en'
    }
  })

  const switchLang = (l) => {
    setLang(l)
    try { localStorage.setItem('lang', l) } catch {}
  }

  return (
    <LanguageContext.Provider value={{ lang, setLang: switchLang, t: translations[lang] }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  return useContext(LanguageContext)
}
