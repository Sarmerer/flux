import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Project } from '@/types/api'
import { projectService } from '@/api/services/project'

/**
 * Simplified project store - only stores current project context
 * Project lists should use useResourceCache in components for better memory management
 */
export const useProjectStore = defineStore('project', () => {
  // Only store the currently selected project (global context)
  const currentProject = ref<Project | null>(null)
  const isLoadingProject = ref(false)

  const setCurrentProject = (project: Project | null) => {
    currentProject.value = project

    // Save to localStorage for persistence across page reloads
    if (project) {
      localStorage.setItem('current-project-id', project.id)
      // Also cache the project data for faster restoration
      localStorage.setItem('current-project-data', JSON.stringify(project))
    } else {
      localStorage.removeItem('current-project-id')
      localStorage.removeItem('current-project-data')
    }
  }

  const clearCurrentProject = () => {
    currentProject.value = null
    localStorage.removeItem('current-project-id')
    localStorage.removeItem('current-project-data')
  }

  const loadProjectById = async (projectId: string) => {
    isLoadingProject.value = true
    try {
      const project = await projectService.getById(projectId)
      setCurrentProject(project)
      return project
    } catch (error) {
      console.error('Failed to load project:', error)
      clearCurrentProject()
      throw error
    } finally {
      isLoadingProject.value = false
    }
  }

  const restoreFromLocalStorage = () => {
    try {
      const savedProjectData = localStorage.getItem('current-project-data')
      if (savedProjectData) {
        const project = JSON.parse(savedProjectData) as Project
        currentProject.value = project
        return project
      }
    } catch (error) {
      console.error('Failed to restore project from localStorage:', error)
      clearCurrentProject()
    }
    return null
  }

  // Try to restore from localStorage on init
  restoreFromLocalStorage()

  return {
    currentProject,
    isLoadingProject,
    setCurrentProject,
    clearCurrentProject,
    loadProjectById,
    restoreFromLocalStorage,
  }
})
