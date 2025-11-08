import { http } from '../../http-client'

export const tableDataService = {
  get(projectId: string, tableId: string, page = 1, limit = 50) {
    return http.get<{ rows: any[]; total: number; page: number; limit: number }>(
      `/projects/${projectId}/tables/${tableId}/data?page=${page}&limit=${limit}`
    )
  },

  insert(projectId: string, tableId: string, data: Record<string, any>) {
    return http.post<Record<string, any>>(`/projects/${projectId}/tables/${tableId}/data`, data)
  },

  update(projectId: string, tableId: string, id: string, data: Record<string, any>) {
    return http.put<{ status: string }>(`/projects/${projectId}/tables/${tableId}/data/${id}`, data)
  },

  delete(projectId: string, tableId: string, id: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${tableId}/data/${id}`)
  },
}
