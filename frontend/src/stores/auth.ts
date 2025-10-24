import { computed, ref } from 'vue'

import { authService } from '@/api/services/auth'
import type { User } from '@/types/auth'
import { defineStore } from 'pinia'

const TOKEN_KEY = 'auth_token'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)
  const hasToken = computed(() => !!localStorage.getItem(TOKEN_KEY))

  const login = async (email: string, password: string): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      await authService.login({ email, password })
      await loadUser()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const register = async (name: string, email: string, password: string): Promise<void> => {
    isLoading.value = true
    error.value = null

    try {
      await authService.register({ name, email, password })
      await login(email, password)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Registration failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const logout = (): void => {
    authService.logout()
    user.value = null
    error.value = null
  }

  const loadUser = async (): Promise<void> => {
    if (!hasToken.value) {
      user.value = null
      return
    }

    isLoading.value = true
    try {
      const userData = await authService.getActiveUser()
      user.value = userData
      error.value = null
    } catch (err) {
      authService.logout()
      user.value = null
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const clearError = (): void => {
    error.value = null
  }

  return {
    user,
    isLoading,
    error,
    isAuthenticated,
    hasToken,
    login,
    register,
    logout,
    loadUser,
    clearError,
  }
})
