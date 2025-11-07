import type { WorkflowAction, WorkflowTrigger } from '@/types'

export interface ValidationResult {
  valid: boolean
  message?: string
}

export function useWorkflowValidation() {
  const validateWorkflowName = (name: string): ValidationResult => {
    if (!name.trim()) {
      return { valid: false, message: 'Workflow name is required' }
    }
    if (name.length > 255) {
      return { valid: false, message: 'Workflow name must be less than 255 characters' }
    }
    return { valid: true }
  }

  const validateActions = (actions: WorkflowAction[]): ValidationResult => {
    if (actions.length === 0) {
      return { valid: false, message: 'At least one action is required' }
    }
    return { valid: true }
  }

  const validateTrigger = (trigger: WorkflowTrigger): ValidationResult => {
    if (trigger.type.startsWith('on_row_') && 'table_id' in trigger && !trigger.table_id) {
      return { valid: false, message: 'Table selection is required for row-based triggers' }
    }

    if (trigger.type === 'scheduled' && 'schedule' in trigger && !trigger.schedule) {
      return { valid: false, message: 'Schedule is required for scheduled triggers' }
    }

    if (trigger.type === 'webhook' && 'webhook_url' in trigger && !trigger.webhook_url) {
      return { valid: false, message: 'Webhook URL is required for webhook triggers' }
    }

    return { valid: true }
  }

  const validateWorkflow = (
    name: string,
    trigger: WorkflowTrigger,
    actions: WorkflowAction[]
  ): ValidationResult => {
    const nameValidation = validateWorkflowName(name)
    if (!nameValidation.valid) return nameValidation

    const triggerValidation = validateTrigger(trigger)
    if (!triggerValidation.valid) return triggerValidation

    const actionsValidation = validateActions(actions)
    if (!actionsValidation.valid) return actionsValidation

    return { valid: true }
  }

  return {
    validateWorkflowName,
    validateActions,
    validateTrigger,
    validateWorkflow,
  }
}
