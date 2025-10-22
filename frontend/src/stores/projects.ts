import { ref } from 'vue'

import { projectService } from '@/api/services/project'
import type { Project } from '@/types/api'
import { defineStore } from 'pinia'

export const useProjectStore = defineStore('project', () => {
  const currentProject = ref<Project | null>(null)
  const isLoading = ref(false)

  const setCurrentProject = (project: Project | null) => {
    currentProject.value = project
  }

  const loadProjectById = async (projectId: string) => {
    isLoading.value = true
    try {
      const project = await projectService.getById(projectId)
      setCurrentProject(project)
      return project
    } catch (error) {
      console.error('Failed to load project:', error)
      setCurrentProject(null)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  return {
    currentProject,
    isLoading,
    setCurrentProject,
    loadProjectById,
  }
})
