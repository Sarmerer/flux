import { computed, type Ref } from 'vue'
import { useTableData } from '@/composables/api'

export function useTableRowOperations(projectId: Ref<string>, tableId: Ref<string | null>) {
  const tableDataService = computed(() => {
    if (!tableId.value) return null
    return useTableData(projectId.value, tableId.value)
  })

  const insertRow = async (data: Record<string, any>) => {
    if (!tableDataService.value) {
      throw new Error('Table service not available')
    }
    await tableDataService.value.insertRow(data)
  }

  const updateRow = async (rowId: string, data: Record<string, any>) => {
    if (!tableDataService.value) {
      throw new Error('Table service not available')
    }
    await tableDataService.value.updateRow(rowId, data)
  }

  const deleteRow = async (rowId: string) => {
    if (!tableDataService.value) {
      throw new Error('Table service not available')
    }
    await tableDataService.value.deleteRow(rowId)
  }

  const deleteRows = async (rows: any[], primaryKeyColumn: string) => {
    if (!tableDataService.value) {
      throw new Error('Table service not available')
    }

    const errors: Error[] = []

    for (const row of rows) {
      try {
        await deleteRow(row[primaryKeyColumn])
      } catch (error) {
        errors.push(error as Error)
      }
    }

    if (errors.length > 0) {
      throw new Error(`Failed to delete ${errors.length} row(s)`)
    }
  }

  return {
    insertRow,
    updateRow,
    deleteRow,
    deleteRows,
  }
}
