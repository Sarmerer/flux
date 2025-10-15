import type { SidebarItem } from '@/types/sidebar'
import { defineStore } from 'pinia'
import { 
  Home, 
  FolderOpen, 
  Table, 
  Database, 
  Workflow, 
  Activity, 
  Settings,
  LogOut,
  User
} from 'lucide-vue-next'

export const useSidebarItemsStore = defineStore('sidebar-items', {
  state: () => ({
    items: [
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
    ] as Array<SidebarItem>,
    
    projectItems: [] as Array<SidebarItem>,
    
    userItems: [
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
    ] as Array<SidebarItem>,
  }),

  actions: {
    addItem(item: SidebarItem) {
      this.items.push(item)
    },

    removeItem(id: SidebarItem['id']) {
      this.items = this.items.filter((item) => item.id !== id)
    },

    setProjectItems(projectId: string) {
      this.projectItems = [
        {
          id: 'tables',
          title: 'Tables',
          url: `/projects/${projectId}/tables`,
          icon: Table,
        },
        {
          id: 'databases',
          title: 'Databases',
          url: `/projects/${projectId}/databases`,
          icon: Database,
        },
        {
          id: 'workflows',
          title: 'Workflows',
          url: `/projects/${projectId}/workflows`,
          icon: Workflow,
        },
        {
          id: 'activity',
          title: 'Activity',
          url: `/projects/${projectId}/activity`,
          icon: Activity,
        },
      ]
    },

    clearProjectItems() {
      this.projectItems = []
    },
  },
})
