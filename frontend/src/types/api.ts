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
  columns?: number
  rows?: number
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

export interface BaseDatabaseMutation {
  id: string
  table_name: string
  status: 'pending' | 'in_progress' | 'completed' | 'failed'
  created_at: string
  completed_at?: string
  error_message?: string
}

export type DatabaseMutation =
  | (BaseDatabaseMutation & { type: 'create_table'; data: { columns: TableColumn[] } })
  | (BaseDatabaseMutation & { type: 'drop_table'; data: Record<string, never> })
  | (BaseDatabaseMutation & { type: 'add_column'; data: { column: TableColumn } })
  | (BaseDatabaseMutation & { type: 'remove_column'; data: { column_name: string } })
  | (BaseDatabaseMutation & { type: 'modify_column'; data: { column: TableColumn } })
  | (BaseDatabaseMutation & { type: 'add_foreign_key'; data: { foreign_key: ForeignKey } })
  | (BaseDatabaseMutation & { type: 'remove_foreign_key'; data: { column_name: string } })

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

export type WorkflowTrigger =
  | { type: 'on_row_created'; table_name: string; conditions?: Record<string, unknown> }
  | { type: 'on_row_updated'; table_name: string; conditions?: Record<string, unknown> }
  | { type: 'on_row_deleted'; table_name: string; conditions?: Record<string, unknown> }
  | { type: 'scheduled'; schedule: string }
  | { type: 'webhook'; webhook_url?: string }

export interface BaseWorkflowAction {
  id: string
  order: number
}

export type WorkflowAction =
  | (BaseWorkflowAction & {
      type: 'send_webhook'
      config: {
        url: string
        method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
        payload?: Record<string, unknown>
        headers?: Record<string, string>
      }
    })
  | (BaseWorkflowAction & {
      type: 'send_email'
      config: {
        to: string
        subject: string
        template: string
        variables?: Record<string, unknown>
      }
    })
  | (BaseWorkflowAction & {
      type: 'update_row'
      config: {
        table: string
        condition: string
        updates: Record<string, unknown>
      }
    })
  | (BaseWorkflowAction & {
      type: 'create_row'
      config: {
        table: string
        values: Record<string, unknown>
      }
    })
  | (BaseWorkflowAction & {
      type: 'delete_row'
      config: {
        table: string
        condition: string
      }
    })

export interface BaseActivityLog {
  id: string
  project_id: string
  entity_type: 'table' | 'row' | 'workflow' | 'database'
  entity_id: string
  entity_name: string
  user_id?: string
  created_at: string
}

export type ActivityLog =
  | (BaseActivityLog & {
      type: 'table_created'
      details: { columns: number; description?: string }
    })
  | (BaseActivityLog & { type: 'table_updated'; details: { changes: string[] } })
  | (BaseActivityLog & { type: 'table_deleted'; details: Record<string, never> })
  | (BaseActivityLog & { type: 'row_created'; details: { table: string; row_id: string } })
  | (BaseActivityLog & {
      type: 'row_updated'
      details: { table: string; row_id: string; fields: string[] }
    })
  | (BaseActivityLog & { type: 'row_deleted'; details: { table: string; row_id: string } })
  | (BaseActivityLog & { type: 'workflow_triggered'; details: { trigger_type: string } })
  | (BaseActivityLog & {
      type: 'workflow_completed'
      details: { duration_ms: number; actions_executed: number }
    })
  | (BaseActivityLog & {
      type: 'workflow_failed'
      details: { error: string; failed_action?: string }
    })

export type TableData = Record<string, unknown>

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

export interface WorkflowCreateRequest {
  name: string
  description?: string
  trigger: WorkflowTrigger
  actions: WorkflowAction[]
  is_active?: boolean
}

export interface WorkflowUpdateRequest {
  name?: string
  description?: string
  trigger?: WorkflowTrigger
  actions?: WorkflowAction[]
  is_active?: boolean
}
