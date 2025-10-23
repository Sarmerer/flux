import { useAuthStore } from '@/stores/auth'
import type { Permission, Role } from '@/types/auth'

export function usePermissions() {
  const authStore = useAuthStore()
  const can = (permission: Permission | Permission[]): boolean => {
    return authStore.hasPermission(permission)
  }

  const canAny = (permissions: Permission[]): boolean => {
    return authStore.hasAnyPermission(permissions)
  }

  const hasRole = (roles: Role | Role[]): boolean => {
    return authStore.hasRole(roles)
  }

  const isAdmin = (): boolean => {
    return authStore.isAdmin
  }

  const getUserRole = (): Role | undefined => {
    return authStore.userRole
  }

  const getUserPermissions = (): Permission[] => {
    return authStore.permissions
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
