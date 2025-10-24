import type { Table } from '@/types/api'

import { http } from '../../http-client'

export const tableService = {
  getAll(projectId: string) {
    return http.get<Table[]>(`/projects/${projectId}/tables`)
  },

  getById(projectId: string, id: string) {
    return http.get<Table>(`/projects/${projectId}/tables/${id}`)
  },

  create(projectId: string, data: { name: string; description?: string }) {
    return http.post<Table>(`/projects/${projectId}/tables`, data)
  },

  update(projectId: string, id: string, data: { name?: string; description?: string }) {
    return http.put<Table>(`/projects/${projectId}/tables/${id}`, data)
  },

  delete(projectId: string, id: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${id}`)
  },
}
