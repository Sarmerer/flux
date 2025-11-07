import { createRouter, createWebHistory } from 'vue-router'

import { ROUTE_NAMES, ROUTE_PATHS } from '@/constants/routes'
import { useActiveProjectStore } from '@/stores/activeProject'
import { useAuthStore } from '@/stores/auth'
import { useProjectMemberStore } from '@/stores/projectMember'
import type { Permission, Role } from '@/types'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    requiresGuest?: boolean
    permissions?: Permission[]
    roles?: Role[]
  }
}

const Login = () => import('@/pages/auth/Login.vue')
const Register = () => import('@/pages/auth/Register.vue')
const Dashboard = () => import('@/pages/Dashboard.vue')
const Activity = () => import('@/pages/Activity.vue')
const Settings = () => import('@/pages/Settings.vue')

const NewProject = () => import('@/pages/project/New.vue')
const ProjectList = () => import('@/pages/project/List.vue')
const ProjectOverview = () => import('@/pages/project/Overview.vue')
const Members = () => import('@/pages/project/Members.vue')

const Tables = () => import('@/pages/table/Builder.vue')

const WorkflowList = () => import('@/pages/workflow/List.vue')
const WorkflowDetail = () => import('@/pages/workflow/Detail.vue')

const routes = [
  {
    path: ROUTE_PATHS.HOME,
    name: ROUTE_NAMES.HOME,
    component: Dashboard,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.LOGIN,
    name: ROUTE_NAMES.LOGIN,
    component: Login,
    meta: { requiresGuest: true },
  },
  {
    path: ROUTE_PATHS.REGISTER,
    name: ROUTE_NAMES.REGISTER,
    component: Register,
    meta: { requiresGuest: true },
  },
  {
    path: ROUTE_PATHS.SETTINGS,
    name: ROUTE_NAMES.SETTINGS,
    component: Settings,
    meta: {
      requiresAuth: true,
      permissions: ['settings.manage' as Permission],
    },
  },

  {
    path: ROUTE_PATHS.PROJECTS,
    name: ROUTE_NAMES.PROJECTS,
    component: ProjectList,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECTS_NEW,
    name: ROUTE_NAMES.PROJECTS_NEW,
    component: NewProject,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECT,
    name: ROUTE_NAMES.PROJECT,
    component: ProjectOverview,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECT_TABLES,
    name: ROUTE_NAMES.PROJECT_TABLES,
    component: Tables,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECT_TABLE_DETAIL,
    name: ROUTE_NAMES.PROJECT_TABLE_DETAIL,
    component: Tables,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECT_WORKFLOWS,
    name: ROUTE_NAMES.PROJECT_WORKFLOWS,
    component: WorkflowList,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECT_WORKFLOW_DETAIL,
    name: ROUTE_NAMES.PROJECT_WORKFLOW_DETAIL,
    component: WorkflowDetail,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECT_ACTIVITY,
    name: ROUTE_NAMES.PROJECT_ACTIVITY,
    component: Activity,
    meta: { requiresAuth: true },
  },
  {
    path: ROUTE_PATHS.PROJECT_MEMBERS,
    name: ROUTE_NAMES.PROJECT_MEMBERS,
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
  const projectMemberStore = useProjectMemberStore()

  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      next({
        path: ROUTE_PATHS.LOGIN,
        query: { redirect: to.fullPath },
      })
      return
    }

    const projectIdFromRoute = to.params.projectId as string
    if (projectIdFromRoute && activeProjectStore.activeProject?.id !== projectIdFromRoute) {
      try {
        await activeProjectStore.loadById(projectIdFromRoute)
        await projectMemberStore.loadMyProjectRole(projectIdFromRoute)
      } catch (error) {
        console.error('Failed to load project context:', error)
        next({ name: ROUTE_NAMES.HOME })
        return
      }
    } else if (!projectIdFromRoute) {
      activeProjectStore.clear()
      projectMemberStore.clearCurrentProject()
    }

    if (to.meta.permissions && projectIdFromRoute) {
      if (!projectMemberStore.hasPermission(to.meta.permissions)) {
        next({ name: ROUTE_NAMES.HOME })
        return
      }
    }

    if (to.meta.roles && projectIdFromRoute) {
      if (!projectMemberStore.hasRole(to.meta.roles)) {
        next({ name: ROUTE_NAMES.HOME })
        return
      }
    }
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next(ROUTE_PATHS.HOME)
    return
  }

  next()
})
