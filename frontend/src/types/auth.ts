export enum Role {
  ADMIN = 'admin',
  MANAGER = 'manager',
  DEVELOPER = 'developer',
  VIEWER = 'viewer',
}

// User permissions
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

// User interface with role and permissions
export interface User {
  id: string
  email: string
  name: string
  role: Role
  permissions: Permission[]
  created_at: string
  updated_at: string
}

// Login/Register types
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
