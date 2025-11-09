import { markRaw } from 'vue'
import type { Component } from 'vue'
import type { ComponentProps } from 'vue-component-type-helpers'

import { useDrawerService } from './useDrawerService'

import DeleteRowContent from '@/components/tables/drawers/DeleteRowContent.vue'
import RowEditorContent from '@/components/tables/drawers/RowEditorContent.vue'
import TableEditorContent from '@/components/tables/drawers/TableEditorContent.vue'
import WorkflowEditorContent from '@/components/workflows/drawers/WorkflowEditorContent.vue'

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

    openInsertRow: (config: Omit<DrawerConfig<typeof RowEditorContent>, 'component'>) =>
      openDrawer({ component: markRaw(RowEditorContent), width: '600px', ...config }),

    openEditRow: (config: Omit<DrawerConfig<typeof RowEditorContent>, 'component'>) =>
      openDrawer({ component: markRaw(RowEditorContent), width: '600px', ...config }),

    openDeleteRow: (config: Omit<DrawerConfig<typeof DeleteRowContent>, 'component'>) =>
      openDrawer({ component: markRaw(DeleteRowContent), width: '500px', ...config }),

    openWorkflowEditor: (config: Omit<DrawerConfig<typeof WorkflowEditorContent>, 'component'>) =>
      openDrawer({ component: markRaw(WorkflowEditorContent), width: '700px', ...config }),

    ...rest,
  }
}
