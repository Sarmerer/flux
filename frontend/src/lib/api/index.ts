import type {
  ActivityLog,
  Database,
  Project,
  Table,
  TableSchema,
  User,
  Workflow,
} from '@/types/api'
import { StatusCodes } from 'http-status-codes'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class ApiClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.token = localStorage.getItem('auth_token')
  }

  get hasToken(): boolean {
    return !!this.token
  }

  setToken(token: string) {
    this.token = token
    localStorage.setItem('auth_token', token)
  }

  clearToken() {
    this.token = null
    localStorage.removeItem('auth_token')
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    }

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: 'Unknown error' }))
      throw new Error(error.message || `HTTP ${response.status}`)
    }

    if (response.status === StatusCodes.NO_CONTENT) {
      return {} as T
    }

    return response.json() as Promise<T>
  }

  // Auth endpoints
  async register(data: { email: string; password: string; name: string }) {
    return this.request<User>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async login(data: { email: string; password: string }) {
    const response = await this.request<{ token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    this.setToken(response.token)
    return response
  }

  async logout() {
    this.clearToken()
  }

  async getProfile(userId: string) {
    return this.request<User>(`/users/${userId}`)
  }

  // Project endpoints
  async getProjects() {
    return this.request<Project[]>('/projects')
  }

  async getProject(projectId: string) {
    return this.request<Project>(`/projects/${projectId}`)
  }

  async createProject(data: { name: string; description?: string }) {
    return this.request<Project>('/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateProject(projectId: string, data: { name: string; description?: string }) {
    return this.request<Project>(`/projects/${projectId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteProject(projectId: string) {
    return this.request<void>(`/projects/${projectId}`, {
      method: 'DELETE',
    })
  }

  // Table endpoints
  async getTables(projectId: string) {
    return this.request<Table[]>(`/projects/${projectId}/tables`)
  }

  async getTable(tableId: string) {
    return this.request<Table>(`/tables/${tableId}`)
  }

  async createTable(projectId: string, data: { name: string; description?: string }) {
    return this.request<Table>(`/projects/${projectId}/tables`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateTable(tableId: string, data: { name: string; description?: string }) {
    return this.request<Table>(`/tables/${tableId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteTable(tableId: string) {
    return this.request<void>(`/tables/${tableId}`, {
      method: 'DELETE',
    })
  }

  // Database endpoints
  async getProjectDatabases(projectId: string) {
    return this.request<Database[]>(`/projects/${projectId}/databases`)
  }

  async createProjectDatabase(
    projectId: string,
    data: {
      name: string
      host: string
      port: number
      username: string
      password: string
      database: string
    }
  ) {
    return this.request<Database>(`/projects/${projectId}/databases`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateProjectDatabase(
    databaseId: string,
    data: {
      name: string
      host: string
      port: number
      username: string
      password: string
      database: string
    }
  ) {
    return this.request<Database>(`/databases/${databaseId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteProjectDatabase(databaseId: string) {
    return this.request<void>(`/databases/${databaseId}`, {
      method: 'DELETE',
    })
  }

  async testDatabaseConnection(databaseId: string) {
    return this.request<{ status: string }>(`/databases/${databaseId}/test`, {
      method: 'POST',
    })
  }

  // Table schema mutation endpoints
  async createTableInDatabase(projectId: string, tableName: string, schema: TableSchema) {
    return this.request<{ status: string }>(`/projects/${projectId}/tables/${tableName}/create`, {
      method: 'POST',
      body: JSON.stringify({ table_name: tableName, schema }),
    })
  }

  async dropTableFromDatabase(projectId: string, tableName: string) {
    return this.request<void>(`/projects/${projectId}/tables/${tableName}/drop`, {
      method: 'DELETE',
    })
  }

  async addColumnToTable(
    projectId: string,
    tableName: string,
    data: {
      column_name: string
      column_type: string
      is_nullable?: boolean
      default_value?: string
    }
  ) {
    return this.request<{ status: string }>(`/projects/${projectId}/tables/${tableName}/columns`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async removeColumnFromTable(projectId: string, tableName: string, columnName: string) {
    return this.request<void>(`/projects/${projectId}/tables/${tableName}/columns/${columnName}`, {
      method: 'DELETE',
    })
  }

  async modifyColumnInTable(
    projectId: string,
    tableName: string,
    data: {
      column_name: string
      column_type?: string
      is_nullable?: boolean
      default_value?: string
    }
  ) {
    return this.request<{ status: string }>(
      `/projects/${projectId}/tables/${tableName}/columns/${data.column_name}`,
      {
        method: 'PUT',
        body: JSON.stringify(data),
      }
    )
  }

  // Database mutation endpoints
  async getTableData(projectId: string, tableName: string, page = 1, limit = 50) {
    return this.request<{ data: any[]; total: number; page: number; limit: number }>(
      `/projects/${projectId}/tables/${tableName}/data?page=${page}&limit=${limit}`
    )
  }

  async insertTableData(projectId: string, tableName: string, data: Record<string, any>) {
    return this.request<{ status: string }>(`/projects/${projectId}/tables/${tableName}/data`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateTableData(
    projectId: string,
    tableName: string,
    id: string,
    data: Record<string, any>
  ) {
    return this.request<{ status: string }>(
      `/projects/${projectId}/tables/${tableName}/data/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify(data),
      }
    )
  }

  async deleteTableData(projectId: string, tableName: string, id: string) {
    return this.request<void>(`/projects/${projectId}/tables/${tableName}/data/${id}`, {
      method: 'DELETE',
    })
  }

  // Workflow endpoints
  async getWorkflows(projectId: string) {
    return this.request<Workflow[]>(`/projects/${projectId}/workflows`)
  }

  async getWorkflow(workflowId: string) {
    return this.request<Workflow>(`/workflows/${workflowId}`)
  }

  async createWorkflow(
    projectId: string,
    data: { name: string; description?: string; trigger: any; actions: any[] }
  ) {
    return this.request<Workflow>(`/projects/${projectId}/workflows`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateWorkflow(
    workflowId: string,
    data: { name: string; description?: string; trigger: any; actions: any[] }
  ) {
    return this.request<Workflow>(`/workflows/${workflowId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteWorkflow(workflowId: string) {
    return this.request<void>(`/workflows/${workflowId}`, {
      method: 'DELETE',
    })
  }

  // Activity log endpoints
  async getActivityLogs(projectId: string, page = 1, limit = 50) {
    return this.request<{ logs: ActivityLog[]; total: number; page: number; limit: number }>(
      `/projects/${projectId}/activity?page=${page}&limit=${limit}`
    )
  }
}

export const apiClient = new ApiClient(API_BASE_URL)
