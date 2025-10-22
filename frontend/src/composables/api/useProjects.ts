import { computed } from 'vue'

import { projectService } from '@/api/services/project'
import { useCacheStore } from '@/stores/cache'
import type { Project } from '@/types/api'

import { useResourceCache } from '../data/useResourceCache'

/**
 * Composable for managing projects data
 * Use this instead of direct API calls in components
 */
export function useProjects() {
  const cacheStore = useCacheStore()

  const {
    data: projects,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Project[]>('projects', {
    fetchFn: () => projectService.getAll(),
    subscribeToUpdates: true,
    events: ['project:created', 'project:updated', 'project:deleted'],
    ttlMs: 60000,
    fetchOnMount: true,
    tags: ['projects', 'global'],
  })

  const createProject = async (data: { name: string; description?: string }) => {
    const newProject = await projectService.create(data)
    refresh()

    cacheStore.invalidateByTag('projects')
    return newProject
  }

  const updateProject = async (id: string, data: { name?: string; description?: string }) => {
    const updated = await projectService.update(id, data)
    refresh()

    cacheStore.invalidate(`project:${id}`)
    cacheStore.invalidateByTag('projects')
    return updated
  }

  const deleteProject = async (id: string) => {
    await projectService.delete(id)
    refresh()

    cacheStore.invalidateByTag(`project:${id}`)
    cacheStore.invalidateByTag('projects')
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
