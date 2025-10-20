import type { Table } from '@/types/api'

import { http } from '../../http-client'

export const tableService = {
  getAll(projectId: string) {
    return http.get<Table[]>(`/projects/${projectId}/tables`)
  },

  getById(id: string) {
    return http.get<Table>(`/tables/${id}`)
  },

  create(projectId: string, data: { name: string; description?: string }) {
    return http.post<Table>(`/projects/${projectId}/tables`, data)
  },

  update(id: string, data: { name?: string; description?: string }) {
    return http.put<Table>(`/tables/${id}`, data)
  },

  delete(id: string) {
    return http.delete<void>(`/tables/${id}`)
  },
}
