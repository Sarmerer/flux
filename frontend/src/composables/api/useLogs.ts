import { computed } from 'vue'

import { logsService, type LogQueryParams, type Log } from '@/api/services/logs'

import { useResourceCache } from '../data/useResourceCache'

export function useLogs(params?: LogQueryParams) {
  const cacheKey = `logs:${JSON.stringify(params || {})}`

  const {
    data: logs,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Log[]>(cacheKey, {
    ttlMs: 30000,
    tags: ['logs'],
    fetch: {
      fn: async () => {
        const response = await logsService.getAll(params)
        return response.logs
      },
      onMount: true,
    },
    toastOnError: {
      title: 'Failed to load logs',
      description: 'Could not fetch system logs',
    },
  })

  return {
    logs: computed(() => logs.value ?? []),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    invalidate,
  }
}

export function useProjectLogs(projectId: string, params?: Omit<LogQueryParams, 'project_id'>) {
  const cacheKey = `project-logs:${projectId}:${JSON.stringify(params || {})}`

  const {
    data: logs,
    loading,
    error,
    refresh,
    invalidate,
  } = useResourceCache<Log[]>(cacheKey, {
    ttlMs: 30000,
    tags: ['logs', `project:${projectId}`],
    fetch: {
      fn: async () => {
        const response = await logsService.getProjectLogs(projectId, params)
        return response.logs
      },
      onMount: true,
    },
    realtime: {
      resourceId: projectId,
      events: ['log:created'],
    },
    toastOnError: {
      title: 'Failed to load project logs',
      description: 'Could not fetch project activity logs',
    },
  })

  return {
    logs: computed(() => logs.value ?? []),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    invalidate,
  }
}
