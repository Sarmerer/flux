import type { Ref } from 'vue'

import type { TableColumn } from '@/types'

import { useToast } from '../ui'

export function useTableRowActions(
  columns: Ref<TableColumn[]>,
  onInsert: (data: Record<string, any>) => Promise<void>,
  onUpdate: (id: string, data: Record<string, any>) => Promise<void>,
  onDelete: (id: string) => Promise<void>
) {
  const toast = useToast()

  const getPrimaryKeyColumn = (): TableColumn | undefined => {
    return columns.value.find((col) => col.is_primary_key)
  }

  const getRowId = (row: Record<string, any>): string | null => {
    const pkColumn = getPrimaryKeyColumn()
    if (!pkColumn) {
      toast.error('Error', 'No primary key found for this table')
      return null
    }
    return row[pkColumn.name]
  }

  const insertRow = async (data: Record<string, any>) => {
    try {
      await onInsert(data)
      toast.success('Success', 'Row inserted successfully')
    } catch (error: any) {
      toast.error('Error', error.message || 'Failed to insert row')
      throw error
    }
  }

  const updateRow = async (row: Record<string, any>, data: Record<string, any>) => {
    const rowId = getRowId(row)
    if (!rowId) return

    try {
      await onUpdate(rowId, data)
      toast.success('Success', 'Row updated successfully')
    } catch (error: any) {
      toast.error('Error', error.message || 'Failed to update row')
      throw error
    }
  }

  const deleteRow = async (row: Record<string, any>) => {
    const rowId = getRowId(row)
    if (!rowId) return

    try {
      await onDelete(rowId)
      toast.success('Success', 'Row deleted successfully')
    } catch (error: any) {
      toast.error('Error', error.message || 'Failed to delete row')
      throw error
    }
  }

  return {
    getPrimaryKeyColumn,
    getRowId,
    insertRow,
    updateRow,
    deleteRow,
  }
}
