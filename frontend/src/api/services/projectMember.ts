import { http } from '../http-client'
import type { ProjectMember, ProjectMemberWithUser, Permission, Role } from '@/types/auth'

interface AddMemberRequest {
  email: string
  role: Role
  permissions: Permission[]
}

interface UpdateMemberRequest {
  role: Role
  permissions: Permission[]
}

export const projectMemberService = {
  getProjectMembers(projectId: string) {
    return http.get<ProjectMemberWithUser[]>(`/projects/${projectId}/members`)
  },

  getMyProjectRole(projectId: string) {
    return http.get<ProjectMember>(`/projects/${projectId}/members/me`)
  },

  addProjectMember(projectId: string, data: AddMemberRequest) {
    return http.post<ProjectMemberWithUser>(`/projects/${projectId}/members`, data)
  },

  updateProjectMember(projectId: string, memberId: string, data: UpdateMemberRequest) {
    return http.put<ProjectMemberWithUser>(`/projects/${projectId}/members/${memberId}`, data)
  },

  removeProjectMember(projectId: string, memberId: string) {
    return http.delete(`/projects/${projectId}/members/${memberId}`)
  },
}
