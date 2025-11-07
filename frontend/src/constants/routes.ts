const PROJECT_BASE = '/projects/:projectId'

export const ROUTE_PATHS = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  SETTINGS: '/settings',
  PROJECTS: '/projects',
  PROJECTS_NEW: '/projects/new',
  PROJECT: PROJECT_BASE,
  PROJECT_TABLES: `${PROJECT_BASE}/tables`,
  PROJECT_TABLE_DETAIL: `${PROJECT_BASE}/tables/:tableId`,
  PROJECT_WORKFLOWS: `${PROJECT_BASE}/workflows`,
  PROJECT_WORKFLOW_DETAIL: `${PROJECT_BASE}/workflows/:workflowId`,
  PROJECT_MEMBERS: `${PROJECT_BASE}/members`,
  PROJECT_ACTIVITY: `${PROJECT_BASE}/activity`,
} as const

export const buildPath = {
  project: (projectId: string) => `/projects/${projectId}`,
  projectTables: (projectId: string) => `/projects/${projectId}/tables`,
  projectTableDetail: (projectId: string, tableId: string) => `/projects/${projectId}/tables/${tableId}`,
  projectWorkflows: (projectId: string) => `/projects/${projectId}/workflows`,
  projectWorkflowDetail: (projectId: string, workflowId: string) => `/projects/${projectId}/workflows/${workflowId}`,
  projectMembers: (projectId: string) => `/projects/${projectId}/members`,
  projectActivity: (projectId: string) => `/projects/${projectId}/activity`,
} as const

export const ROUTE_NAMES = {
  HOME: 'Dashboard',
  LOGIN: 'Login',
  REGISTER: 'Register',
  SETTINGS: 'Settings',
  PROJECTS: 'Projects',
  PROJECTS_NEW: 'ProjectsNew',
  PROJECT: 'Project',
  PROJECT_TABLES: 'Tables',
  PROJECT_TABLE_DETAIL: 'TableDetail',
  PROJECT_WORKFLOWS: 'Workflows',
  PROJECT_WORKFLOW_DETAIL: 'WorkflowDetail',
  PROJECT_MEMBERS: 'ProjectMembers',
  PROJECT_ACTIVITY: 'Activity',
} as const

export const ROUTE_SEGMENTS = {
  TABLES: 'tables',
  WORKFLOWS: 'workflows',
  MEMBERS: 'members',
  ACTIVITY: 'activity',
} as const
