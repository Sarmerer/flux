import type { TableSchema } from '@/types'

import { http } from '../../http-client'

interface ColumnDefinition {
  name: string
  type: string
  nullable?: boolean
  default?: string
  primary_key?: boolean
}

interface ForeignKeyDefinition {
  name: string
  columns: string[]
  reference_table: string
  reference_columns: string[]
}

export const tableSchemaService = {
  updateSchema(projectId: string, tableId: string, schema: TableSchema) {
    return http.put<{ status: string }>(`/api/v1/projects/${projectId}/tables/${tableId}/schema/update`, {
      schema,
    })
  },

  addColumn(projectId: string, tableId: string, column: ColumnDefinition) {
    return http.post<{ status: string }>(
      `/api/v1/projects/${projectId}/tables/${tableId}/schema/columns`,
      { column }
    )
  },

  removeColumn(projectId: string, tableId: string, columnName: string) {
    return http.delete<void>(`/api/v1/projects/${projectId}/tables/${tableId}/schema/columns`, {
      body: { column_name: columnName },
    })
  },

  modifyColumn(projectId: string, tableId: string, columnName: string, column: ColumnDefinition) {
    return http.put<{ status: string }>(`/api/v1/projects/${projectId}/tables/${tableId}/schema/columns`, {
      column_name: columnName,
      column,
    })
  },

  addForeignKey(projectId: string, tableId: string, foreignKey: ForeignKeyDefinition) {
    return http.post<{ status: string }>(
      `/api/v1/projects/${projectId}/tables/${tableId}/schema/foreign-keys`,
      { foreign_key: foreignKey }
    )
  },

  removeForeignKey(projectId: string, tableId: string, foreignKeyName: string) {
    return http.delete<void>(`/api/v1/projects/${projectId}/tables/${tableId}/schema/foreign-keys`, {
      body: { foreign_key_name: foreignKeyName },
    })
  },
}
