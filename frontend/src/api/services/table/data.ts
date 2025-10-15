import { http } from '../../http-client'

export const tableDataService = {
  get(projectId: string, tableName: string, page = 1, limit = 50) {
    return http.get<{ data: any[]; total: number; page: number; limit: number }>(
      `/projects/${projectId}/tables/${tableName}/data?page=${page}&limit=${limit}`
    )
  },

  insert(projectId: string, tableName: string, data: Record<string, any>) {
    return http.post<{ status: string }>(`/projects/${projectId}/tables/${tableName}/data`, data)
  },

  update(projectId: string, tableName: string, id: string, data: Record<string, any>) {
    return http.put<{ status: string }>(
      `/projects/${projectId}/tables/${tableName}/data/${id}`,
      data
    )
  },

  delete(projectId: string, tableName: string, id: string) {
    return http.delete<void>(`/projects/${projectId}/tables/${tableName}/data/${id}`)
  },
}
