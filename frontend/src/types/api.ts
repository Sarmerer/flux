export interface User {
  id: string
  email: string
  name: string
  created_at: string
  updated_at: string
}

export interface Project {
  id: string
  name: string
  description?: string
  owner_id: string
  created_at: string
  updated_at: string
}

export interface Table {
  id: string
  name: string
  description?: string
  project_id: string
  created_at: string
  updated_at: string
}

export interface Database {
  id: string
  name: string
  host: string
  port: number
  username: string
  database: string
  project_id: string
  created_at: string
  updated_at: string
}

export interface TableColumn {
  name: string
  type: string
  is_nullable: boolean
  default_value?: string
  is_primary_key?: boolean
  is_foreign_key?: boolean
  foreign_table?: string
  foreign_column?: string
}

export interface TableSchema {
  columns: TableColumn[]
  primary_keys: string[]
  foreign_keys: ForeignKey[]
}

export interface ForeignKey {
  column: string
  referenced_table: string
  referenced_column: string
}

export interface DatabaseMutation {
  id: string
  type: 'create_table' | 'drop_table' | 'add_column' | 'remove_column' | 'modify_column' | 'add_foreign_key' | 'remove_foreign_key'
  table_name: string
  data: any
  status: 'pending' | 'in_progress' | 'completed' | 'failed'
  created_at: string
  completed_at?: string
  error_message?: string
}

export interface Workflow {
  id: string
  name: string
  description?: string
  project_id: string
  trigger: WorkflowTrigger
  actions: WorkflowAction[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface WorkflowTrigger {
  type: 'on_row_created' | 'on_row_updated' | 'on_row_deleted' | 'scheduled' | 'webhook'
  table_name?: string
  conditions?: Record<string, any>
  schedule?: string
}

export interface WorkflowAction {
  id: string
  type: 'send_webhook' | 'send_email' | 'update_row' | 'create_row' | 'delete_row'
  config: Record<string, any>
  order: number
}

export interface ActivityLog {
  id: string
  project_id: string
  type: 'table_created' | 'table_updated' | 'table_deleted' | 'row_created' | 'row_updated' | 'row_deleted' | 'workflow_triggered' | 'workflow_completed' | 'workflow_failed'
  entity_type: 'table' | 'row' | 'workflow' | 'database'
  entity_id: string
  entity_name: string
  details: Record<string, any>
  user_id?: string
  created_at: string
}

export interface TableData {
  [key: string]: any
}

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
