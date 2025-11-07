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
  | { type: 'on_row_created'; table_id: string; conditions?: Record<string, unknown> }
  | { type: 'on_row_updated'; table_id: string; conditions?: Record<string, unknown> }
  | { type: 'on_row_deleted'; table_id: string; conditions?: Record<string, unknown> }
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
