import type { TableSchema } from '@/types/api'

import { http } from '../../http-client'

export const tableSchemaService = {
  createTable(projectId: string, tableId: string, tableName: string, schema: TableSchema) {
    return http.post<{ status: string }>(
      `/projects/${projectId}/tables/${tableId}/schema/create`,
      {
        table_name: tableName,
        schema,
      }
    )
  },

  dropTable(projectId: string, tableId: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${tableId}/schema/drop`)
  },

  addColumn(projectId: string, tableId: string, data: any) {
    return http.post<{ status: string }>(
      `/projects/${projectId}/tables/${tableId}/schema/columns/add`,
      data
    )
  },

  removeColumn(projectId: string, tableId: string, columnName: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${tableId}/schema/columns/remove`, {
      body: { column_name: columnName },
    })
  },

  modifyColumn(projectId: string, tableId: string, data: any) {
    return http.put<{ status: string }>(
      `/projects/${projectId}/tables/${tableId}/schema/columns/modify`,
      data
    )
  },
}
