import { FolderOpen, Home, LogOut, Settings } from 'lucide-vue-next'
import { computed, ref } from 'vue'

import { ROUTE_PATHS } from '@/constants/routes'
import type { SidebarItem } from '@/types/sidebar'
import { defineStore } from 'pinia'

export const useSidebarItemsStore = defineStore('sidebar-items', () => {
  const items = ref<SidebarItem[]>([
    {
      id: 'dashboard',
      title: 'Dashboard',
      url: ROUTE_PATHS.HOME,
      icon: Home,
    },
    {
      id: 'projects',
      title: 'Projects',
      url: ROUTE_PATHS.PROJECTS,
      icon: FolderOpen,
    },
  ])

  const userItems = ref<SidebarItem[]>([
    {
      id: 'settings',
      title: 'Settings',
      url: ROUTE_PATHS.SETTINGS,
      icon: Settings,
    },
    {
      id: 'logout',
      title: 'Logout',
      url: '/logout',
      icon: LogOut,
      action: 'logout',
    },
  ])

  const currentProjectId = ref<string | null>(null)
  const tableCount = ref<number>(0)
  const workflowCount = ref<number>(0)

  const addItem = (item: SidebarItem) => {
    items.value.push(item)
  }

  const removeItem = (id: SidebarItem['id']) => {
    items.value = items.value.filter((item) => item.id !== id)
  }

  const setProjectItems = (projectId: string, tables: number = 0, workflows: number = 0) => {
    currentProjectId.value = projectId
    tableCount.value = tables
    workflowCount.value = workflows
  }

  const updateProjectCounts = (tables: number, workflows: number) => {
    tableCount.value = tables
    workflowCount.value = workflows
  }

  const clearProjectItems = () => {
    currentProjectId.value = null
    tableCount.value = 0
    workflowCount.value = 0
  }

  return {
    items,
    userItems,
    currentProjectId: computed(() => currentProjectId.value),
    tableCount: computed(() => tableCount.value),
    workflowCount: computed(() => workflowCount.value),
    addItem,
    removeItem,
    setProjectItems,
    updateProjectCounts,
    clearProjectItems,
  }
})
