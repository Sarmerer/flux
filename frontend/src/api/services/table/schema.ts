import type { TableSchema } from '@/types/api'

import { http } from '../../http-client'

export const tableSchemaService = {
  createTable(projectId: string, tableName: string, schema: TableSchema) {
    return http.post<{ status: string }>(`/projects/${projectId}/tables/${tableName}/create`, {
      table_name: tableName,
      schema,
    })
  },

  dropTable(projectId: string, tableName: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${tableName}/drop`)
  },

  addColumn(projectId: string, tableName: string, data: any) {
    return http.post<{ status: string }>(`/projects/${projectId}/tables/${tableName}/columns`, data)
  },

  removeColumn(projectId: string, tableName: string, columnName: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${tableName}/columns/${columnName}`)
  },

  modifyColumn(projectId: string, tableName: string, data: any) {
    return http.put<{ status: string }>(
      `/projects/${projectId}/tables/${tableName}/columns/${data.column_name}`,
      data
    )
  },
}
