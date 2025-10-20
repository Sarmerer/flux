<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  Table as TableIcon,
  Workflow as WorkflowIcon,
  Plus,
  MoreVertical,
  Edit,
  Trash2,
  Settings as SettingsIcon,
  Play,
  Pause,
  User,
  Home,
  LogOut,
  Eye,
  FolderOpen
} from 'lucide-vue-next'

import { useAuthStore } from '@/stores/auth'
import { tableService } from '@/api/services/table'
import { workflowService } from '@/api/services/workflow'
import { useToast } from '@/composables/ui'
import ProjectSelector from '@/components/ProjectSelector.vue'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

import { Button } from '@/components/ui/button'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toast = useToast()

// State
const tables = ref<any[]>([])
const workflows = ref<any[]>([])
const isLoadingTables = ref(false)
const isLoadingWorkflows = ref(false)

// Computed
const currentProjectId = computed(() => route.params.projectId as string)

const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/')
}

// Methods
const loadTables = async () => {
  if (!currentProjectId.value) {
    tables.value = []
    return
  }

  isLoadingTables.value = true
  try {
    const response = await tableService.getAll(currentProjectId.value)
    tables.value = Array.isArray(response) ? response : []
  } catch (error: any) {
    console.error('Failed to load tables:', error)
    tables.value = []
  } finally {
    isLoadingTables.value = false
  }
}

const loadWorkflows = async () => {
  if (!currentProjectId.value) {
    workflows.value = []
    return
  }

  isLoadingWorkflows.value = true
  try {
    const response = await workflowService.getAll(currentProjectId.value)
    workflows.value = Array.isArray(response) ? response : []
  } catch (error: any) {
    console.error('Failed to load workflows:', error)
    workflows.value = []
  } finally {
    isLoadingWorkflows.value = false
  }
}

const loadProjectData = async () => {
  await Promise.all([loadTables(), loadWorkflows()])
}

// Table Actions
const handleCreateTable = () => {
  if (!currentProjectId.value) {
    toast.warning('No Project Selected', 'Please select a project first')
    return
  }
  router.push(`/projects/${currentProjectId.value}/tables?action=create`)
}

const handleViewTable = (tableId: string) => {
  router.push(`/projects/${currentProjectId.value}/tables/${tableId}/data`)
}

const handleEditTable = (tableId: string) => {
  router.push(`/projects/${currentProjectId.value}/tables/${tableId}/builder`)
}

const handleDeleteTable = async (tableId: string, tableName: string) => {
  if (!confirm(`Are you sure you want to delete table "${tableName}"?`)) {
    return
  }

  try {
    await tableService.delete(tableId)
    toast.success('Deleted', `Table "${tableName}" has been deleted`)
    await loadTables()
  } catch (error: any) {
    toast.error('Error', `Failed to delete table: ${error.message}`)
  }
}

// Workflow Actions
const handleCreateWorkflow = () => {
  if (!currentProjectId.value) {
    toast.warning('No Project Selected', 'Please select a project first')
    return
  }
  router.push(`/projects/${currentProjectId.value}/workflows?action=create`)
}

const handleEditWorkflow = (workflowId: string) => {
  router.push(`/projects/${currentProjectId.value}/workflows/${workflowId}/builder`)
}

const handleToggleWorkflow = async (workflowId: string, isActive: boolean, workflowName: string) => {
  try {
    await workflowService.toggle(currentProjectId.value, workflowId)
    toast.success(
      isActive ? 'Deactivated' : 'Activated',
      `Workflow "${workflowName}" has been ${isActive ? 'deactivated' : 'activated'}`
    )
    await loadWorkflows()
  } catch (error: any) {
    toast.error('Error', `Failed to toggle workflow: ${error.message}`)
  }
}

const handleDeleteWorkflow = async (workflowId: string, workflowName: string) => {
  if (!confirm(`Are you sure you want to delete workflow "${workflowName}"?`)) {
    return
  }

  try {
    await workflowService.delete(currentProjectId.value, workflowId)
    toast.success('Deleted', `Workflow "${workflowName}" has been deleted`)
    await loadWorkflows()
  } catch (error: any) {
    toast.error('Error', `Failed to delete workflow: ${error.message}`)
  }
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

// Lifecycle
onMounted(() => {
  if (currentProjectId.value) {
    loadProjectData()
  }
})

// Watch for route changes
watch(currentProjectId, (newProjectId) => {
  if (newProjectId) {
    loadProjectData()
  } else {
    tables.value = []
    workflows.value = []
  }
})
</script>

<template>
  <Sidebar collapsible="icon">
    <SidebarHeader class="p-4 space-y-3">
      <div class="flex items-center space-x-2">
        <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
          <span class="text-primary-foreground font-bold text-sm">F</span>
        </div>
        <span class="font-semibold text-foreground">Flow</span>
      </div>

      <!-- Project Selector -->
      <ProjectSelector v-if="currentProjectId" />
    </SidebarHeader>

    <SidebarContent>
      <!-- Main Navigation -->
      <SidebarGroup>
        <SidebarGroupLabel>Navigation</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton @click="router.push('/dashboard')" :isActive="isActive('/dashboard')">
                <Home class="w-4 h-4" />
                <span>Dashboard</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton @click="router.push('/projects')" :isActive="route.path === '/projects'">
                <FolderOpen class="w-4 h-4" />
                <span>Projects</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton @click="router.push('/settings')" :isActive="isActive('/settings')">
                <SettingsIcon class="w-4 h-4" />
                <span>Settings</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- Tables Section -->
      <SidebarGroup v-if="currentProjectId">
        <SidebarGroupLabel class="flex items-center justify-between">
          <span>Tables</span>
          <Button
            variant="ghost"
            size="sm"
            class="h-6 w-6 p-0"
            @click="handleCreateTable"
          >
            <Plus class="w-4 h-4" />
          </Button>
        </SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <div v-if="isLoadingTables" class="px-2 py-2 text-sm text-muted-foreground">
              Loading...
            </div>
            <SidebarMenuItem v-for="table in tables" :key="table.id">
              <div class="group relative w-full">
                <SidebarMenuButton
                  @click="handleViewTable(table.id)"
                  :isActive="isActive(`/projects/${currentProjectId}/tables/${table.id}`)"
                  class="pr-8"
                >
                  <TableIcon class="w-4 h-4" />
                  <span class="flex-1 truncate">{{ table.name }}</span>
                </SidebarMenuButton>

                <!-- Table Actions -->
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      class="absolute right-1 top-1 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <MoreVertical class="w-4 h-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem @click="handleViewTable(table.id)">
                      <Eye class="w-4 h-4 mr-2" />
                      View Data
                    </DropdownMenuItem>
                    <DropdownMenuItem @click="handleEditTable(table.id)">
                      <Edit class="w-4 h-4 mr-2" />
                      Edit Schema
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      @click="handleDeleteTable(table.id, table.name)"
                      class="text-destructive focus:text-destructive"
                    >
                      <Trash2 class="w-4 h-4 mr-2" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </SidebarMenuItem>
            <div v-if="!isLoadingTables && !tables.length" class="px-2 py-2">
              <p class="text-xs text-muted-foreground italic">No tables yet</p>
            </div>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- Workflows Section -->
      <SidebarGroup v-if="currentProjectId">
        <SidebarGroupLabel class="flex items-center justify-between">
          <span>Workflows</span>
          <Button
            variant="ghost"
            size="sm"
            class="h-6 w-6 p-0"
            @click="handleCreateWorkflow"
          >
            <Plus class="w-4 h-4" />
          </Button>
        </SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <div v-if="isLoadingWorkflows" class="px-2 py-2 text-sm text-muted-foreground">
              Loading...
            </div>
            <SidebarMenuItem v-for="workflow in workflows" :key="workflow.id">
              <div class="group relative w-full">
                <SidebarMenuButton
                  @click="handleEditWorkflow(workflow.id)"
                  :isActive="isActive(`/projects/${currentProjectId}/workflows/${workflow.id}`)"
                  class="pr-8"
                >
                  <WorkflowIcon
                    class="w-4 h-4"
                    :class="workflow.is_active ? 'text-success' : 'text-muted-foreground'"
                  />
                  <span class="flex-1 truncate">{{ workflow.name }}</span>
                </SidebarMenuButton>

                <!-- Workflow Actions -->
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      class="absolute right-1 top-1 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <MoreVertical class="w-4 h-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem @click="handleEditWorkflow(workflow.id)">
                      <Edit class="w-4 h-4 mr-2" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem @click="handleToggleWorkflow(workflow.id, workflow.is_active, workflow.name)">
                      <component :is="workflow.is_active ? Pause : Play" class="w-4 h-4 mr-2" />
                      {{ workflow.is_active ? 'Deactivate' : 'Activate' }}
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      @click="handleDeleteWorkflow(workflow.id, workflow.name)"
                      class="text-destructive focus:text-destructive"
                    >
                      <Trash2 class="w-4 h-4 mr-2" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </SidebarMenuItem>
            <div v-if="!isLoadingWorkflows && !workflows.length" class="px-2 py-2">
              <p class="text-xs text-muted-foreground italic">No workflows yet</p>
            </div>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- No Project Selected State -->
      <SidebarGroup v-if="!currentProjectId">
        <SidebarGroupContent>
          <div class="px-4 py-6 text-center">
            <FolderOpen class="w-8 h-8 mx-auto text-muted-foreground mb-2" />
            <p class="text-sm text-muted-foreground mb-3">No project selected</p>
            <Button variant="outline" size="sm" @click="router.push('/projects')">
              View Projects
            </Button>
          </div>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <!-- User Footer -->
    <SidebarFooter class="p-4 border-t">
      <div class="space-y-2">
        <div class="flex items-center space-x-2 text-sm text-muted-foreground px-2">
          <User class="w-4 h-4" />
          <span class="truncate">{{ authStore.user?.name || 'User' }}</span>
        </div>
        <Button variant="ghost" size="sm" class="w-full justify-start" @click="handleLogout">
          <LogOut class="w-4 h-4 mr-2" />
          Logout
        </Button>
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
