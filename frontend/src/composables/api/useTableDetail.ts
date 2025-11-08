import { computed, ref, watch } from 'vue'

import { tableDataService } from '@/api/services/table/data'
import { tableService } from '@/api/services/table/index'
import type {
  Table,
  TableColumn,
  TableSchemaInfo,
  ForeignKeyInfo,
  ColumnInfo,
} from '@/types'

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

      const tableDataResponse = await tableDataService.get(projectId, tableId, 1, 50)

      const columns = extractColumnsFromSchemaInfo(table.schema)

      data.value = {
        table,
        columns,
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

  const extractColumnsFromSchemaInfo = (
    schemaInfo?: TableSchemaInfo
  ): TableColumn[] => {
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
