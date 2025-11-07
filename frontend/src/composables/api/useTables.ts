import { computed } from 'vue'

import { tableService } from '@/api/services/table'
import { useCacheStore } from '@/stores/cache'
import type { Table } from '@/types'

import { useResourceCache } from '../data/useResourceCache'

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
    ttlMs: 60000,
    tags: ['tables', `project:${projectId}`],
    fetch: {
      fn: () => tableService.getAll(projectId),
      onMount: true,
    },
    realtime: {
      resourceId: projectId,
      events: ['table:created', 'table:updated', 'table:deleted'],
    },
  })

  const createTable = async (data: {
    name: string
    description?: string
    schema?: {
      columns: Array<{ name: string; type: string; nullable?: boolean; default?: string }>
      primary_key?: string[]
      indexes?: any[]
      foreign_keys?: any[]
    }
  }) => {
    const newTable = await tableService.create(projectId, data)
    refresh()

    cacheStore.invalidateByTag(`project:${projectId}`)
    return newTable
  }

  const updateTable = async (tableId: string, data: { name?: string; description?: string }) => {
    const updated = await tableService.update(projectId, tableId, data)
    refresh()

    cacheStore.invalidate(`table:${tableId}`)
    cacheStore.invalidateByTag(`project:${projectId}`)
    return updated
  }

  const deleteTable = async (tableId: string) => {
    await tableService.delete(projectId, tableId)
    refresh()

    cacheStore.invalidate(`table:${tableId}`)
    cacheStore.invalidateByTag(`project:${projectId}`)
  }

  const getTableById = async (tableId: string) => {
    return await tableService.getById(projectId, tableId)
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
