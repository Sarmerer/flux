import { computed, ref } from 'vue'

import type { Project } from '@/types/api'
import { defineStore } from 'pinia'

import { apiClient } from '@/lib/api'

export const useProjectsStore = defineStore('projects', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const getProjectById = computed(
    () => (id: string) => projects.value.find((project) => project.id === id)
  )

  const fetchProjects = async () => {
    isLoading.value = true
    error.value = null

    try {
      projects.value = await apiClient.getProjects()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch projects'
      throw err
    } finally {
      projects.value = projects.value || []
      isLoading.value = false
    }
  }

  const fetchProject = async (projectId: string) => {
    isLoading.value = true
    error.value = null

    try {
      const project = await apiClient.getProject(projectId)
      currentProject.value = project

      // Update in projects list if it exists
      const index = projects.value.findIndex((p) => p.id === projectId)
      if (index !== -1) {
        projects.value[index] = project
      } else {
        projects.value.push(project)
      }

      return project
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch project'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const createProject = async (data: { name: string; description?: string }) => {
    isLoading.value = true
    error.value = null

    try {
      const project = await apiClient.createProject(data)
      projects.value.push(project)
      return project
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create project'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const updateProject = async (projectId: string, data: { name: string; description?: string }) => {
    isLoading.value = true
    error.value = null

    try {
      const project = await apiClient.updateProject(projectId, data)

      // Update in projects list
      const index = projects.value.findIndex((p) => p.id === projectId)
      if (index !== -1) {
        projects.value[index] = project
      }

      // Update current project if it's the same
      if (currentProject.value?.id === projectId) {
        currentProject.value = project
      }

      return project
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update project'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const deleteProject = async (projectId: string) => {
    isLoading.value = true
    error.value = null

    try {
      await apiClient.deleteProject(projectId)

      // Remove from projects list
      projects.value = projects.value.filter((p) => p.id !== projectId)

      // Clear current project if it's the same
      if (currentProject.value?.id === projectId) {
        currentProject.value = null
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete project'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const setCurrentProject = (project: Project | null) => {
    currentProject.value = project
  }

  const clearError = () => {
    error.value = null
  }

  return {
    projects,
    currentProject,
    isLoading,
    error,
    getProjectById,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
    setCurrentProject,
    clearError,
  }
})
