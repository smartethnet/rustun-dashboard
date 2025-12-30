import { createI18n } from 'vue-i18n'
import en from './en'
import zh from './zh'

// Get saved locale from localStorage or use browser language
const savedLocale = localStorage.getItem('locale')
const browserLocale = navigator.language.toLowerCase()
const defaultLocale = savedLocale || (browserLocale.startsWith('zh') ? 'zh' : 'en')

const i18n = createI18n({
  legacy: false, // Use Composition API
  locale: defaultLocale,
  fallbackLocale: 'en',
  messages: {
    en,
    zh,
  },
})

export default i18n

