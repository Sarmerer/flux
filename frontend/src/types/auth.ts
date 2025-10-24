export type Role = 'admin' | 'manager' | 'developer' | 'viewer'

export type Permission =
  | 'projects.create'
  | 'projects.edit'
  | 'projects.delete'
  | 'projects.view'
  | 'workflows.create'
  | 'workflows.edit'
  | 'workflows.delete'
  | 'tables.create'
  | 'tables.edit'
  | 'tables.delete'
  | 'settings.manage'
  | 'users.manage'

export interface User {
  id: string
  email: string
  name: string
  created_at: string
  updated_at: string
}

export interface ProjectMember {
  id: string
  project_id: string
  user_id: string
  role: Role
  permissions: Permission[]
  created_at: string
  updated_at: string
}

export interface ProjectMemberWithUser {
  id: string
  project_id: string
  user: User
  role: Role
  permissions: Permission[]
  created_at: string
  updated_at: string
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterData {
  email: string
  password: string
  name: string
}

export interface AuthResponse {
  user: User
  token: string
}
