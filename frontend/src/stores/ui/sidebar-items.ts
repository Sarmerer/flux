import type { SidebarItem } from '@/types/sidebar'
import { defineStore } from 'pinia'

export const useSidebarItemsStore = defineStore('sidebar-items', {
  state: () => ({
    items: [
      {
        title: 'Home',
        url: '/',
        icon: 'Home',
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
  },
})
