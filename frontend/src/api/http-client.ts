import { safeJson } from '@/lib/utils'

const BASE_URL = import.meta.env.VITE_API_URL ?? ''

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export interface RequestOptions extends RequestInit {
  body?: any
}

async function request<T = unknown>(url: string, options: RequestOptions = {}): Promise<T> {
  const token = localStorage.getItem('token')
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options.headers,
  }

  const response = await fetch(url.startsWith('/') ? `${BASE_URL}${url}` : url, {
    ...options,
    headers,
    body:
      options.body && typeof options.body === 'object'
        ? JSON.stringify(options.body)
        : options.body,
  })

  if (!response.ok) {
    const body = await safeJson<{ message?: string }>(response)
    throw new Error(body?.message || `HTTP ${response.status}: ${response.statusText}`)
  }

  const data = await safeJson<T>(response)
  if (data === null) throw new Error('Invalid JSON in response')

  return data
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
