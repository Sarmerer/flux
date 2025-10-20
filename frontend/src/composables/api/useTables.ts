import { computed } from 'vue'

import { tableService } from '@/api/services/table'
import type { Table } from '@/types/api'

import { useResourceCache } from '../data/useResourceCache'

/**
 * Composable for managing tables data for a specific project
 * Use this instead of direct API calls in components
 */
export function useTables(projectId: string) {
  const cacheKey = `tables:${projectId}`

  const {
    data: tables,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Table[]>(cacheKey, {
    fetchFn: () => tableService.getAll(projectId),
    subscribeToUpdates: true,
    scope: 'project',
    resourceId: projectId,
    events: ['table:created', 'table:updated', 'table:deleted'],
    ttl: 60000, // 1 minute
    refetchOnMount: true,
    optimisticUpdate: true,
  })

  const createTable = async (data: { name: string; description?: string }) => {
    const newTable = await tableService.create(projectId, data)
    refresh() // Refresh to get updated list
    return newTable
  }

  const updateTable = async (tableId: string, data: { name?: string; description?: string }) => {
    const updated = await tableService.update(tableId, data)
    refresh() // Refresh to get updated list
    return updated
  }

  const deleteTable = async (tableId: string) => {
    await tableService.delete(tableId)
    refresh() // Refresh to get updated list
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
