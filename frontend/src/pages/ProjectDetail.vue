<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectsStore } from '@/stores/projects'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Database, 
  Table, 
  Workflow, 
  Activity, 
  Settings,
  Plus,
  ExternalLink,
  Calendar,
  User
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectsStore = useProjectsStore()

const projectId = computed(() => route.params.id as string)
const project = computed(() => projectsStore.currentProject)

const stats = ref({
  tables: 0,
  databases: 0,
  workflows: 0,
  lastActivity: '2 hours ago'
})

onMounted(async () => {
  try {
    await projectsStore.fetchProject(projectId.value)
  } catch (error) {
    console.error('Failed to fetch project:', error)
  }
})

const handleCreateTable = () => {
  router.push(`/projects/${projectId.value}/tables`)
}

const handleCreateDatabase = () => {
  router.push(`/projects/${projectId.value}/databases`)
}

const handleCreateWorkflow = () => {
  router.push(`/projects/${projectId.value}/workflows`)
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div v-if="project" class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">{{ project.name }}</h1>
        <p class="text-gray-600">{{ project.description || 'No description' }}</p>
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
          <p class="text-xs text-muted-foreground">
            Database tables
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Databases</CardTitle>
          <Database class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.databases }}</div>
          <p class="text-xs text-muted-foreground">
            Connected databases
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Workflows</CardTitle>
          <Workflow class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.workflows }}</div>
          <p class="text-xs text-muted-foreground">
            Active workflows
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Last Activity</CardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats.lastActivity }}</div>
          <p class="text-xs text-muted-foreground">
            Time since last update
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Main Content -->
    <Tabs default-value="overview" class="space-y-6">
      <TabsList>
        <TabsTrigger value="overview">Overview</TabsTrigger>
        <TabsTrigger value="tables">Tables</TabsTrigger>
        <TabsTrigger value="databases">Databases</TabsTrigger>
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
              <Button 
                variant="outline" 
                class="w-full justify-start"
                @click="handleCreateTable"
              >
                <Table class="w-4 h-4 mr-2" />
                Create New Table
              </Button>
              <Button 
                variant="outline" 
                class="w-full justify-start"
                @click="handleCreateDatabase"
              >
                <Database class="w-4 h-4 mr-2" />
                Connect Database
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
                <div class="flex items-center space-x-1 text-sm text-gray-500">
                  <Calendar class="w-4 h-4" />
                  <span>{{ project?.created_at ? new Date(project.created_at).toLocaleDateString() : 'Unknown' }}</span>
                </div>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">Owner</span>
                <div class="flex items-center space-x-1 text-sm text-gray-500">
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
            <div class="text-center py-8">
              <Table class="mx-auto h-12 w-12 text-gray-400" />
              <h3 class="mt-2 text-sm font-medium text-gray-900">No tables yet</h3>
              <p class="mt-1 text-sm text-gray-500">Get started by creating your first table.</p>
              <div class="mt-6">
                <Button @click="handleCreateTable">Create Table</Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="databases">
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <div>
                <CardTitle>Databases</CardTitle>
                <CardDescription>Manage your database connections</CardDescription>
              </div>
              <Button @click="handleCreateDatabase">
                <Plus class="w-4 h-4 mr-2" />
                Connect Database
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div class="text-center py-8">
              <Database class="mx-auto h-12 w-12 text-gray-400" />
              <h3 class="mt-2 text-sm font-medium text-gray-900">No databases connected</h3>
              <p class="mt-1 text-sm text-gray-500">Connect a database to start managing your data.</p>
              <div class="mt-6">
                <Button @click="handleCreateDatabase">Connect Database</Button>
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
            <div class="text-center py-8">
              <Workflow class="mx-auto h-12 w-12 text-gray-400" />
              <h3 class="mt-2 text-sm font-medium text-gray-900">No workflows yet</h3>
              <p class="mt-1 text-sm text-gray-500">Create workflows to automate your database operations.</p>
              <div class="mt-6">
                <Button @click="handleCreateWorkflow">Create Workflow</Button>
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
              <Activity class="mx-auto h-12 w-12 text-gray-400" />
              <h3 class="mt-2 text-sm font-medium text-gray-900">No activity yet</h3>
              <p class="mt-1 text-sm text-gray-500">Activity will appear here as you work on your project.</p>
            </div>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
