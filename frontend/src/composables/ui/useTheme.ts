import { onMounted, ref, computed } from 'vue'

type Theme = 'dark' | 'light' | 'system'

const THEME_KEY = 'flow-theme'
const currentTheme = ref<Theme>('dark')

export function useTheme() {
  const getSystemTheme = (): 'dark' | 'light' => {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }

  const applyTheme = (theme: Theme) => {
    let actualTheme: 'dark' | 'light'

    if (theme === 'system') {
      actualTheme = getSystemTheme()
    } else {
      actualTheme = theme
    }

    // Update document class
    if (actualTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }

    // Save to localStorage
    localStorage.setItem(THEME_KEY, theme)
  }

  const setTheme = (theme: Theme) => {
    currentTheme.value = theme
    applyTheme(theme)
  }

  const isDark = computed(() => {
    if (currentTheme.value === 'system') {
      return getSystemTheme() === 'dark'
    }
    return currentTheme.value === 'dark'
  })

  const toggleTheme = () => {
    setTheme(isDark.value ? 'light' : 'dark')
  }

  // Load theme from localStorage or use default
  const loadTheme = () => {
    const saved = localStorage.getItem(THEME_KEY) as Theme | null
    const theme = saved || 'dark' // Default to dark
    setTheme(theme)
  }

  // Listen for system theme changes
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const handleSystemThemeChange = () => {
    if (currentTheme.value === 'system') {
      applyTheme('system')
    }
  }

  onMounted(() => {
    loadTheme()
    mediaQuery.addEventListener('change', handleSystemThemeChange)
  })

  return {
    theme: currentTheme,
    setTheme,
    toggleTheme,
  }
}
