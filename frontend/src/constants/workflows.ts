import type { Component } from 'vue'
import { Database, Globe, Mail, Settings } from 'lucide-vue-next'

export interface TypeOption {
  value: string
  label: string
  icon: Component
}

export const WORKFLOW_TRIGGER_TYPES: TypeOption[] = [
  { value: 'on_row_created', label: 'On Row Created', icon: Database },
  { value: 'on_row_updated', label: 'On Row Updated', icon: Database },
  { value: 'on_row_deleted', label: 'On Row Deleted', icon: Database },
  { value: 'scheduled', label: 'Scheduled', icon: Settings },
  { value: 'webhook', label: 'Webhook', icon: Globe },
]

export const WORKFLOW_ACTION_TYPES: TypeOption[] = [
  { value: 'send_webhook', label: 'Send Webhook', icon: Globe },
  { value: 'send_email', label: 'Send Email', icon: Mail },
  { value: 'update_row', label: 'Update Row', icon: Database },
  { value: 'create_row', label: 'Create Row', icon: Database },
  { value: 'delete_row', label: 'Delete Row', icon: Database },
]

export const DEFAULT_WORKFLOW_ACTION_CONFIG = {
  send_webhook: {
    url: '',
    method: 'POST',
    payload: {},
  },
  send_email: {
    to: '',
    subject: '',
    template: '',
  },
  update_row: {
    table: '',
    condition: '',
    updates: {},
  },
  create_row: {
    table: '',
    data: {},
  },
  delete_row: {
    table: '',
    condition: '',
  },
}
