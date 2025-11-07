<script setup lang="ts">
import ErrorState from '@/components/common/ErrorState.vue'
import LoadingWrapper from '@/components/common/LoadingWrapper.vue'
import { Clock, FolderOpen, Plus, Table, Workflow } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { ROUTE_PATHS, buildPath } from '@/constants/routes'
import { useAuthStore } from '@/stores/auth'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { useProjects } from '@/composables/api/useProjects'

const router = useRouter()
const authStore = useAuthStore()

const { projects, loading, error, refresh } = useProjects()

const totalProjects = computed(() => projects.value?.length || 0)

const totalTables = computed(() => {
  if (!projects.value) return 0
  return projects.value.reduce((sum: number, project: any) => sum + (project.table_count || 0), 0)
})

const activeWorkflows = computed(() => {
  if (!projects.value) return 0
  return projects.value.reduce(
    (sum: number, project: any) => sum + (project.workflow_count || 0),
    0
  )
})

const handleProjectCreate = () => {
  router.push(ROUTE_PATHS.PROJECTS_NEW)
}

const handleProjectClick = (projectId: string) => {
  router.push(buildPath.project(projectId))
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
      <Button @click="handleProjectCreate" class="flex items-center space-x-2">
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
          <div class="text-2xl font-bold">{{ totalProjects }}</div>
          <p class="text-xs text-muted-foreground">Your active projects</p>
        </CardContent>
      </Card>

      <Card class="grow">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Tables</CardTitle>
          <Table class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalTables }}</div>
          <p class="text-xs text-muted-foreground">Across all projects</p>
        </CardContent>
      </Card>

      <Card class="grow">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Active Workflows</CardTitle>
          <Workflow class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ activeWorkflows }}</div>
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
          <LoadingWrapper :is-loading="loading" loading-text="Loading projects...">
            <ErrorState
              v-if="error"
              :error="error"
              title="Failed to load projects"
              :handleRetry="refresh"
            />
            <div v-else class="space-y-3">
              <div
                v-for="project in projects?.slice(0, 5)"
                :key="project.id"
                class="flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors"
                @click="handleProjectClick(project.id)"
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

              <div v-if="projects && projects.length > 5" class="pt-2">
                <Button variant="ghost" size="sm" class="w-full" @click="router.push(ROUTE_PATHS.PROJECTS)">
                  View All Projects ({{ projects.length }})
                </Button>
              </div>
            </div>
          </LoadingWrapper>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>Latest updates across your projects</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="text-center py-8 text-muted-foreground">
            <Clock class="h-12 w-12 mx-auto mb-3 opacity-50" />
            <p class="text-sm">Activity feed coming soon</p>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
