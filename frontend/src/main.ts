import App from './App.vue'
import { createApp } from 'vue'

import { createPinia } from 'pinia'

import './index.css'
import { router } from './router'
import { useAuthStore } from './stores/auth'

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)
app.use(router)

// Initialize auth store to restore session
const authStore = useAuthStore()
authStore.loadUser()

app.mount('#app')
