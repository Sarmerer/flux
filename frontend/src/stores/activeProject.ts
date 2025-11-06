import { ref } from 'vue'

import { projectService } from '@/api/services/project'
import type { Project } from '@/types/api'
import { defineStore } from 'pinia'

export const useActiveProjectStore = defineStore('project', () => {
  const activeProject = ref<Project | null>(null)
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  async function loadById(projectId: string): Promise<Project | null> {
    isLoading.value = true
    error.value = null

    try {
      const project = await projectService.getById(projectId)
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
  }

  return {
    activeProject,
    isLoading,
    error,
    loadById,
    clear,
  }
})
