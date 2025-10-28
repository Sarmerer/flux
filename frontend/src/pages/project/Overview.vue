<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import {
  Activity,
  ArrowRight,
  Calendar,
  CheckCircle,
  Clock,
  Database,
  Settings,
  Table,
  Users,
  Workflow,
  Zap,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useTableStore } from '@/stores/tables'
import { useWorkflowStore } from '@/stores/workflows'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { activityService } from '@/api/services/activity'
import { useFormatting } from '@/composables/formatting'
import { useActivityFormatting } from '@/composables/formatting/useActivityFormatting'
import type { ActivityLog } from '@/types/api'

const router = useRouter()
const { formatDate } = useFormatting()
const { getActivityIcon: getActivityIconFromType } = useActivityFormatting()

const activeProjectStore = useActiveProjectStore()
const activeProject = computed(() => activeProjectStore.activeProject)
const activeProjectId = computed(() => activeProject.value?.id)

const tablesStore = useTableStore()
const workflowsStore = useWorkflowStore()

const recentActivity = ref<ActivityLog[]>([])
const isLoadingActivity = ref(false)

const statsCards = computed(() => [
  {
    title: 'Tables',
    value: tablesStore.length,
    description: 'Database tables',
    icon: Table,
    color: 'text-blue-600',
  },
  {
    title: 'Workflows',
    value: workflowsStore.length,
    description: 'Active workflows',
    icon: Workflow,
    color: 'text-purple-600',
  },
  {
    title: 'Members',
    value: '3',
    description: 'Team members',
    icon: Users,
    color: 'text-green-600',
  },
])

const quickActions = [
  {
    icon: Table,
    title: 'Create Table',
    description: 'Add a new database table',
    action: () => router.push(`/projects/${activeProjectId.value}/tables`),
  },
  {
    icon: Workflow,
    title: 'Create Workflow',
    description: 'Automate your database',
    action: () => router.push(`/projects/${activeProjectId.value}/workflows`),
  },
  {
    icon: Users,
    title: 'Invite Members',
    description: 'Collaborate with your team',
    action: () => router.push(`/projects/${activeProjectId.value}/members`),
  },
]

const loadRecentActivity = async () => {
  if (!activeProjectId.value) return

  isLoadingActivity.value = true
  try {
    const response = await activityService.get(activeProjectId.value, 1, 5)
    recentActivity.value = response.logs || []
  } catch (error: any) {
    console.error('Failed to load activity:', error)
  } finally {
    isLoadingActivity.value = false
  }
}

onMounted(() => {
  loadRecentActivity()
})

const formatActivityTime = (timestamp: string) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return 'Just now'
  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  return `${days}d ago`
}

const getActivityDescription = (activity: ActivityLog): string => {
  const entityName = activity.entity_name
  switch (activity.type) {
    case 'table_created':
      return `Created table "${entityName}"`
    case 'table_updated':
      return `Updated table "${entityName}"`
    case 'table_deleted':
      return `Deleted table "${entityName}"`
    case 'row_created':
      return `Added row to "${entityName}"`
    case 'row_updated':
      return `Updated row in "${entityName}"`
    case 'row_deleted':
      return `Deleted row from "${entityName}"`
    case 'workflow_triggered':
      return `Triggered workflow "${entityName}"`
    case 'workflow_completed':
      return `Completed workflow "${entityName}"`
    case 'workflow_failed':
      return `Failed workflow "${entityName}"`
    default:
      return `Activity on "${entityName}"`
  }
}
</script>

<template>
  <div class="h-full overflow-auto">
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
        <div class="border-b bg-gradient-to-b from-muted/30 to-background">
          <div class="px-6 py-4">
            <div class="flex items-start justify-between mb-4">
              <div class="flex-1">
                <h1 class="text-2xl font-bold text-foreground mb-1">{{ activeProject.name }}</h1>
                <p class="text-muted-foreground text-xs max-w-2xl">
                  {{ activeProject.description || 'No description provided' }}
                </p>
                <div class="flex items-center gap-3 mt-2">
                  <div class="flex items-center text-xs text-muted-foreground">
                    <Calendar class="w-3 h-3 mr-1" />
                    {{ formatDate(activeProject.created_at) }}
                  </div>
                  <Badge variant="secondary" class="flex items-center gap-1 h-5 text-xs">
                    <CheckCircle class="w-2.5 h-2.5" />
                    Active
                  </Badge>
                </div>
              </div>
              <Button variant="outline" size="sm">
                <Settings class="w-4 h-4 mr-2" />
                Settings
              </Button>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
              <Card
                v-for="stat in statsCards"
                :key="stat.title"
                class="border-border/50 hover:border-border/80 transition-colors"
              >
                <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-1.5 px-4 pt-3">
                  <CardTitle class="text-xs font-medium">{{ stat.title }}</CardTitle>
                  <component :is="stat.icon" :class="['h-4 w-4', stat.color]" />
                </CardHeader>
                <CardContent class="px-4 pb-3">
                  <div class="text-2xl font-bold">{{ stat.value }}</div>
                  <p class="text-xs text-muted-foreground">{{ stat.description }}</p>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>

        <div class="p-4">
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <div class="lg:col-span-2 space-y-4">
              <Card>
                <CardHeader class="pb-3">
                  <CardTitle class="text-base">Quick Actions</CardTitle>
                  <CardDescription class="text-xs">Get started with your project</CardDescription>
                </CardHeader>
                <CardContent class="pb-4">
                  <div class="grid grid-cols-3 gap-2">
                    <button
                      v-for="action in quickActions"
                      :key="action.title"
                      @click="action.action"
                      class="flex flex-col items-center gap-2 p-3 rounded-lg border border-border hover:border-primary/50 hover:bg-muted/50 transition-all group"
                    >
                      <div class="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                        <component :is="action.icon" class="w-4 h-4 text-primary" />
                      </div>
                      <div class="text-center">
                        <div class="text-xs font-semibold">{{ action.title }}</div>
                      </div>
                    </button>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader class="pb-3">
                  <CardTitle class="text-base">Resources</CardTitle>
                  <CardDescription class="text-xs">Your project's tables and workflows</CardDescription>
                </CardHeader>
                <CardContent class="space-y-2 pb-4">
                  <button
                    @click="router.push(`/projects/${activeProjectId}/tables`)"
                    class="w-full flex items-center justify-between p-3 rounded-lg border border-border hover:border-primary/50 hover:bg-muted/50 transition-all group"
                  >
                    <div class="flex items-center gap-2.5">
                      <div class="w-8 h-8 rounded-lg bg-blue-600/10 flex items-center justify-center">
                        <Database class="w-4 h-4 text-blue-600" />
                      </div>
                      <div class="text-left">
                        <div class="text-sm font-semibold">Tables</div>
                        <div class="text-xs text-muted-foreground">
                          {{ tablesStore.length }} {{ tablesStore.length === 1 ? 'table' : 'tables' }}
                        </div>
                      </div>
                    </div>
                    <ArrowRight class="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
                  </button>

                  <button
                    @click="router.push(`/projects/${activeProjectId}/workflows`)"
                    class="w-full flex items-center justify-between p-3 rounded-lg border border-border hover:border-primary/50 hover:bg-muted/50 transition-all group"
                  >
                    <div class="flex items-center gap-2.5">
                      <div class="w-8 h-8 rounded-lg bg-purple-600/10 flex items-center justify-center">
                        <Zap class="w-4 h-4 text-purple-600" />
                      </div>
                      <div class="text-left">
                        <div class="text-sm font-semibold">Workflows</div>
                        <div class="text-xs text-muted-foreground">
                          {{ workflowsStore.length }} {{ workflowsStore.length === 1 ? 'workflow' : 'workflows' }}
                        </div>
                      </div>
                    </div>
                    <ArrowRight class="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
                  </button>
                </CardContent>
              </Card>
            </div>

            <div class="space-y-4">
              <Card>
                <CardHeader class="pb-3">
                  <div class="flex items-center justify-between">
                    <div>
                      <CardTitle class="text-base">Recent Activity</CardTitle>
                      <CardDescription class="text-xs">Latest updates</CardDescription>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      class="h-7 text-xs"
                      @click="router.push(`/projects/${activeProjectId}/activity`)"
                    >
                      View All
                    </Button>
                  </div>
                </CardHeader>
                <CardContent class="pb-4">
                  <LoadingWrapper :is-loading="isLoadingActivity" loading-text="Loading...">
                    <div v-if="recentActivity.length > 0" class="space-y-3">
                      <div
                        v-for="activity in recentActivity"
                        :key="activity.id"
                        class="flex items-start gap-2.5"
                      >
                        <div class="flex-shrink-0 w-7 h-7 rounded-lg bg-muted flex items-center justify-center">
                          <component :is="getActivityIconFromType(activity)" class="w-3.5 h-3.5 text-muted-foreground" />
                        </div>
                        <div class="flex-1 min-w-0">
                          <p class="text-xs text-foreground leading-snug">{{ getActivityDescription(activity) }}</p>
                          <p class="text-xs text-muted-foreground flex items-center gap-1 mt-0.5">
                            <Clock class="h-2.5 w-2.5" />
                            <span>{{ formatActivityTime(activity.created_at) }}</span>
                          </p>
                        </div>
                      </div>
                    </div>
                    <div v-else class="text-center py-6">
                      <Activity class="w-7 h-7 text-muted-foreground mx-auto mb-1.5" />
                      <p class="text-xs text-muted-foreground">No recent activity</p>
                    </div>
                  </LoadingWrapper>
                </CardContent>
              </Card>

              <Card>
                <CardHeader class="pb-3">
                  <CardTitle class="text-base">Project Info</CardTitle>
                  <CardDescription class="text-xs">Project details</CardDescription>
                </CardHeader>
                <CardContent class="space-y-2 pb-4">
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-muted-foreground">Owner</span>
                    <span class="font-medium">{{ activeProject.owner_id || 'Unknown' }}</span>
                  </div>
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-muted-foreground">Created</span>
                    <span class="font-medium">{{ formatDate(activeProject.created_at) }}</span>
                  </div>
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-muted-foreground">Last Updated</span>
                    <span class="font-medium">{{ formatDate(activeProject.updated_at) }}</span>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </template>
    </LoadingWrapper>
  </div>
</template>
