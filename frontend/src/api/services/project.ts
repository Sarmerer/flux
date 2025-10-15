import type { Project } from '@/types/api'

import { http } from '../http-client'

export const projectService = {
  getAll() {
    return http.get<Project[]>('/projects')
  },

  getById(id: string) {
    return http.get<Project>(`/projects/${id}`)
  },

  create(data: { name: string; description?: string }) {
    return http.post<Project>('/projects', data)
  },

  update(id: string, data: { name: string; description?: string }) {
    return http.put<Project>(`/projects/${id}`, data)
  },

  delete(id: string) {
    return http.delete<void>(`/projects/${id}`)
  },
}
