export interface Table {
  id: string
  name: string
  description?: string
  project_id: string
  schema: string
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
  is_identity?: boolean
  unique?: boolean
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

export interface TableDetail {
  id: string
  name: string
  description?: string
  columns: TableColumn[]
  rows: Record<string, any>[]
  total?: number
}

export type TableData = Record<string, unknown>

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
