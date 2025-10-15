import type { Workflow } from '@/types/api'

import { http } from '../http-client'

export const workflowService = {
  getAll(projectId: string) {
    return http.get<Workflow[]>(`/projects/${projectId}/workflows`)
  },

  getById(id: string) {
    return http.get<Workflow>(`/workflows/${id}`)
  },

  create(projectId: string, data: any) {
    return http.post<Workflow>(`/projects/${projectId}/workflows`, data)
  },

  update(id: string, data: any) {
    return http.put<Workflow>(`/workflows/${id}`, data)
  },

  delete(id: string) {
    return http.delete<void>(`/workflows/${id}`)
  },
}
