import { computed } from 'vue'

import { projectMemberService } from '@/api/services/projectMember'
import { useCacheStore } from '@/stores/cache'
import type { ProjectMemberWithUser, Permission, Role } from '@/types'

import { useResourceCache } from '../data/useResourceCache'

export function useProjectMembers(projectId: string) {
  const cacheKey = `project-members:${projectId}`
  const cacheStore = useCacheStore()

  const {
    data: members,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<ProjectMemberWithUser[]>(cacheKey, {
    ttlMs: 60000,
    tags: ['project-members', `project:${projectId}`],
    fetch: {
      fn: () => projectMemberService.getProjectMembers(projectId),
      onMount: true,
    },
    realtime: {
      resourceId: projectId,
      events: ['member:added', 'member:updated', 'member:removed'],
    },
    toastOnError: {
      title: 'Failed to load members',
      description: 'Could not fetch project members',
    },
  })

  const addMember = async (data: { email: string; role: Role; permissions: Permission[] }) => {
    const newMember = await projectMemberService.addProjectMember(projectId, data)
    refresh()
    cacheStore.invalidateByTag(`project:${projectId}`)
    return newMember
  }

  const updateMember = async (
    memberId: string,
    data: { role: Role; permissions: Permission[] }
  ) => {
    const updated = await projectMemberService.updateProjectMember(projectId, memberId, data)
    refresh()
    cacheStore.invalidateByTag(`project:${projectId}`)
    return updated
  }

  const removeMember = async (memberId: string) => {
    await projectMemberService.removeProjectMember(projectId, memberId)
    refresh()
    cacheStore.invalidateByTag(`project:${projectId}`)
  }

  return {
    members: computed(() => members.value ?? []),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    invalidate,
    addMember,
    updateMember,
    removeMember,
  }
}
