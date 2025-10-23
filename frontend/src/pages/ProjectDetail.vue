<script setup lang="ts">
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'
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
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useProjectStore } from '@/stores/projects'
import { useSidebarItemsStore } from '@/stores/ui/sidebar-items'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { useProjects } from '@/composables/api'
import { useTables } from '@/composables/api'
import { useWorkflows } from '@/composables/api'
import { useFormatting } from '@/composables/formatting'
import { useToast } from '@/composables/ui'

const route = useRoute()
const router = useRouter()
const projectsStore = useProjectStore()
const sidebarStore = useSidebarItemsStore()
const toast = useToast()
const { formatDate } = useFormatting()

const projectId = computed(() => route.params.projectId as string)
const project = computed(() => projectsStore.currentProject)

const currentTab = computed(() => (route.meta.tab as string) || 'overview')

const { getProjectById } = useProjects()
const { tables, loading: tablesLoading } = useTables(projectId.value)
const { workflows, loading: workflowsLoading } = useWorkflows(projectId.value)

watch(
  [tables, workflows],
  ([newTables, newWorkflows]) => {
    if (projectId.value) {
      sidebarStore.setProjectItems(
        projectId.value,
        newTables?.length || 0,
        newWorkflows?.length || 0
      )
    }
  },
  { immediate: true }
)

const stats = computed(() => ({
  tables: tables.value?.length || 0,
  workflows: workflows.value?.length || 0,
  lastActivity: '2 hours ago',
}))

const isLoading = computed(() => !project.value && projectsStore.isLoading)

onMounted(async () => {
  if (!project.value || project.value.id !== projectId.value) {
    try {
      const fetchedProject = await getProjectById(projectId.value)
      projectsStore.setCurrentProject(fetchedProject)
    } catch (error: any) {
      console.error('Failed to fetch project:', error)
      toast.error('Error', error.message || 'Failed to load project')
    }
  }

  if (projectId.value) {
    sidebarStore.setProjectItems(
      projectId.value,
      tables.value?.length || 0,
      workflows.value?.length || 0
    )
  }
})

const handleCreateTable = () => {
  router.push(`/projects/${projectId.value}/tables`)
}

const handleCreateWorkflow = () => {
  router.push(`/projects/${projectId.value}/workflows`)
}

const handleTabChange = (tab: string | number) => {
  const tabString = String(tab)
  if (tabString === 'overview') {
    router.push(`/projects/${projectId.value}/overview`)
  } else {
    router.push(`/projects/${projectId.value}/overview/${tabString}`)
  }
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Loading State -->
    <LoadingSpinner v-if="isLoading" text="Loading project..." />

    <template v-else-if="project">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-foreground">{{ project.name }}</h1>
          <p class="text-muted-foreground">{{ project.description || 'No description' }}</p>
        </div>
        <div class="flex items-center space-x-2">
          <Button variant="outline">
            <Settings class="w-4 h-4 mr-2" />
            Settings
          </Button>
          <Button>
            <Plus class="w-4 h-4 mr-2" />
            Add Resource
          </Button>
        </div>
      </div>

      <!-- Stats -->
      <div v-if="project" class="grid grid-cols-1 md:grid-cols-4 gap-6">
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

      <!-- Main Content -->
      <Tabs :model-value="currentTab" @update:model-value="handleTabChange" class="space-y-6">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="tables">Tables</TabsTrigger>
          <TabsTrigger value="workflows">Workflows</TabsTrigger>
          <TabsTrigger value="activity">Activity</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" class="space-y-6">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Quick Actions -->
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
                <CardDescription>Get started with your project</CardDescription>
              </CardHeader>
              <CardContent class="space-y-3">
                <Button variant="outline" class="w-full justify-start" @click="handleCreateTable">
                  <Table class="w-4 h-4 mr-2" />
                  Create New Table
                </Button>
                <Button
                  variant="outline"
                  class="w-full justify-start"
                  @click="handleCreateWorkflow"
                >
                  <Workflow class="w-4 h-4 mr-2" />
                  Create Workflow
                </Button>
              </CardContent>
            </Card>

            <!-- Project Info -->
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
                    <span>{{ formatDate(project.created_at) }}</span>
                  </div>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm font-medium">Owner</span>
                  <div class="flex items-center space-x-1 text-sm text-muted-foreground">
                    <User class="w-4 h-4" />
                    <span>{{ project?.owner_id || 'Unknown' }}</span>
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
                <Button @click="handleCreateTable">
                  <Plus class="w-4 h-4 mr-2" />
                  New Table
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div v-if="tablesLoading" class="text-center py-8">
                <div
                  class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"
                ></div>
                <p class="mt-2 text-sm text-muted-foreground">Loading tables...</p>
              </div>
              <div v-else-if="tables.length === 0" class="text-center py-8">
                <Table class="mx-auto h-12 w-12 text-muted-foreground" />
                <h3 class="mt-2 text-sm font-medium text-foreground">No tables yet</h3>
                <p class="mt-1 text-sm text-muted-foreground">
                  Get started by creating your first table.
                </p>
                <div class="mt-6">
                  <Button @click="handleCreateTable">Create Table</Button>
                </div>
              </div>
              <div v-else class="space-y-4">
                <div
                  v-for="table in tables"
                  :key="table.id"
                  class="flex items-center justify-between p-4 border rounded-lg hover:bg-accent cursor-pointer"
                  @click="router.push(`/projects/${projectId}/tables/${table.id}/data`)"
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
                <Button @click="handleCreateWorkflow">
                  <Plus class="w-4 h-4 mr-2" />
                  New Workflow
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div v-if="workflowsLoading" class="text-center py-8">
                <div
                  class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"
                ></div>
                <p class="mt-2 text-sm text-muted-foreground">Loading workflows...</p>
              </div>
              <div v-else-if="workflows.length === 0" class="text-center py-8">
                <Workflow class="mx-auto h-12 w-12 text-muted-foreground" />
                <h3 class="mt-2 text-sm font-medium text-foreground">No workflows yet</h3>
                <p class="mt-1 text-sm text-muted-foreground">
                  Create workflows to automate your database operations.
                </p>
                <div class="mt-6">
                  <Button @click="handleCreateWorkflow">Create Workflow</Button>
                </div>
              </div>
              <div v-else class="space-y-4">
                <div
                  v-for="workflow in workflows"
                  :key="workflow.id"
                  class="flex items-center justify-between p-4 border rounded-lg hover:bg-accent cursor-pointer"
                  @click="router.push(`/projects/${projectId}/workflows/${workflow.id}/builder`)"
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

    <!-- Error State -->
    <div v-else class="text-center py-12">
      <h3 class="text-lg font-medium text-foreground">Project not found</h3>
      <p class="text-sm text-muted-foreground mt-2">
        The project you're looking for doesn't exist.
      </p>
      <Button variant="outline" class="mt-4" @click="router.push('/projects')">
        Back to Projects
      </Button>
    </div>
  </div>
</template>
