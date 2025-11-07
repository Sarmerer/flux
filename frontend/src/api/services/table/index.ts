import type { Table } from '@/types'

import { http } from '../../http-client'

interface TableColumn {
  name: string
  type: string
  nullable?: boolean
  default?: string
}

interface TableSchema {
  columns: TableColumn[]
  primary_key?: string[]
  indexes?: any[]
  foreign_keys?: any[]
}

interface CreateTableData {
  name: string
  description?: string
  schema?: TableSchema
}

export const tableService = {
  getAll(projectId: string) {
    return http.get<Table[]>(`/projects/${projectId}/tables`)
  },

  getById(projectId: string, id: string) {
    return http.get<Table>(`/projects/${projectId}/tables/${id}`)
  },

  create(projectId: string, data: CreateTableData) {
    return http.post<Table>(`/projects/${projectId}/tables`, data)
  },

  update(projectId: string, id: string, data: { name?: string; description?: string }) {
    return http.put<Table>(`/projects/${projectId}/tables/${id}`, data)
  },

  delete(projectId: string, id: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${id}`)
  },
}
