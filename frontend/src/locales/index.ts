import { createI18n } from 'vue-i18n'

// 引入中文语言包
import zhCN from './zh-CN'
// 引入英文语言包
import enUS from './en-US'

const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: localStorage.getItem('locale') || 'zh-CN', // 默认中文
  fallbackLocale: 'en-US', // 回退语言
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  }
})

export default i18n