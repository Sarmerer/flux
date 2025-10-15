<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectsStore } from '@/stores/projects'
import { useAuthStore } from '@/stores/auth'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { 
  FolderOpen, 
  Plus, 
  Activity, 
  Database, 
  Table, 
  Workflow,
  TrendingUp,
  Clock
} from 'lucide-vue-next'

const router = useRouter()
const projectsStore = useProjectsStore()
const authStore = useAuthStore()

const recentActivity = ref([
  {
    id: '1',
    type: 'project_created',
    message: 'Created project "E-commerce Platform"',
    timestamp: '2 hours ago',
    icon: FolderOpen
  },
  {
    id: '2',
    type: 'table_created',
    message: 'Added table "users" to E-commerce Platform',
    timestamp: '4 hours ago',
    icon: Table
  },
  {
    id: '3',
    type: 'workflow_triggered',
    message: 'Workflow "New User Welcome" executed',
    timestamp: '6 hours ago',
    icon: Workflow
  }
])

const stats = computed(() => ({
  totalProjects: projectsStore.projects.length,
  totalTables: projectsStore.projects.reduce((acc, project) => acc + (project as any).table_count || 0, 0),
  totalWorkflows: projectsStore.projects.reduce((acc, project) => acc + (project as any).workflow_count || 0, 0),
  activeConnections: projectsStore.projects.reduce((acc, project) => acc + (project as any).database_count || 0, 0)
}))

onMounted(async () => {
  try {
    await projectsStore.fetchProjects()
  } catch (error) {
    console.error('Failed to fetch projects:', error)
  }
})

const handleCreateProject = () => {
  router.push('/projects')
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-gray-600">Welcome back, {{ authStore.user?.name || 'User' }}!</p>
      </div>
      <Button @click="handleCreateProject" class="flex items-center space-x-2">
        <Plus class="w-4 h-4" />
        <span>New Project</span>
      </Button>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Projects</CardTitle>
          <FolderOpen class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.totalProjects }}</div>
          <p class="text-xs text-muted-foreground">
            +2 from last month
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Tables</CardTitle>
          <Table class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.totalTables }}</div>
          <p class="text-xs text-muted-foreground">
            +12 from last month
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Active Workflows</CardTitle>
          <Workflow class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.totalWorkflows }}</div>
          <p class="text-xs text-muted-foreground">
            +3 from last month
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Database Connections</CardTitle>
          <Database class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.activeConnections }}</div>
          <p class="text-xs text-muted-foreground">
            All systems operational
          </p>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent Projects -->
      <Card>
        <CardHeader>
          <CardTitle>Recent Projects</CardTitle>
          <CardDescription>Your most recently accessed projects</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="projectsStore.isLoading" class="flex items-center justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
          <div v-else-if="projectsStore.projects.length === 0" class="text-center py-8">
            <FolderOpen class="mx-auto h-12 w-12 text-gray-400" />
            <h3 class="mt-2 text-sm font-medium text-gray-900">No projects</h3>
            <p class="mt-1 text-sm text-gray-500">Get started by creating a new project.</p>
            <div class="mt-6">
              <Button @click="handleCreateProject">Create Project</Button>
            </div>
          </div>
          <div v-else class="space-y-3">
            <div 
              v-for="project in projectsStore.projects.slice(0, 5)" 
              :key="project.id"
              class="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50 cursor-pointer"
              @click="router.push(`/projects/${project.id}`)"
            >
              <div class="flex items-center space-x-3">
                <FolderOpen class="h-5 w-5 text-blue-600" />
                <div>
                  <h4 class="text-sm font-medium">{{ project.name }}</h4>
                  <p class="text-xs text-gray-500">{{ project.description || 'No description' }}</p>
                </div>
              </div>
              <Badge variant="secondary">Active</Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Recent Activity -->
      <Card>
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>Latest updates across your projects</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div 
              v-for="activity in recentActivity" 
              :key="activity.id"
              class="flex items-start space-x-3"
            >
              <div class="flex-shrink-0">
                <component :is="activity.icon" class="h-5 w-5 text-blue-600" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-900">{{ activity.message }}</p>
                <p class="text-xs text-gray-500 flex items-center space-x-1">
                  <Clock class="h-3 w-3" />
                  <span>{{ activity.timestamp }}</span>
                </p>
              </div>
            </div>
          </div>
          <div class="mt-4">
            <Button variant="outline" size="sm" class="w-full">
              View All Activity
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
