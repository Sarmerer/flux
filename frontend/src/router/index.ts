import { createRouter, createWebHistory } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useAuthStore } from '@/stores/auth'
import type { Permission, Role } from '@/types/auth'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    requiresGuest?: boolean
    permissions?: Permission[]
    roles?: Role[]
    tab?: string
  }
}

const Login = () => import('@/pages/auth/Login.vue')
const Register = () => import('@/pages/auth/Register.vue')
const Dashboard = () => import('@/pages/Dashboard.vue')
const Activity = () => import('@/pages/Activity.vue')
const Settings = () => import('@/pages/Settings.vue')

const NewProject = () => import('@/pages/project/New.vue')
const ProjectList = () => import('@/pages/project/List.vue')
const ProjectOverview = () => import('@/pages/project/overview/index.vue')
const Members = () => import('@/pages/project/Members.vue')

const Tables = () => import('@/pages/project/Tables.vue')
const TableBuilder = () => import('@/pages/table/Builder.vue')
const TableData = () => import('@/pages/table/Data.vue')

const Workflows = () => import('@/pages/workflow/List.vue')
const WorkflowBuilder = () => import('@/pages/workflow/Builder.vue')

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
    component: ProjectList,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/new',
    name: 'ProjectsNew',
    component: NewProject,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId',
    name: 'ProjectDetail',
    redirect: { name: 'ProjectOverview' },
  },
  {
    path: '/projects/:projectId/overview',
    name: 'ProjectOverview',
    component: ProjectOverview,
    meta: { requiresAuth: true, tab: 'overview' },
  },
  {
    path: '/projects/:projectId/overview/tables',
    name: 'ProjectTablesView',
    component: ProjectOverview,
    meta: { requiresAuth: true, tab: 'tables' },
  },
  {
    path: '/projects/:projectId/overview/workflows',
    name: 'ProjectWorkflowsView',
    component: ProjectOverview,
    meta: { requiresAuth: true, tab: 'workflows' },
  },
  {
    path: '/projects/:projectId/overview/activity',
    name: 'ProjectActivityView',
    component: ProjectOverview,
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
    component: TableData,
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
  {
    path: '/projects/:projectId/members',
    name: 'ProjectMembers',
    component: Members,
    meta: { requiresAuth: true },
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const activeProjectStore = useActiveProjectStore()

  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      next({
        path: '/login',
        query: { redirect: to.fullPath },
      })
      return
    }

    const projectIdFromRoute = to.params.projectId as string
    if (projectIdFromRoute && activeProjectStore.activeProject?.id !== projectIdFromRoute) {
      await activeProjectStore.loadById(projectIdFromRoute)
    } else if (!projectIdFromRoute) {
      activeProjectStore.clear()
    }
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/')
    return
  }

  next()
})
