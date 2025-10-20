import { computed } from 'vue'

import { workflowService } from '@/api/services/workflow'
import type { Workflow, WorkflowCreateRequest, WorkflowUpdateRequest } from '@/types/api'

import { useResourceCache } from '../data/useResourceCache'

/**
 * Composable for managing workflows data for a specific project
 * Use this instead of direct API calls in components
 */
export function useWorkflows(projectId: string) {
  const cacheKey = `workflows:${projectId}`

  const {
    data: workflows,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Workflow[]>(cacheKey, {
    fetchFn: () => workflowService.getAll(projectId),
    subscribeToUpdates: true,
    scope: 'project',
    resourceId: projectId,
    events: ['workflow:created', 'workflow:updated', 'workflow:deleted'],
    ttl: 60000, // 1 minute
    refetchOnMount: true,
    optimisticUpdate: true,
  })

  const createWorkflow = async (data: WorkflowCreateRequest) => {
    const newWorkflow = await workflowService.create(projectId, data)
    refresh() // Refresh to get updated list
    return newWorkflow
  }

  const updateWorkflow = async (workflowId: string, data: WorkflowUpdateRequest) => {
    const updated = await workflowService.update(projectId, workflowId, data)
    refresh() // Refresh to get updated list
    return updated
  }

  const deleteWorkflow = async (workflowId: string) => {
    await workflowService.delete(projectId, workflowId)
    refresh() // Refresh to get updated list
  }

  const getWorkflowById = async (workflowId: string) => {
    return await workflowService.getById(projectId, workflowId)
  }

  return {
    workflows: computed(() => workflows.value ?? []),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    invalidate,
    createWorkflow,
    updateWorkflow,
    deleteWorkflow,
    getWorkflowById,
  }
}
