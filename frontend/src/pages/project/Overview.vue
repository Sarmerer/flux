<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import {
  Activity,
  Calendar,
  Settings,
  Table,
  User,
  Workflow,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useTableStore } from '@/stores/tables'
import { useWorkflowStore } from '@/stores/workflows'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { useFormatting } from '@/composables/formatting'

const router = useRouter()

const { formatDate } = useFormatting()

const activeProjectStore = useActiveProjectStore()
const activeProject = computed(() => activeProjectStore.activeProject)
const activeProjectId = computed(() => activeProject.value?.id)

const tablesStore = useTableStore()
const workflowsStore = useWorkflowStore()

const stats = computed(() => ({
  tables: tablesStore.length,
  workflows: workflowsStore.length,
  lastActivity: '2 hours ago',
}))

const onTableCreateClick = () => {
  router.push(`/projects/${activeProjectId.value}/tables`)
}

const onWorkflowCreateClick = () => {
  router.push(`/projects/${activeProjectId.value}/workflows`)
}
</script>

<template>
  <div class="p-6 space-y-6">
    <LoadingWrapper :is-loading="activeProjectStore.isLoading" loading-text="Loading project...">
      <ErrorState
        v-if="activeProjectStore.error"
        :error="activeProjectStore.error"
        title="Project not found"
        description="The project you're looking for doesn't exist."
      >
        <Button size="sm" @click="router.push('/projects')">Back to Projects</Button>
      </ErrorState>

      <template v-else-if="activeProject">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-foreground">{{ activeProject.name }}</h1>
            <p class="text-muted-foreground">{{ activeProject.description || 'No description' }}</p>
          </div>
          <div class="flex items-center space-x-2">
            <Button variant="outline">
              <Settings class="w-4 h-4 mr-2" />
              Settings
            </Button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">Tables</CardTitle>
              <Table class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">{{ stats.tables }}</div>
              <p class="text-xs text-muted-foreground">Database tables</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">Workflows</CardTitle>
              <Workflow class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">{{ stats.workflows }}</div>
              <p class="text-xs text-muted-foreground">Active workflows</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">Last Activity</CardTitle>
              <Activity class="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">{{ stats.lastActivity }}</div>
              <p class="text-xs text-muted-foreground">Time since last update</p>
            </CardContent>
          </Card>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Quick Actions</CardTitle>
              <CardDescription>Get started with your project</CardDescription>
            </CardHeader>
            <CardContent class="space-y-3">
              <Button
                variant="outline"
                class="w-full justify-start"
                @click="onTableCreateClick"
              >
                <Table class="w-4 h-4 mr-2" />
                Create New Table
              </Button>
              <Button
                variant="outline"
                class="w-full justify-start"
                @click="onWorkflowCreateClick"
              >
                <Workflow class="w-4 h-4 mr-2" />
                Create Workflow
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Project Information</CardTitle>
              <CardDescription>Details about this project</CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">Created</span>
                <div class="flex items-center space-x-1 text-sm text-muted-foreground">
                  <Calendar class="w-4 h-4" />
                  <span>{{ formatDate(activeProject.created_at) }}</span>
                </div>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">Owner</span>
                <div class="flex items-center space-x-1 text-sm text-muted-foreground">
                  <User class="w-4 h-4" />
                  <span>{{ activeProject.owner_id || 'Unknown' }}</span>
                </div>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">Status</span>
                <Badge variant="secondary">Active</Badge>
              </div>
            </CardContent>
          </Card>
        </div>
      </template>
    </LoadingWrapper>
  </div>
</template>
