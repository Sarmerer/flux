import type { User } from '@/types/auth'

import { http } from '../http-client'

export const authService = {
  register(data: { email: string; password: string; name: string }) {
    return http.post<User>('/auth/register', data)
  },

  async login(data: { email: string; password: string }) {
    const res = await http.post<{ token: string }>('/auth/login', data)
    localStorage.setItem('auth_token', res.token)
    return res
  },

  logout() {
    localStorage.removeItem('auth_token')
  },

  getActiveUser() {
    return http.get<User>('/me')
  },

  getProfile(userId: string) {
    return http.get<User>(`/users/${userId}`)
  },
}
