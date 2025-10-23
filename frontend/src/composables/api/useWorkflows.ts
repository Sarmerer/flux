import { computed } from 'vue'

import { workflowService } from '@/api/services/workflow'
import { useCacheStore } from '@/stores/cache'
import type { Workflow, WorkflowCreateRequest, WorkflowUpdateRequest } from '@/types/api'

import { useResourceCache } from '../data/useResourceCache'

export function useWorkflows(projectId: string) {
  const cacheKey = `workflows:${projectId}`
  const cacheStore = useCacheStore()

  const {
    data: workflows,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Workflow[]>(cacheKey, {
    ttlMs: 60000,
    tags: ['workflows', `project:${projectId}`],
    fetch: {
      fn: () => workflowService.getAll(projectId),
      onMount: true,
    },
    realtime: {
      resourceId: projectId,
      events: ['workflow:created', 'workflow:updated', 'workflow:deleted'],
    },
  })

  const createWorkflow = async (data: WorkflowCreateRequest) => {
    const newWorkflow = await workflowService.create(projectId, data)
    refresh()

    cacheStore.invalidateByTag('workflows')
    return newWorkflow
  }

  const updateWorkflow = async (workflowId: string, data: WorkflowUpdateRequest) => {
    const updated = await workflowService.update(projectId, workflowId, data)
    refresh()

    cacheStore.invalidate(`workflow:${workflowId}`)
    return updated
  }

  const deleteWorkflow = async (workflowId: string) => {
    await workflowService.delete(projectId, workflowId)
    refresh()

    cacheStore.invalidate(`workflow:${workflowId}`)
    cacheStore.invalidateByTag(`project:${projectId}`)
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
