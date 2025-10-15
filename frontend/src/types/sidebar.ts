import type { Component } from 'vue'

export interface SidebarItem {
  id: string
  title: string
  url: string
  icon?: Component | string
  children?: SidebarItem[]
  action?: string
}
