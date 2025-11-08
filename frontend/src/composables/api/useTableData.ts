import { ref } from 'vue'

import { tableDataService } from '@/api/services/table/data'
import { useCacheStore } from '@/stores/cache'

export function useTableData(projectId: string, tableId: string) {
  const cacheStore = useCacheStore()
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const insertRow = async (data: Record<string, any>) => {
    loading.value = true
    error.value = null

    try {
      const result = await tableDataService.insert(projectId, tableId, data)
      cacheStore.invalidateByTag(`table:${tableId}`)
      cacheStore.invalidate(`table-detail:${projectId}:${tableId}`)
      return result
    } catch (err) {
      error.value = err as Error
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateRow = async (rowId: string, data: Record<string, any>) => {
    loading.value = true
    error.value = null

    try {
      const result = await tableDataService.update(projectId, tableId, rowId, data)
      cacheStore.invalidateByTag(`table:${tableId}`)
      cacheStore.invalidate(`table-detail:${projectId}:${tableId}`)
      return result
    } catch (err) {
      error.value = err as Error
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteRow = async (rowId: string) => {
    loading.value = true
    error.value = null

    try {
      await tableDataService.delete(projectId, tableId, rowId)
      cacheStore.invalidateByTag(`table:${tableId}`)
      cacheStore.invalidate(`table-detail:${projectId}:${tableId}`)
    } catch (err) {
      error.value = err as Error
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    insertRow,
    updateRow,
    deleteRow,
  }
}
