import { computed } from 'vue'

import { tableDataService } from '@/api/services/table/data'
import { tableService } from '@/api/services/table/index'
import type { ColumnInfo, ForeignKeyInfo, Table, TableColumn, TableSchemaInfo } from '@/types'

import { useResourceCache } from '../data/useResourceCache'

export interface TableDetailData {
  table: Table
  columns: TableColumn[]
  rows: Record<string, any>[]
  total: number
}

export function useTableDetail(projectId: string, tableId: string) {
  const cacheKey = `table-detail:${projectId}:${tableId}`

  const { data, loading, error, refresh, invalidate } = useResourceCache<TableDetailData>(
    cacheKey,
    {
      ttlMs: 60000,
      tags: ['table-details', `table:${tableId}`, `project:${projectId}`],
      fetch: {
        fn: async () => {
          const table = await tableService.getById(projectId, tableId)

          if (!table) {
            throw new Error('Table not found')
          }

          const tableDataResponse = await tableDataService.get(projectId, tableId, 1, 50)

          const columns = extractColumnsFromSchemaInfo(table.schema)

          return {
            table,
            columns,
            rows: tableDataResponse.rows,
            total: tableDataResponse.total,
          }
        },
        onMount: true,
      },
      realtime: {
        resourceId: tableId,
        events: ['table:updated', 'table:data:created', 'table:data:updated', 'table:data:deleted'],
      },
      toastOnError: {
        title: 'Failed to load table',
        description: 'Could not fetch table details',
      },
    }
  )

  return {
    data: computed(() => data.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    invalidate,
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
