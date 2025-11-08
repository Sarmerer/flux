import { http } from '../http-client'

export interface Log {
  id: string
  timestamp: string
  level: string
  operation: string
  message: string
  project_id?: string
  database_id?: string
  table_id?: string
  workflow_id?: string
  user_id?: string
  metadata?: Record<string, any>
}

export interface LogsResponse {
  logs: Log[]
  total: number
  limit: number
  offset: number
}

export interface LogQueryParams {
  project_id?: string
  database_id?: string
  table_id?: string
  workflow_id?: string
  user_id?: string
  level?: string
  operation?: string
  start_time?: string
  end_time?: string
  limit?: number
  offset?: number
  order_by?: 'timestamp_asc' | 'timestamp_desc'
}

export const logsService = {
  getAll(params?: LogQueryParams) {
    const queryParams = new URLSearchParams()

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value))
        }
      })
    }

    const queryString = queryParams.toString()
    const url = queryString ? `/api/v1/logs?${queryString}` : '/api/v1/logs'

    return http.get<LogsResponse>(url)
  },

  getProjectLogs(projectId: string, params?: Omit<LogQueryParams, 'project_id'>) {
    const queryParams = new URLSearchParams()

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value))
        }
      })
    }

    const queryString = queryParams.toString()
    const url = queryString
      ? `/api/v1/projects/${projectId}/logs?${queryString}`
      : `/api/v1/projects/${projectId}/logs`

    return http.get<LogsResponse & { project_id: string }>(url)
  },

  cleanup(days: number) {
    return http.delete<{ deleted: number; message: string }>(`/api/v1/logs/cleanup?days=${days}`)
  },
}
