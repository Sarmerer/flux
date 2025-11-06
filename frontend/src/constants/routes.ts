export const ROUTE_PATHS = {
  HOME: '/',
  PROJECTS: '/projects',
  SETTINGS: '/settings',
  PROJECT_DETAIL: (projectId: string) => `/projects/${projectId}`,
  PROJECT_TABLES: (projectId: string) => `/projects/${projectId}/tables`,
  PROJECT_WORKFLOWS: (projectId: string) => `/projects/${projectId}/workflows`,
  PROJECT_MEMBERS: (projectId: string) => `/projects/${projectId}/members`,
  PROJECT_ACTIVITY: (projectId: string) => `/projects/${projectId}/activity`,
} as const

export const ROUTE_SEGMENTS = {
  TABLES: 'tables',
  WORKFLOWS: 'workflows',
  MEMBERS: 'members',
  ACTIVITY: 'activity',
} as const
