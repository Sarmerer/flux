import { computed, ref } from 'vue'

import { workflowService } from '@/api/services/workflow'
import type { Workflow } from '@/types/api'
import { defineStore } from 'pinia'

export const useWorkflowStore = defineStore('workflows', () => {
  const workflows = ref<Workflow[]>([])
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  const length = computed(() => workflows.value.length)

  const loadByProjectId = async (projectId: string) => {
    isLoading.value = true
    error.value = null
    try {
      workflows.value = await workflowService.getAll(projectId)
      return workflows
    } catch (_error) {
      console.error('Failed to load workflows:', _error)
      error.value = new Error('Failed to load workflows')
      workflows.value = []
      throw _error
    } finally {
      isLoading.value = false
    }
  }

  const clear = () => {
    workflows.value = []
  }

  return {
    workflows,
    length,
    isLoading,
    error,
    loadByProjectId,
    clear,
  }
})
