import type { WorkflowAction, WorkflowTrigger } from '../models/workflow'

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterData {
  email: string
  password: string
  name: string
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
