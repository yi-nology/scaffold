import { createI18n } from 'vue-i18n'
import enUS from './en-US'
import zhCN from './zh-CN'

// 获取浏览器语言
function getBrowserLocale(): string {
  const lang = navigator.language || (navigator as any).userLanguage
  if (lang.startsWith('zh')) {
    return 'zh-CN'
  }
  return 'en-US'
}

// 获取存储的语言或浏览器语言
function getDefaultLocale(): string {
  const stored = localStorage.getItem('scaffold-locale')
  if (stored && ['zh-CN', 'en-US'].includes(stored)) {
    return stored
  }
  return getBrowserLocale()
}

const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: 'en-US',
  messages: {
    'en-US': enUS,
    'zh-CN': zhCN
  }
})

export default i18n

// 切换语言并持久化
export function setLocale(locale: string) {
  if (i18n.global.locale) {
    (i18n.global.locale as any).value = locale
  }
  localStorage.setItem('scaffold-locale', locale)
}

export function getLocale(): string {
  return (i18n.global.locale as any).value || 'en-US'
}
