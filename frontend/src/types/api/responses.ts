import type { User } from '../models/user'

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface ApiError {
  message: string
  code?: string
  details?: Record<string, any>
  field?: string
}

export interface AuthResponse {
  user: User
  token: string
}
