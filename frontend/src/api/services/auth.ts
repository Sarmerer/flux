import type { User } from '@/types/api'

import { http } from '../http-client'

export const authService = {
  register(data: { email: string; password: string; name: string }) {
    return http.post<User>('/auth/register', data)
  },

  async login(data: { email: string; password: string }) {
    const res = await http.post<{ token: string }>('/auth/login', data)
    localStorage.setItem('token', res.token)
    return res
  },

  logout() {
    localStorage.removeItem('token')
  },

  getProfile(userId: string) {
    return http.get<User>(`/users/${userId}`)
  },
}
