<script setup lang="ts">
import ProjectSelector from '@/components/ProjectSelector.vue'
import {
  Edit,
  Eye,
  FolderOpen,
  Home,
  LogOut,
  MoreVertical,
  Pause,
  Play,
  Plus,
  Settings as SettingsIcon,
  Table as TableIcon,
  Trash2,
  User,
  Workflow as WorkflowIcon,
} from 'lucide-vue-next'
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { tableService } from '@/api/services/table'
import { workflowService } from '@/api/services/workflow'
import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/projects'
import { useTableStore } from '@/stores/tables'
import { useSidebarItemsStore } from '@/stores/ui/sidebar-items'
import { useWorkflowStore } from '@/stores/workflows'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
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

import { useToast } from '@/composables/ui'

const router = useRouter()
const route = useRoute()
const sidebarStore = useSidebarItemsStore()
const authStore = useAuthStore()
const projectStore = useProjectStore()
const tableStore = useTableStore()
const workflowStore = useWorkflowStore()
const toast = useToast()

const currentProjectId = computed(() => {
  return (route.params.projectId as string) || projectStore.currentProject?.id
})

const tables = computed(() => tableStore.currentTables)
const workflows = computed(() => workflowStore.currentWorkflows)
const isLoadingTables = computed(() => tableStore.isLoading)
const isLoadingWorkflows = computed(() => workflowStore.isLoading)

const loadProjectData = async (projectId: string) => {
  await Promise.all([tableStore.loadTables(projectId), workflowStore.loadWorkflows(projectId)])

  sidebarStore.updateProjectCounts(tableStore.tableCount, workflowStore.workflowCount)
}

watch(
  currentProjectId,
  async (newProjectId) => {
    if (newProjectId) {
      tableStore.setCurrentProject(newProjectId)
      workflowStore.setCurrentProject(newProjectId)
      await loadProjectData(newProjectId)
    } else {
      sidebarStore.clearProjectItems()
      tableStore.setCurrentProject(null)
      workflowStore.setCurrentProject(null)
    }
  },
  { immediate: true }
)

onMounted(() => {
  if (currentProjectId.value) {
    loadProjectData(currentProjectId.value)
  }
})

const isActive = (itemUrl: string) => {
  return route.path === itemUrl || route.path.startsWith(itemUrl + '/')
}

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
    if (currentProjectId.value) {
      tableStore.removeTable(currentProjectId.value, tableId)
      sidebarStore.updateProjectCounts(tableStore.tableCount, workflowStore.workflowCount)
    }
  } catch (error: any) {
    toast.error('Error', `Failed to delete table: ${error.message}`)
  }
}

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

const handleToggleWorkflow = async (
  workflowId: string,
  isActive: boolean,
  workflowName: string
) => {
  try {
    await workflowService.toggle(currentProjectId.value!, workflowId)
    toast.success(
      isActive ? 'Deactivated' : 'Activated',
      `Workflow "${workflowName}" has been ${isActive ? 'deactivated' : 'activated'}`
    )
    if (currentProjectId.value) {
      await workflowStore.loadWorkflows(currentProjectId.value)
    }
  } catch (error: any) {
    toast.error('Error', `Failed to toggle workflow: ${error.message}`)
  }
}

const handleDeleteWorkflow = async (workflowId: string, workflowName: string) => {
  if (!confirm(`Are you sure you want to delete workflow "${workflowName}"?`)) {
    return
  }

  try {
    await workflowService.delete(currentProjectId.value!, workflowId)
    toast.success('Deleted', `Workflow "${workflowName}" has been deleted`)
    if (currentProjectId.value) {
      workflowStore.removeWorkflow(currentProjectId.value, workflowId)
      sidebarStore.updateProjectCounts(tableStore.tableCount, workflowStore.workflowCount)
    }
  } catch (error: any) {
    toast.error('Error', `Failed to delete workflow: ${error.message}`)
  }
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
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
              <SidebarMenuButton
                @click="router.push('/dashboard')"
                :isActive="isActive('/dashboard')"
              >
                <Home class="w-4 h-4" />
                <span>Dashboard</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton
                @click="router.push('/projects')"
                :isActive="route.path === '/projects'"
              >
                <FolderOpen class="w-4 h-4" />
                <span>Projects</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton
                @click="router.push('/settings')"
                :isActive="isActive('/settings')"
              >
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
          <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="handleCreateTable">
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
          <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="handleCreateWorkflow">
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
                    <DropdownMenuItem
                      @click="handleToggleWorkflow(workflow.id, workflow.is_active, workflow.name)"
                    >
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
