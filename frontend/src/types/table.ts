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

export interface Table {
  id: string
  name: string
  description?: string
  project_id: string
  created_at: string
  updated_at: string
}

export interface TableDetail {
  id: string
  name: string
  description?: string
  columns: TableColumn[]
  rows: Record<string, any>[]
  total?: number
}
