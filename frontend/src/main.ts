import App from './App.vue'
import { createApp } from 'vue'

import { createPinia } from 'pinia'

import './index.css'
import { router } from './router'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

const savedTheme = localStorage.getItem('flow-theme') || 'dark'
if (
  savedTheme === 'dark' ||
  (savedTheme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
) {
  document.documentElement.classList.add('dark')
} else {
  document.documentElement.classList.remove('dark')
}

app.mount('#app')

const authStore = useAuthStore()
if (authStore.hasToken) {
  authStore.loadUser().catch(() => {
    console.log('Session expired or invalid token')
  })
}
