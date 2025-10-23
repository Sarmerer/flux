import { computed, ref } from 'vue'

import { authService } from '@/api/services/auth'
import type { Permission, Role, User } from '@/types/auth'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const hasToken = computed(() => !!localStorage.getItem('auth_token'))
  const userRole = computed(() => user.value?.role)
  const permissions = computed(() => user.value?.permissions ?? [])

  const login = async (email: string, password: string) => {
    isLoading.value = true
    error.value = null

    try {
      await authService.login({ email, password })

      user.value = {
        id: '1',
        email: email,
        name: 'User',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const register = async (name: string, email: string, password: string) => {
    isLoading.value = true
    error.value = null

    try {
      const newUser = await authService.register({ name, email, password })
      user.value = newUser

      await login(email, password)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Registration failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    isLoading.value = true
    try {
      authService.logout()
    } finally {
      user.value = null
      error.value = null
      isLoading.value = false
    }
  }

  const clearError = () => {
    error.value = null
  }

  const loadUser = async () => {
    if (!hasToken.value) return

    isLoading.value = true
    try {
      user.value = {
        id: '1',
        email: 'user@example.com',
        name: 'User',
        role: 'admin' as Role,
        permissions: [],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
    } catch (err) {
      authService.logout()
      user.value = null
    } finally {
      isLoading.value = false
    }
  }

  const hasPermission = (permission: Permission | Permission[]): boolean => {
    if (!user.value) return false

    // Admins have all permissions
    if (user.value.role === 'admin') return true

    const perms = Array.isArray(permission) ? permission : [permission]
    return perms.every((p) => user.value!.permissions.includes(p))
  }

  const hasAnyPermission = (permissions: Permission[]): boolean => {
    if (!user.value) return false

    // Admins have all permissions
    if (user.value.role === 'admin') return true

    return permissions.some((p) => user.value!.permissions.includes(p))
  }

  const hasRole = (roles: Role | Role[]): boolean => {
    if (!user.value) return false

    const roleArray = Array.isArray(roles) ? roles : [roles]
    return roleArray.includes(user.value.role)
  }

  const can = (permission: Permission | Permission[]): boolean => {
    return hasPermission(permission)
  }

  return {
    // State
    user,
    isLoading,
    error,

    // Computed
    isAuthenticated,
    isAdmin,
    hasToken,
    userRole,
    permissions,

    // Actions
    login,
    register,
    logout,
    clearError,
    loadUser,

    // Permission checks
    hasPermission,
    hasAnyPermission,
    hasRole,
    can,
  }
})
