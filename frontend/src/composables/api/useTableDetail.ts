import { computed, ref, watch } from 'vue'

import { tableDataService } from '@/api/services/table/data'
import { tableService } from '@/api/services/table/index'
import type { Table, TableColumn } from '@/types/api'

export interface TableDetailData {
  table: Table
  columns: TableColumn[]
  rows: Record<string, any>[]
  total: number
}

export function useTableDetail(projectId: string, tableId: string) {
  const data = ref<TableDetailData | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const fetchTableDetail = async () => {
    loading.value = true
    error.value = null

    try {
      const table = await tableService.getById(projectId, tableId)

      if (!table) {
        throw new Error('Table not found')
      }

      const tableDataResponse = await tableDataService.get(projectId, table.name, 1, 50)

      data.value = {
        table,
        columns: inferColumnsFromData(tableDataResponse.data),
        rows: tableDataResponse.data,
        total: tableDataResponse.total,
      }
    } catch (err) {
      error.value = err as Error
      data.value = null
    } finally {
      loading.value = false
    }
  }

  const inferColumnsFromData = (rows: any[]): TableColumn[] => {
    if (rows.length === 0) return []

    const firstRow = rows[0]
    return Object.keys(firstRow).map((key) => ({
      name: key,
      type: inferType(firstRow[key]),
      is_nullable: true,
      is_primary_key: key === 'id',
    }))
  }

  const inferType = (value: any): string => {
    if (value === null || value === undefined) return 'varchar'
    if (typeof value === 'number') return Number.isInteger(value) ? 'integer' : 'numeric'
    if (typeof value === 'boolean') return 'boolean'
    if (value instanceof Date || /^\d{4}-\d{2}-\d{2}/.test(value)) return 'timestamp'
    return 'varchar'
  }

  watch(
    () => [projectId, tableId],
    () => {
      if (projectId && tableId) {
        fetchTableDetail()
      }
    },
    { immediate: true }
  )

  return {
    data: computed(() => data.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh: fetchTableDetail,
  }
}
