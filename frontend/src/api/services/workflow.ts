import type { Workflow, WorkflowCreateRequest, WorkflowUpdateRequest } from '@/types/api'

import { http } from '../http-client'

export const workflowService = {
  /**
   * Get all workflows for a project
   */
  getAll(projectId: string) {
    return http.get<Workflow[]>(`/projects/${projectId}/workflows`)
  },

  /**
   * Get a specific workflow by ID
   */
  getById(projectId: string, id: string) {
    return http.get<Workflow>(`/projects/${projectId}/workflows/${id}`)
  },

  /**
   * Create a new workflow
   */
  create(projectId: string, data: WorkflowCreateRequest) {
    return http.post<Workflow>(`/projects/${projectId}/workflows`, data)
  },

  /**
   * Update an existing workflow
   */
  update(projectId: string, id: string, data: WorkflowUpdateRequest) {
    return http.put<Workflow>(`/projects/${projectId}/workflows/${id}`, data)
  },

  /**
   * Delete a workflow
   */
  delete(projectId: string, id: string) {
    return http.delete<void>(`/projects/${projectId}/workflows/${id}`)
  },

  /**
   * Toggle workflow active/inactive status
   */
  toggle(projectId: string, id: string) {
    return http.post<Workflow>(`/projects/${projectId}/workflows/${id}/toggle`, {})
  },

  /**
   * Execute/test a workflow manually
   */
  execute(projectId: string, id: string) {
    return http.post<{ message: string; result: any }>(
      `/projects/${projectId}/workflows/${id}/execute`,
      {}
    )
  },
}
