import { ref } from 'vue'
import { defineStore } from 'pinia'

/**
 * UI state store - only for truly global UI state
 */
export const useUIStore = defineStore('ui', () => {
  const sidebarOpen = ref(true)
  const sidebarCollapsed = ref(false)

  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
    localStorage.setItem('sidebar-open', String(sidebarOpen.value))
  }

  const toggleSidebarCollapse = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem('sidebar-collapsed', String(sidebarCollapsed.value))
  }

  const setSidebarOpen = (open: boolean) => {
    sidebarOpen.value = open
    localStorage.setItem('sidebar-open', String(open))
  }

  // Load from localStorage
  const savedOpen = localStorage.getItem('sidebar-open')
  if (savedOpen !== null) {
    sidebarOpen.value = savedOpen === 'true'
  }

  const savedCollapsed = localStorage.getItem('sidebar-collapsed')
  if (savedCollapsed !== null) {
    sidebarCollapsed.value = savedCollapsed === 'true'
  }

  return {
    sidebarOpen,
    sidebarCollapsed,
    toggleSidebar,
    toggleSidebarCollapse,
    setSidebarOpen,
  }
})
