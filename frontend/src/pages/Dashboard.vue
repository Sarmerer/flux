<script setup lang="ts">
import { Clock, FolderOpen, Plus, Table, Workflow } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

import { projectService } from '@/api/services/project'
import { useAuthStore } from '@/stores/auth'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { useResourceCache } from '@/composables/data'

const router = useRouter()
const authStore = useAuthStore()

const {
  data: projects,
  loading,
  error,
} = useResourceCache('dashboard-projects', {
  fetch: { fn: () => projectService.getAll(), onMount: true },
})

const recentActivity = [
  {
    id: '1',
    type: 'project_created',
    message: 'Created project "E-commerce Platform"',
    timestamp: '2 hours ago',
    icon: FolderOpen,
  },
  {
    id: '2',
    type: 'table_created',
    message: 'Added table "users" to E-commerce Platform',
    timestamp: '4 hours ago',
    icon: Table,
  },
  {
    id: '3',
    type: 'workflow_triggered',
    message: 'Workflow "New User Welcome" executed',
    timestamp: '6 hours ago',
    icon: Workflow,
  },
]

const onCreateProjectClick = () => {
  router.push('/projects/new')
}

const onProjectClick = (projectId: string) => {
  router.push(`/projects/${projectId}`)
}
</script>

<template>
  <div class="p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Dashboard</h1>
        <p class="text-gray-600 dark:text-gray-400">
          Welcome back, {{ authStore.user?.name || 'User' }}!
        </p>
      </div>
      <Button @click="onCreateProjectClick" class="flex items-center space-x-2">
        <Plus class="w-4 h-4" />
        <span>New Project</span>
      </Button>
    </div>

    <div class="flex gap-6">
      <Card class="grow">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Projects</CardTitle>
          <FolderOpen class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ 0 }}</div>
          <p class="text-xs text-muted-foreground">Your active projects</p>
        </CardContent>
      </Card>

      <Card class="grow">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Tables</CardTitle>
          <Table class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ 0 }}</div>
          <p class="text-xs text-muted-foreground">Across all projects</p>
        </CardContent>
      </Card>

      <Card class="grow">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Active Workflows</CardTitle>
          <Workflow class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ 0 }}</div>
          <p class="text-xs text-muted-foreground">Automation running</p>
        </CardContent>
      </Card>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Recent Projects</CardTitle>
          <CardDescription>Your most recently accessed projects</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="loading" class="flex items-center justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>

          <div v-else-if="error" class="text-center py-8">
            <div class="text-red-600 dark:text-red-400">
              <p class="text-sm font-medium">Failed to load projects</p>
              <p class="text-xs mt-1">{{ error.message }}</p>
            </div>
          </div>

          <div v-else-if="!projects || projects.length === 0" class="text-center py-8">
            <FolderOpen class="mx-auto h-12 w-12 text-gray-400 dark:text-gray-600" />
            <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-gray-100">No projects</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Get started by creating a new project.
            </p>
            <div class="mt-6">
              <Button @click="onCreateProjectClick">Create Project</Button>
            </div>
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="project in projects.slice(0, 5)"
              :key="project.id"
              class="flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors"
              @click="onProjectClick(project.id)"
            >
              <div class="flex items-center space-x-3">
                <FolderOpen class="h-5 w-5 text-blue-600 dark:text-blue-400" />
                <div>
                  <h4 class="text-sm font-medium text-gray-900 dark:text-gray-100">
                    {{ project.name }}
                  </h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400">
                    {{ project.description || 'No description' }}
                  </p>
                </div>
              </div>
              <Badge variant="secondary">Active</Badge>
            </div>

            <div v-if="projects.length > 5" class="pt-2">
              <Button variant="ghost" size="sm" class="w-full" @click="router.push('/projects')">
                View All Projects ({{ projects.length }})
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

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
                <component :is="activity.icon" class="h-5 w-5 text-blue-600 dark:text-blue-400" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-900 dark:text-gray-100">{{ activity.message }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 flex items-center space-x-1">
                  <Clock class="h-3 w-3" />
                  <span>{{ activity.timestamp }}</span>
                </p>
              </div>
            </div>
          </div>
          <div class="mt-4">
            <Button variant="outline" size="sm" class="w-full" @click="router.push('/activity')">
              View All Activity
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
