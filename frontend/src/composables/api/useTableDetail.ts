import { computed, ref, watch } from 'vue'

import { tableDataService } from '@/api/services/table/data'
import { tableService } from '@/api/services/table/index'
import type { ColumnInfo, ForeignKeyInfo, Table, TableColumn, TableSchemaInfo } from '@/types'

import { useToast } from '../ui'

export interface TableDetailData {
  table: Table
  columns: TableColumn[]
  rows: Record<string, any>[]
  total: number
}

export interface UseTableDetailOptions {
  page?: number
  limit?: number
  autoFetch?: boolean
}

export function useTableDetail(
  projectId: string,
  tableId: string,
  options: UseTableDetailOptions = {}
) {
  const { page: initialPage = 1, limit: initialLimit = 50, autoFetch = true } = options

  const toast = useToast()
  const data = ref<TableDetailData | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const currentPage = ref(initialPage)
  const pageLimit = ref(initialLimit)

  const columns = computed(() => data.value?.columns ?? [])
  const rows = computed(() => data.value?.rows ?? [])
  const totalRows = computed(() => data.value?.total ?? 0)

  const fetchTableDetail = async (forceRefresh = false) => {
    if (!projectId || !tableId) {
      data.value = null
      return
    }

    if (loading.value && !forceRefresh) return

    loading.value = true
    error.value = null

    try {
      const [table, tableDataResponse] = await Promise.all([
        tableService.getById(projectId, tableId),
        tableDataService.get(projectId, tableId, currentPage.value, pageLimit.value),
      ])

      if (!table) {
        throw new Error('Table not found')
      }

      const extractedColumns = extractColumnsFromSchemaInfo(table.schema)

      data.value = {
        table,
        columns: extractedColumns,
        rows: tableDataResponse.data,
        total: tableDataResponse.total,
      }

      error.value = null
    } catch (err) {
      error.value = err as Error
      data.value = null
      console.error('[useTableDetail] Failed to fetch table details:', err)

      toast.error('Failed to load table', (err as Error).message)
    } finally {
      loading.value = false
    }
  }

  const refresh = () => fetchTableDetail(true)

  const goToPage = (page: number) => {
    currentPage.value = page
    fetchTableDetail()
  }

  const setPageLimit = (limit: number) => {
    pageLimit.value = limit
    currentPage.value = 1
    fetchTableDetail()
  }

  watch([() => projectId, () => tableId], () => {
    currentPage.value = 1
    data.value = null
    if (autoFetch && projectId && tableId) {
      fetchTableDetail()
    }
  })

  watch([currentPage, pageLimit], () => {
    if (autoFetch && projectId && tableId) {
      fetchTableDetail()
    }
  })

  if (autoFetch && projectId && tableId) {
    fetchTableDetail()
  }

  return {
    data: computed(() => data.value),
    columns,
    rows,
    totalRows,
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    currentPage: computed(() => currentPage.value),
    pageLimit: computed(() => pageLimit.value),
    refresh,
    goToPage,
    setPageLimit,
  }
}

function extractColumnsFromSchemaInfo(schemaInfo?: TableSchemaInfo): TableColumn[] {
  if (!schemaInfo?.columns || !Array.isArray(schemaInfo.columns)) {
    return []
  }

  const foreignKeyMap = new Map<string, ForeignKeyInfo>()
  schemaInfo.foreign_keys?.forEach((fk: ForeignKeyInfo) => {
    fk.column_names.forEach((colName: string) => {
      foreignKeyMap.set(colName, fk)
    })
  })

  return schemaInfo.columns.map((col: ColumnInfo) => {
    const foreignKey = foreignKeyMap.get(col.name)
    const isPrimaryKey = schemaInfo.primary_keys?.includes(col.name) || false

    return {
      name: col.name,
      type: col.data_type,
      is_nullable: col.is_nullable,
      is_primary_key: isPrimaryKey,
      default_value: col.default_value || undefined,
      is_foreign_key: !!foreignKey,
      foreign_table: foreignKey?.referenced_table,
      foreign_column: foreignKey?.referenced_columns[0],
    }
  })
}
