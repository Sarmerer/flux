import type { ActivityLog } from '@/types'

import { http } from '../http-client'

export const activityService = {
  get(projectId: string, page = 1, limit = 50) {
    return http.get<{ logs: ActivityLog[]; total: number; page: number; limit: number }>(
      `/projects/${projectId}/activity?page=${page}&limit=${limit}`
    )
  },
}
