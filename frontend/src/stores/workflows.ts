import { computed, ref } from 'vue'

import { workflowService } from '@/api/services/workflow'
import type { Workflow } from '@/types/api'
import { defineStore } from 'pinia'

export const useWorkflowStore = defineStore('workflows', () => {
  const workflows = ref<Workflow[]>([])
  const isLoading = ref(false)

  const length = computed(() => workflows.value.length)

  const loadByProjectId = async (projectId: string) => {
    isLoading.value = true
    try {
      workflows.value = await workflowService.getAll(projectId)
      return workflows
    } catch (error) {
      console.error('Failed to load workflows:', error)
      workflows.value = []
      throw error
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
    loadByProjectId,
    clear,
  }
})
