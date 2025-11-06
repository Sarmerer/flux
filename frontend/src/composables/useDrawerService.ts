import { nanoid } from 'nanoid'
import type { Component } from 'vue'
import { readonly, shallowRef } from 'vue'

export interface DrawerLayer {
  id: string
  component: Component
  props?: Record<string, any>
  width?: string
  onSave?: (data: any) => void
  onClose?: () => void
}

const drawerStack = shallowRef<DrawerLayer[]>([])

export const useDrawerService = () => {
  const openDrawer = (layer: Omit<DrawerLayer, 'id'>) => {
    const id = nanoid()
    drawerStack.value = [...drawerStack.value, { ...layer, id }]
    return id
  }

  const closeDrawer = (id?: string) => {
    if (id) {
      const index = drawerStack.value.findIndex((l) => l.id === id)
      if (index !== -1) {
        const layer = drawerStack.value[index]
        layer?.onClose?.()
        drawerStack.value = drawerStack.value.filter((l) => l.id !== id)
      }
    } else {
      const lastLayer = drawerStack.value[drawerStack.value.length - 1]
      if (lastLayer) {
        lastLayer.onClose?.()
        drawerStack.value = drawerStack.value.slice(0, -1)
      }
    }
  }

  const closeAll = () => {
    drawerStack.value.forEach((layer) => layer.onClose?.())
    drawerStack.value = []
  }

  const replaceDrawer = (id: string, layer: Omit<DrawerLayer, 'id'>) => {
    const index = drawerStack.value.findIndex((l) => l.id === id)
    if (index !== -1) {
      const newStack = [...drawerStack.value]
      newStack[index] = { ...layer, id }
      drawerStack.value = newStack
    }
  }

  return {
    drawerStack: readonly(drawerStack),
    openDrawer,
    closeDrawer,
    closeAll,
    replaceDrawer,
  }
}
