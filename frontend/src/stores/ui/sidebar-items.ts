import type { SidebarItem } from '@/types/sidebar'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import {
  Home,
  FolderOpen,
  Table,
  Workflow,
  Activity,
  Settings,
  LogOut
} from 'lucide-vue-next'

export const useSidebarItemsStore = defineStore('sidebar-items', () => {
  // Static main navigation items
  const items = ref<SidebarItem[]>([
    {
      id: 'dashboard',
      title: 'Dashboard',
      url: '/dashboard',
      icon: Home,
    },
    {
      id: 'projects',
      title: 'Projects',
      url: '/projects',
      icon: FolderOpen,
    },
  ])

  // Project-scoped navigation items with reactive counts
  const projectItems = ref<SidebarItem[]>([])
  const currentProjectId = ref<string | null>(null)
  const tableCount = ref<number>(0)
  const workflowCount = ref<number>(0)

  // User navigation items
  const userItems = ref<SidebarItem[]>([
    {
      id: 'settings',
      title: 'Settings',
      url: '/settings',
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

  // Computed project items with badges
  const projectItemsWithBadges = computed((): SidebarItem[] => {
    if (!currentProjectId.value) return []

    return [
      {
        id: 'project-overview',
        title: 'Overview',
        url: `/projects/${currentProjectId.value}/overview`,
        icon: Home,
      },
      {
        id: 'project-tables',
        title: 'Tables',
        url: `/projects/${currentProjectId.value}/tables`,
        icon: Table,
        badge: tableCount.value > 0 ? tableCount.value.toString() : undefined,
      },
      {
        id: 'project-workflows',
        title: 'Workflows',
        url: `/projects/${currentProjectId.value}/workflows`,
        icon: Workflow,
        badge: workflowCount.value > 0 ? workflowCount.value.toString() : undefined,
      },
      {
        id: 'project-activity',
        title: 'Activity',
        url: `/projects/${currentProjectId.value}/activity`,
        icon: Activity,
      },
    ]
  })

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
    projectItems.value = projectItemsWithBadges.value
  }

  const updateProjectCounts = (tables: number, workflows: number) => {
    tableCount.value = tables
    workflowCount.value = workflows
    // Update project items with new counts
    if (currentProjectId.value) {
      projectItems.value = projectItemsWithBadges.value
    }
  }

  const clearProjectItems = () => {
    currentProjectId.value = null
    tableCount.value = 0
    workflowCount.value = 0
    projectItems.value = []
  }

  return {
    items,
    projectItems: computed(() => projectItemsWithBadges.value),
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
