import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/projects'
import type { Permission, Role } from '@/types/auth'

// Extend route meta with permission and role fields
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    requiresGuest?: boolean
    permissions?: Permission[]
    roles?: Role[]
    openCreateModal?: boolean
    tab?: string
  }
}

const Login = () => import('@/pages/auth/Login.vue')
const Register = () => import('@/pages/auth/Register.vue')
const Dashboard = () => import('@/pages/Dashboard.vue')
const Projects = () => import('@/pages/Projects.vue')
const ProjectDetail = () => import('@/pages/ProjectDetail.vue')
const Tables = () => import('@/pages/Tables.vue')
const TableBuilder = () => import('@/pages/TableBuilder.vue')
const DataGrid = () => import('@/pages/DataGrid.vue')
const Workflows = () => import('@/pages/Workflows.vue')
const WorkflowBuilder = () => import('@/pages/WorkflowBuilder.vue')
const Activity = () => import('@/pages/Activity.vue')
const Settings = () => import('@/pages/Settings.vue')

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresGuest: true },
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: {
      requiresAuth: true,
      permissions: ['settings.manage' as Permission],
    },
  },

  {
    path: '/projects',
    name: 'Projects',
    component: Projects,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/new',
    name: 'ProjectsNew',
    component: Projects,
    meta: { requiresAuth: true, openCreateModal: true },
  },
  {
    path: '/projects/:projectId',
    name: 'ProjectDetail',
    redirect: { name: 'ProjectOverview' },
  },
  {
    path: '/projects/:projectId/overview',
    name: 'ProjectOverview',
    component: ProjectDetail,
    meta: { requiresAuth: true, tab: 'overview' },
  },
  {
    path: '/projects/:projectId/overview/tables',
    name: 'ProjectTablesView',
    component: ProjectDetail,
    meta: { requiresAuth: true, tab: 'tables' },
  },
  {
    path: '/projects/:projectId/overview/workflows',
    name: 'ProjectWorkflowsView',
    component: ProjectDetail,
    meta: { requiresAuth: true, tab: 'workflows' },
  },
  {
    path: '/projects/:projectId/overview/activity',
    name: 'ProjectActivityView',
    component: ProjectDetail,
    meta: { requiresAuth: true, tab: 'activity' },
  },
  {
    path: '/projects/:projectId/tables',
    name: 'Tables',
    component: Tables,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/tables/:tableId/builder',
    name: 'TableBuilder',
    component: TableBuilder,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/tables/:tableId/data',
    name: 'DataGrid',
    component: DataGrid,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/workflows',
    name: 'Workflows',
    component: Workflows,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/workflows/:workflowId/builder',
    name: 'WorkflowBuilder',
    component: WorkflowBuilder,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/activity',
    name: 'Activity',
    component: Activity,
    meta: { requiresAuth: true },
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const projectStore = useProjectStore()

  // Check authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }

  // Check permissions
  if (to.meta.permissions && to.meta.permissions.length > 0) {
    const hasPermission = authStore.hasPermission(to.meta.permissions)
    if (!hasPermission) {
      console.warn('User does not have required permissions:', to.meta.permissions)
      // Redirect to dashboard or unauthorized page
      next('/')
      return
    }
  }

  // Check roles
  if (to.meta.roles && to.meta.roles.length > 0) {
    const hasRole = authStore.hasRole(to.meta.roles)
    if (!hasRole) {
      console.warn('User does not have required role:', to.meta.roles)
      // Redirect to dashboard or unauthorized page
      next('/')
      return
    }
  }

  // Load project context if needed
  const projectIdFromRoute = to.params.projectId as string

  if (projectIdFromRoute) {
    if (!projectStore.currentProject || projectStore.currentProject.id !== projectIdFromRoute) {
      try {
        await projectStore.loadProjectById(projectIdFromRoute)
      } catch (error) {
        console.error('Failed to load project context:', error)
      }
    }
  } else {
    projectStore.setCurrentProject(null)
  }

  next()
})
