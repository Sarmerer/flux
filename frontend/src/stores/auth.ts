import { computed, ref } from 'vue'

import type { User } from '@/types/api'
import { defineStore } from 'pinia'

import { authService } from '@/api/services/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.email === 'admin@flow.com')
  const hasToken = computed(() => !!localStorage.getItem('auth_token'))

  const login = async (email: string, password: string) => {
    isLoading.value = true
    error.value = null

    try {
      await authService.login({ email, password })
      // Get user profile after successful login
      // For now, we'll create a mock user since we don't have the user ID yet
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
      // Auto-login after registration
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
      // For now, we'll create a mock user since we don't have a profile endpoint
      // In a real app, you would decode the JWT token or call a profile endpoint
      user.value = {
        id: '1',
        email: 'user@example.com',
        name: 'User',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
    } catch (err) {
      // Token might be invalid, clear it
      authService.logout()
      user.value = null
    } finally {
      isLoading.value = false
    }
  }

  return {
    user,
    isLoading,
    error,
    isAuthenticated,
    isAdmin,
    hasToken,
    login,
    register,
    logout,
    clearError,
    loadUser,
  }
})
