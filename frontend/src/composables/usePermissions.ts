import type { Permission, Role } from '@/types/auth'
import { useAuthStore } from '@/stores/auth'

/**
 * Composable for checking user permissions and roles
 * Use this in components to control UI elements and actions based on permissions
 *
 * @example
 * ```vue
 * <script setup>
 * import { usePermissions } from '@/composables/usePermissions'
 *
 * const { can, hasRole, isAdmin } = usePermissions()
 * </script>
 *
 * <template>
 *   <button v-if="can('projects.create')" @click="createProject">
 *     Create Project
 *   </button>
 *
 *   <input
 *     v-model="projectName"
 *     :disabled="!can('projects.edit')"
 *   />
 *
 *   <div v-if="isAdmin()">
 *     <h2>Admin Settings</h2>
 *   </div>
 * </template>
 * ```
 */
export function usePermissions() {
  const authStore = useAuthStore()

  /**
   * Check if user has a specific permission
   * @param permission Single permission or array of permissions (all must be present)
   * @returns true if user has the permission(s)
   */
  const can = (permission: Permission | Permission[]): boolean => {
    return authStore.hasPermission(permission)
  }

  /**
   * Check if user has any of the specified permissions
   * @param permissions Array of permissions (at least one must be present)
   * @returns true if user has any of the permissions
   */
  const canAny = (permissions: Permission[]): boolean => {
    return authStore.hasAnyPermission(permissions)
  }

  /**
   * Check if user has a specific role
   * @param roles Single role or array of roles
   * @returns true if user has one of the roles
   */
  const hasRole = (roles: Role | Role[]): boolean => {
    return authStore.hasRole(roles)
  }

  /**
   * Check if user is an admin
   * @returns true if user is an admin
   */
  const isAdmin = (): boolean => {
    return authStore.isAdmin
  }

  /**
   * Get the current user's role
   * @returns The user's role or undefined if not authenticated
   */
  const getUserRole = (): Role | undefined => {
    return authStore.userRole
  }

  /**
   * Get all permissions for the current user
   * @returns Array of user permissions
   */
  const getUserPermissions = (): Permission[] => {
    return authStore.permissions
  }

  /**
   * Check if field should be disabled based on permission
   * Useful for form fields that should be disabled for users without edit permission
   */
  const isFieldDisabled = (permission: Permission): boolean => {
    return !can(permission)
  }

  /**
   * Check if action button should be hidden based on permission
   * Useful for action buttons that should be hidden for users without permission
   */
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
