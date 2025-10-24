import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { ProjectMember, ProjectMemberWithUser, Permission, Role } from '@/types/auth'
import { projectMemberService } from '@/api/services/projectMember'

export const useProjectMemberStore = defineStore('projectMember', () => {
  const currentProjectMember = ref<ProjectMember | null>(null)
  const projectMembers = ref<Map<string, ProjectMemberWithUser[]>>(new Map())
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const currentRole = computed(() => currentProjectMember.value?.role)
  const currentPermissions = computed(() => currentProjectMember.value?.permissions ?? [])

  const isAdmin = computed(() => currentRole.value === 'admin')

  const loadProjectMembers = async (projectId: string): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      const members = await projectMemberService.getProjectMembers(projectId)
      projectMembers.value.set(projectId, members)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load project members'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const loadMyProjectRole = async (projectId: string): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      const member = await projectMemberService.getMyProjectRole(projectId)
      currentProjectMember.value = member
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load project role'
      currentProjectMember.value = null
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const addProjectMember = async (
    projectId: string,
    email: string,
    role: Role,
    permissions: Permission[]
  ): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      await projectMemberService.addProjectMember(projectId, { email, role, permissions })
      await loadProjectMembers(projectId)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to add project member'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const updateProjectMember = async (
    projectId: string,
    memberId: string,
    role: Role,
    permissions: Permission[]
  ): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      await projectMemberService.updateProjectMember(projectId, memberId, { role, permissions })
      await loadProjectMembers(projectId)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update project member'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const removeProjectMember = async (projectId: string, memberId: string): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      await projectMemberService.removeProjectMember(projectId, memberId)
      await loadProjectMembers(projectId)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to remove project member'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const hasPermission = (permission: Permission | Permission[]): boolean => {
    if (!currentProjectMember.value) return false
    if (isAdmin.value) return true

    const perms = Array.isArray(permission) ? permission : [permission]
    return perms.every((p) => currentPermissions.value.includes(p))
  }

  const hasAnyPermission = (permissions: Permission[]): boolean => {
    if (!currentProjectMember.value) return false
    if (isAdmin.value) return true

    return permissions.some((p) => currentPermissions.value.includes(p))
  }

  const hasRole = (roles: Role | Role[]): boolean => {
    if (!currentProjectMember.value) return false

    const roleArray = Array.isArray(roles) ? roles : [roles]
    return roleArray.includes(currentProjectMember.value.role)
  }

  const clearCurrentProject = (): void => {
    currentProjectMember.value = null
    error.value = null
  }

  const getProjectMembers = (projectId: string): ProjectMemberWithUser[] => {
    return projectMembers.value.get(projectId) ?? []
  }

  return {
    currentProjectMember,
    projectMembers,
    isLoading,
    error,
    currentRole,
    currentPermissions,
    isAdmin,
    loadProjectMembers,
    loadMyProjectRole,
    addProjectMember,
    updateProjectMember,
    removeProjectMember,
    hasPermission,
    hasAnyPermission,
    hasRole,
    clearCurrentProject,
    getProjectMembers,
  }
})
