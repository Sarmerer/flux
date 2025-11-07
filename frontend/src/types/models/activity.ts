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
