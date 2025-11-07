import { useProjectMemberStore } from '@/stores/projectMember'
import type { Permission, Role } from '@/types'

export function usePermissions() {
  const projectMemberStore = useProjectMemberStore()

  const can = (permission: Permission | Permission[]): boolean => {
    return projectMemberStore.hasPermission(permission)
  }

  const canAny = (permissions: Permission[]): boolean => {
    return projectMemberStore.hasAnyPermission(permissions)
  }

  const hasRole = (roles: Role | Role[]): boolean => {
    return projectMemberStore.hasRole(roles)
  }

  const isAdmin = (): boolean => {
    return projectMemberStore.isAdmin
  }

  const getUserRole = (): Role | undefined => {
    return projectMemberStore.currentRole
  }

  const getUserPermissions = (): Permission[] => {
    return projectMemberStore.currentPermissions
  }

  const isFieldDisabled = (permission: Permission): boolean => {
    return !can(permission)
  }

  const shouldHideAction = (permission: Permission): boolean => {
    return !can(permission)
  }

  return {
    can,
    canAny,
    hasRole,
    isAdmin,
    getUserRole,
    getUserPermissions,
    isFieldDisabled,
    shouldHideAction,
  }
}
