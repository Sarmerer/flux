import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/projects'

// Lazy load routes
const Home = () => import('@/pages/Home.vue')
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
  { path: '/', name: 'Home', component: Home },
  { path: '/login', name: 'Login', component: Login, meta: { requiresGuest: true } },
  { path: '/register', name: 'Register', component: Register, meta: { requiresGuest: true } },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects',
    name: 'Projects',
    component: Projects,
    meta: { requiresAuth: true },
  },
  {
    path: '/projects/:id',
    name: 'ProjectDetail',
    redirect: { name: 'ProjectOverview' },
  },
  {
    path: '/projects/:id/overview',
    name: 'ProjectOverview',
    component: ProjectDetail,
    meta: { requiresAuth: true, tab: 'overview' },
  },
  {
    path: '/projects/:id/overview/tables',
    name: 'ProjectTablesView',
    component: ProjectDetail,
    meta: { requiresAuth: true, tab: 'tables' },
  },
  {
    path: '/projects/:id/overview/workflows',
    name: 'ProjectWorkflowsView',
    component: ProjectDetail,
    meta: { requiresAuth: true, tab: 'workflows' },
  },
  {
    path: '/projects/:id/overview/activity',
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
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: { requiresAuth: true },
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guards
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const projectStore = useProjectStore()

  // Auth guard
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }

  // Project context restoration for project-specific routes
  const projectIdFromRoute = (to.params.projectId as string) || (to.params.id as string)

  if (projectIdFromRoute) {
    // If navigating to a project route and no current project or different project
    if (!projectStore.currentProject || projectStore.currentProject.id !== projectIdFromRoute) {
      try {
        // Try to load the project (will use cache if available)
        await projectStore.loadProjectById(projectIdFromRoute)
      } catch (error) {
        console.error('Failed to load project context:', error)
        // Continue navigation anyway - the page will handle the error
      }
    }
  }

  next()
})
