import { createRouter, createWebHistory } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useAuthStore } from '@/stores/auth'
import { useProjectMemberStore } from '@/stores/projectMember'
import type { Permission, Role } from '@/types/auth'

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

const Workflows = () => import('@/pages/workflow/Builder.vue')

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
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/tables',
    name: 'Tables',
    component: Tables,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/tables/:tableId',
    name: 'TableDetail',
    component: Tables,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/workflows',
    name: 'Workflows',
    component: Workflows,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:projectId/workflows/:workflowId',
    name: 'WorkflowDetail',
    component: Workflows,
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
  const projectMemberStore = useProjectMemberStore()

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
      try {
        await activeProjectStore.loadById(projectIdFromRoute)
        await projectMemberStore.loadMyProjectRole(projectIdFromRoute)
      } catch (error) {
        console.error('Failed to load project context:', error)
        next({ name: 'Dashboard' })
        return
      }
    } else if (!projectIdFromRoute) {
      activeProjectStore.clear()
      projectMemberStore.clearCurrentProject()
    }

    if (to.meta.permissions && projectIdFromRoute) {
      if (!projectMemberStore.hasPermission(to.meta.permissions)) {
        next({ name: 'Dashboard' })
        return
      }
    }

    if (to.meta.roles && projectIdFromRoute) {
      if (!projectMemberStore.hasRole(to.meta.roles)) {
        next({ name: 'Dashboard' })
        return
      }
    }
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/')
    return
  }

  next()
})
