import App from './App.vue'
import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'

import { createPinia } from 'pinia'

import { THEME_KEY } from './constants/auth'
import './index.css'
import { router } from './router'
import { useAuthStore } from './stores/auth'

const savedTheme = localStorage.getItem(THEME_KEY) || 'dark'
if (
  savedTheme === 'dark' ||
  (savedTheme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
) {
  document.documentElement.classList.add('dark')
} else {
  document.documentElement.classList.remove('dark')
}

const app = createApp(App)
const i18n = createI18n({ legacy: false, locale: 'en-US', fallbackLocale: 'en-US' })
const pinia = createPinia()

app.use(i18n)
app.use(pinia)

const authStore = useAuthStore()

async function initializeApp() {
  if (authStore.hasToken) {
    try {
      await authStore.loadUser()
    } catch (error) {
      console.log('Session expired or invalid token')
      authStore.clearError()
    }
  }

  app.use(router)
  app.mount('#app')
}

initializeApp()
