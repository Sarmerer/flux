import { computed } from 'vue'

import { tableService } from '@/api/services/table'
import { useCacheStore } from '@/stores/cache'
import type { Table } from '@/types/api'

import { useResourceCache } from '../data/useResourceCache'

/**
 * Composable for managing tables data for a specific project
 * Use this instead of direct API calls in components
 */
export function useTables(projectId: string) {
  const cacheKey = `tables:${projectId}`
  const cacheStore = useCacheStore()

  const {
    data: tables,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Table[]>(cacheKey, {
    fetchFn: () => tableService.getAll(projectId),
    subscribeToUpdates: true,
    resourceId: projectId,
    events: ['table:created', 'table:updated', 'table:deleted'],
    ttlMs: 60000,
    fetchOnMount: true,
    tags: ['tables', `project:${projectId}`],
  })

  const createTable = async (data: { name: string; description?: string }) => {
    const newTable = await tableService.create(projectId, data)
    refresh()

    cacheStore.invalidateByTag(`project:${projectId}`)
    return newTable
  }

  const updateTable = async (tableId: string, data: { name?: string; description?: string }) => {
    const updated = await tableService.update(tableId, data)
    refresh()

    cacheStore.invalidate(`table:${tableId}`)
    cacheStore.invalidateByTag(`project:${projectId}`)
    return updated
  }

  const deleteTable = async (tableId: string) => {
    await tableService.delete(tableId)
    refresh()

    cacheStore.invalidate(`table:${tableId}`)
    cacheStore.invalidateByTag(`project:${projectId}`)
  }

  const getTableById = async (tableId: string) => {
    return await tableService.getById(tableId)
  }

  return {
    tables: computed(() => tables.value ?? []),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    invalidate,
    createTable,
    updateTable,
    deleteTable,
    getTableById,
  }
}
