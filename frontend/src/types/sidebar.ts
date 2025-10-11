export interface SidebarItem {
  id: string
  title: string
  url: string
  icon?: string
  children?: SidebarItem[]
}
