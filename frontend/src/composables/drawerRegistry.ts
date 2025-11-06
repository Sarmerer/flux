import { markRaw } from 'vue'
import type { Component } from 'vue'
import type { ComponentProps } from 'vue-component-type-helpers'

import { useDrawerService } from './useDrawerService'

import TableEditorContent from '@/components/tables/drawers/TableEditorContent.vue'

type DrawerConfig<T extends Component> = {
  component: T
  props?: ComponentProps<T>
  width?: string
  onSave?: (data: any) => void
  onClose?: () => void
}

export const useDrawers = () => {
  const { openDrawer, ...rest } = useDrawerService()

  return {
    openTableEditor: (config: Omit<DrawerConfig<typeof TableEditorContent>, 'component'>) =>
      openDrawer({ component: markRaw(TableEditorContent), ...config }),

    ...rest,
  }
}
