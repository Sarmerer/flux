import { AUTH_TOKEN_KEY } from '@/constants/auth'
import type { User } from '@/types'

import { http } from '../http-client'

export const authService = {
  register(data: { email: string; password: string; name: string }) {
    return http.post<User>('/auth/register', data)
  },

  async login(data: { email: string; password: string }) {
    const res = await http.post<{ token: string }>('/auth/login', data)
    localStorage.setItem(AUTH_TOKEN_KEY, res.token)
    return res
  },

  logout() {
    localStorage.removeItem(AUTH_TOKEN_KEY)
  },

  getActiveUser() {
    return http.get<User>('/me')
  },

  getProfile(userId: string) {
    return http.get<User>(`/users/${userId}`)
  },
}
