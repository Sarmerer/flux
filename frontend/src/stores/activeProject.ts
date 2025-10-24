import { ref } from 'vue'

import { projectService } from '@/api/services/project'
import type { Project } from '@/types/api'
import { defineStore } from 'pinia'

import { useTableStore } from './tables'
import { useWorkflowStore } from './workflows'

export const useActiveProjectStore = defineStore('project', () => {
  const activeProject = ref<Project | null>(null)
  const isLoading = ref(false)

  async function loadById(projectId: string): Promise<Project> {
    isLoading.value = true

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
    } catch (error) {
      console.error('Failed to load project by ID:', error)
      clear()
      throw error
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
    loadById,
    clear,
  }
})
