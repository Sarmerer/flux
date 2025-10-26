import type { TableSchema } from '@/types/api'

import { http } from '../../http-client'

export const tableSchemaService = {
  createTable(projectId: string, tableName: string, schema: TableSchema) {
    return http.post<{ status: string }>(
      `/projects/${projectId}/tables/${tableName}/schema/create`,
      {
        table_name: tableName,
        schema,
      }
    )
  },

  dropTable(projectId: string, tableName: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${tableName}/schema/drop`)
  },

  addColumn(projectId: string, tableName: string, data: any) {
    return http.post<{ status: string }>(
      `/projects/${projectId}/tables/${tableName}/schema/columns/add`,
      data
    )
  },

  removeColumn(projectId: string, tableName: string, columnName: string) {
    return http.delete<void>(
      `/projects/${projectId}/tables/${tableName}/schema/columns/remove`,
      { column_name: columnName }
    )
  },

  modifyColumn(projectId: string, tableName: string, data: any) {
    return http.put<{ status: string }>(
      `/projects/${projectId}/tables/${tableName}/schema/columns/modify`,
      data
    )
  },
}
