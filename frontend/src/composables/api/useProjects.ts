import { computed } from 'vue'

import { projectService } from '@/api/services/project'
import type { Project } from '@/types/api'

import { useResourceCache } from '../data/useResourceCache'

/**
 * Composable for managing projects data
 * Use this instead of direct API calls in components
 */
export function useProjects() {
  const {
    data: projects,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Project[]>('projects', {
    fetchFn: () => projectService.getAll(),
    subscribeToUpdates: true,
    scope: 'global',
    events: ['project:created', 'project:updated', 'project:deleted'],
    ttl: 60000, // 1 minute
    refetchOnMount: true,
    optimisticUpdate: true,
  })

  const createProject = async (data: { name: string; description?: string }) => {
    const newProject = await projectService.create(data)
    refresh() // Refresh to get updated list
    return newProject
  }

  const updateProject = async (id: string, data: { name?: string; description?: string }) => {
    const updated = await projectService.update(id, data)
    refresh() // Refresh to get updated list
    return updated
  }

  const deleteProject = async (id: string) => {
    await projectService.delete(id)
    refresh() // Refresh to get updated list
  }

  const getProjectById = async (id: string) => {
    return await projectService.getById(id)
  }

  return {
    projects: computed(() => projects.value ?? []),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    invalidate,
    createProject,
    updateProject,
    deleteProject,
    getProjectById,
  }
}
