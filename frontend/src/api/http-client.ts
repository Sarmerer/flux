import { safeJson } from '@/lib/utils'
import { AUTH_TOKEN_KEY } from '@/constants/auth'

const BASE_URL = import.meta.env.VITE_API_URL ?? ''

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export interface RequestOptions extends RequestInit {
  body?: any
}

async function request<T = unknown>(url: string, options: RequestOptions = {}): Promise<T> {
  url = url.startsWith('/') ? `${BASE_URL}${url}` : url
  const token = localStorage.getItem(AUTH_TOKEN_KEY)
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options.headers,
  }

  const response = await fetch(url, {
    ...options,
    headers,
    body:
      options.body && typeof options.body === 'object'
        ? JSON.stringify(options.body)
        : options.body,
  })

  if (!response.ok) {
    if (response.status === 401 || response.status === 403) {
      const wasAuthenticated = !!localStorage.getItem(AUTH_TOKEN_KEY)
      localStorage.removeItem(AUTH_TOKEN_KEY)

      if (wasAuthenticated && window.location.pathname !== '/login' && window.location.pathname !== '/register') {
        window.location.href = '/login'
      }
      throw new Error('Session expired. Please login again.')
    }

    const body = await safeJson<{ message?: string }>(response)
    throw new Error(body?.message || `HTTP ${response.status}: ${response.statusText}`)
  }

  if (response.status === 204 || response.headers.get('content-length') === '0') {
    return undefined as T
  }

  const contentType = response.headers.get('content-type')
  if (!contentType || !contentType.includes('application/json')) {
    const text = await response.text()
    console.error('Non-JSON response received:', {
      url: url.startsWith('/') ? `${BASE_URL}${url}` : url,
      contentType,
      status: response.status,
      text: text.substring(0, 200),
    })
    throw new Error(`Expected JSON response but got ${contentType || 'unknown'}`)
  }

  try {
    const data = (await response.json()) as T
    return data
  } catch (error) {
    console.error('JSON parse error:', { url, error })
    throw new Error('Invalid JSON in response')
  }
}

export const http = {
  get: <T>(url: string, opts?: RequestOptions) => request<T>(url, { ...opts, method: 'GET' }),
  post: <T>(url: string, body?: any, opts?: RequestOptions) =>
    request<T>(url, { ...opts, method: 'POST', body }),
  put: <T>(url: string, body?: any, opts?: RequestOptions) =>
    request<T>(url, { ...opts, method: 'PUT', body }),
  patch: <T>(url: string, body?: any, opts?: RequestOptions) =>
    request<T>(url, { ...opts, method: 'PATCH', body }),
  delete: <T>(url: string, opts?: RequestOptions) => request<T>(url, { ...opts, method: 'DELETE' }),
}
