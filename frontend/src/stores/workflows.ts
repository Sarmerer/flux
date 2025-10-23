import { computed, ref } from 'vue'

import { workflowService } from '@/api/services/workflow'
import type { Workflow } from '@/types/api'
import { defineStore } from 'pinia'

export const useWorkflowStore = defineStore('workflows', () => {
  const workflowsByProject = ref<Map<string, Workflow[]>>(new Map())
  const isLoading = ref(false)
  const currentProjectId = ref<string | null>(null)

  const currentWorkflows = computed(() => {
    if (!currentProjectId.value) return []
    return workflowsByProject.value.get(currentProjectId.value) || []
  })

  const workflowCount = computed(() => currentWorkflows.value.length)

  const setCurrentProject = (projectId: string | null) => {
    currentProjectId.value = projectId
  }

  const loadWorkflows = async (projectId: string) => {
    isLoading.value = true
    try {
      const workflows = await workflowService.getAll(projectId)
      workflowsByProject.value.set(projectId, Array.isArray(workflows) ? workflows : [])
      currentProjectId.value = projectId
      return workflows
    } catch (error) {
      console.error('Failed to load workflows:', error)
      workflowsByProject.value.set(projectId, [])
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const addWorkflow = (projectId: string, workflow: Workflow) => {
    const workflows = workflowsByProject.value.get(projectId) || []
    workflowsByProject.value.set(projectId, [...workflows, workflow])
  }

  const updateWorkflow = (projectId: string, workflowId: string, updates: Partial<Workflow>) => {
    const workflows = workflowsByProject.value.get(projectId)
    if (!workflows) return

    const index = workflows.findIndex((w) => w.id === workflowId)
    if (index !== -1) {
      const updatedWorkflow = { ...workflows[index], ...updates } as Workflow
      workflows[index] = updatedWorkflow
      workflowsByProject.value.set(projectId, [...workflows])
    }
  }

  const removeWorkflow = (projectId: string, workflowId: string) => {
    const workflows = workflowsByProject.value.get(projectId)
    if (!workflows) return

    workflowsByProject.value.set(
      projectId,
      workflows.filter((w) => w.id !== workflowId)
    )
  }

  const clearProject = (projectId: string) => {
    workflowsByProject.value.delete(projectId)
  }

  const clearAll = () => {
    workflowsByProject.value.clear()
    currentProjectId.value = null
  }

  return {
    currentWorkflows,
    workflowCount,
    isLoading,
    currentProjectId,
    setCurrentProject,
    loadWorkflows,
    addWorkflow,
    updateWorkflow,
    removeWorkflow,
    clearProject,
    clearAll,
  }
})
