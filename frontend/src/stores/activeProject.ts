import { ref } from 'vue'

import { projectService } from '@/api/services/project'
import type { Project } from '@/types/api'
import { defineStore } from 'pinia'

import { useTableStore } from './tables'
import { useWorkflowStore } from './workflows'

export const useActiveProjectStore = defineStore('project', () => {
  const activeProject = ref<Project | null>(null)
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  async function loadById(projectId: string): Promise<Project | null> {
    isLoading.value = true
    error.value = null

    const tableStore = useTableStore()
    const workflowStore = useWorkflowStore()

    try {
      const [project] = await Promise.all([
        projectService.getById(projectId),
        tableStore.loadByProjectId(projectId),
        workflowStore.loadByProjectId(projectId),
      ])
      activeProject.value = project
      return project
    } catch (_error) {
      console.error('Failed to load project by ID:', _error)
      error.value = new Error('Failed to load active project')
      clear()
      return null
    } finally {
      isLoading.value = false
    }
  }

  function clear() {
    activeProject.value = null

    const tableStore = useTableStore()
    const workflowStore = useWorkflowStore()

    tableStore.clear()
    workflowStore.clear()
  }

  return {
    activeProject,
    isLoading,
    error,
    loadById,
    clear,
  }
})
