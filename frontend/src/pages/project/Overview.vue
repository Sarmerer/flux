<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import {
  Activity,
  Calendar,
  ExternalLink,
  Plus,
  Settings,
  Table,
  User,
  Workflow,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useTableStore } from '@/stores/tables'
import { useWorkflowStore } from '@/stores/workflows'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { useFormatting } from '@/composables/formatting'

const route = useRoute()
const router = useRouter()

const { formatDate } = useFormatting()

const activeProjectStore = useActiveProjectStore()
const activeProject = computed(() => activeProjectStore.activeProject)
const activeProjectId = computed(() => activeProject.value?.id)

const tablesStore = useTableStore()
const workflowsStore = useWorkflowStore()

const currentTab = computed(() => (route.meta.tab as string) || 'overview')

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

const onTabChange = (tab: string | number) => {
  const tabString = String(tab)
  if (tabString === 'overview') {
    router.push(`/projects/${activeProjectId.value}/overview`)
  } else {
    router.push(`/projects/${activeProjectId.value}/overview/${tabString}`)
  }
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

        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
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

        <Tabs :model-value="currentTab" @update:model-value="onTabChange" class="space-y-6">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="tables">Tables</TabsTrigger>
            <TabsTrigger value="workflows">Workflows</TabsTrigger>
            <TabsTrigger value="activity">Activity</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" class="space-y-6">
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
          </TabsContent>

          <TabsContent value="tables">
            <Card>
              <CardHeader>
                <div class="flex items-center justify-between">
                  <div>
                    <CardTitle>Tables</CardTitle>
                    <CardDescription>Manage your database tables</CardDescription>
                  </div>
                  <Button @click="onTableCreateClick">
                    <Plus class="w-4 h-4 mr-2" />
                    New Table
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <LoadingWrapper
                  :is-loading="tablesStore.isLoading"
                  loading-text="Loading tables..."
                >
                  <ErrorState
                    v-if="tablesStore.error"
                    :error="tablesStore.error"
                    title="Failed to load tables"
                    :on-retry="
                      activeProjectId
                        ? () => tablesStore.loadByProjectId(activeProjectId!)
                        : undefined
                    "
                  />
                  <div v-else class="space-y-4">
                    <div
                      v-for="table in tablesStore.tables"
                      :key="table.id"
                      class="flex items-center justify-between p-4 border rounded-lg hover:bg-accent cursor-pointer"
                      @click="router.push(`/projects/${activeProjectId}/tables/${table.id}/data`)"
                    >
                      <div class="flex items-center space-x-3">
                        <div class="p-2 bg-blue-100 rounded-lg">
                          <Table class="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <h4 class="font-medium">{{ table.name }}</h4>
                          <p class="text-sm text-muted-foreground">
                            {{ table.description || 'No description' }}
                          </p>
                        </div>
                      </div>
                      <ExternalLink class="h-4 w-4 text-muted-foreground" />
                    </div>
                  </div>
                </LoadingWrapper>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="workflows">
            <Card>
              <CardHeader>
                <div class="flex items-center justify-between">
                  <div>
                    <CardTitle>Workflows</CardTitle>
                    <CardDescription>Automate your database operations</CardDescription>
                  </div>
                  <Button @click="onWorkflowCreateClick">
                    <Plus class="w-4 h-4 mr-2" />
                    New Workflow
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <LoadingWrapper
                  :is-loading="workflowsStore.isLoading"
                  loading-text="Loading workflows..."
                >
                  <ErrorState
                    v-if="workflowsStore.error"
                    :error="workflowsStore.error"
                    title="Failed to load workflows"
                    :on-retry="
                      activeProjectId
                        ? () => workflowsStore.loadByProjectId(activeProjectId!)
                        : undefined
                    "
                  />
                  <div v-else class="space-y-4">
                    <div
                      v-for="workflow in workflowsStore.workflows"
                      :key="workflow.id"
                      class="flex items-center justify-between p-4 border rounded-lg hover:bg-accent cursor-pointer"
                      @click="
                        router.push(`/projects/${activeProjectId}/workflows/${workflow.id}/builder`)
                      "
                    >
                      <div class="flex items-center space-x-3">
                        <div class="p-2 bg-purple-100 rounded-lg">
                          <Workflow class="h-5 w-5 text-purple-600" />
                        </div>
                        <div>
                          <h4 class="font-medium">{{ workflow.name }}</h4>
                          <p class="text-sm text-muted-foreground">
                            {{ workflow.description || 'No description' }}
                          </p>
                        </div>
                      </div>
                      <ExternalLink class="h-4 w-4 text-muted-foreground" />
                    </div>
                  </div>
                </LoadingWrapper>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="activity">
            <Card>
              <CardHeader>
                <CardTitle>Activity Log</CardTitle>
                <CardDescription>Recent activity in this project</CardDescription>
              </CardHeader>
              <CardContent>
                <div class="text-center py-8">
                  <Activity class="mx-auto h-12 w-12 text-muted-foreground" />
                  <h3 class="mt-2 text-sm font-medium text-foreground">No activity yet</h3>
                  <p class="mt-1 text-sm text-muted-foreground">
                    Activity will appear here as you work on your project.
                  </p>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </template>
    </LoadingWrapper>
  </div>
</template>
