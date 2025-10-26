import type { Workflow, WorkflowCreateRequest, WorkflowUpdateRequest } from '@/types/api'

import { http } from '../http-client'

export const workflowService = {
  getAll(projectId: string) {
    return http.get<Workflow[]>(`/projects/${projectId}/workflows`)
  },

  getById(projectId: string, id: string) {
    return http.get<Workflow>(`/projects/${projectId}/workflows/${id}`)
  },

  create(projectId: string, data: WorkflowCreateRequest) {
    return http.post<Workflow>(`/projects/${projectId}/workflows`, data)
  },

  update(projectId: string, id: string, data: WorkflowUpdateRequest) {
    return http.put<Workflow>(`/projects/${projectId}/workflows/${id}`, data)
  },

  delete(projectId: string, id: string) {
    return http.delete<void>(`/projects/${projectId}/workflows/${id}`)
  },

  toggle(projectId: string, id: string, isActive: boolean) {
    return http.post<Workflow>(`/projects/${projectId}/workflows/${id}/toggle`, { is_active: isActive })
  },

  execute(projectId: string, id: string) {
    return http.post<{ message: string; result: any }>(
      `/projects/${projectId}/workflows/${id}/execute`,
      {}
    )
  },
}
