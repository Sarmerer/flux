<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  ChevronDown,
  ChevronRight,
  FolderOpen,
  Table as TableIcon,
  Workflow as WorkflowIcon,
  Plus,
  MoreVertical,
  Edit,
  Trash2,
  Copy,
  Settings as SettingsIcon,
  Play,
  Pause,
  User,
  Home,
  LogOut
} from 'lucide-vue-next'

import { useAuthStore } from '@/stores/auth'
import { projectService } from '@/api/services/project'
import { tableService } from '@/api/services/table'
import { workflowService } from '@/api/services/workflow'
import { useToast } from '@/composables/ui'

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
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
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
const projects = ref<any[]>([])
const expandedProjects = ref<Set<string>>(new Set())
const projectTables = ref<Map<string, any[]>>(new Map())
const projectWorkflows = ref<Map<string, any[]>>(new Map())
const isLoading = ref(false)

// Computed
const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/')
}

const isProjectExpanded = (projectId: string) => {
  return expandedProjects.value.has(projectId)
}

// Methods
const toggleProject = async (projectId: string) => {
  if (expandedProjects.value.has(projectId)) {
    expandedProjects.value.delete(projectId)
  } else {
    expandedProjects.value.add(projectId)
    // Load tables and workflows if not already loaded
    if (!projectTables.value.has(projectId)) {
      await loadProjectData(projectId)
    }
  }
}

const loadProjects = async () => {
  isLoading.value = true
  try {
    const response = await projectService.getAll()
    projects.value = response
  } catch (error: any) {
    console.error('Failed to load projects:', error)
    toast.error('Error', 'Failed to load projects')
  } finally {
    isLoading.value = false
  }
}

const loadProjectData = async (projectId: string) => {
  try {
    // Load tables
    const tablesResponse = await tableService.getAll(projectId)
    projectTables.value.set(projectId, tablesResponse)

    // Load workflows
    const workflowsResponse = await workflowService.getAll(projectId)
    projectWorkflows.value.set(projectId, workflowsResponse)
  } catch (error: any) {
    console.error('Failed to load project data:', error)
  }
}

// Actions
const handleCreateProject = () => {
  router.push('/projects?action=create')
}

const handleEditProject = (projectId: string) => {
  router.push(`/projects/${projectId}?action=edit`)
}

const handleDeleteProject = async (projectId: string, projectName: string) => {
  if (!confirm(`Are you sure you want to delete "${projectName}"? This action cannot be undone.`)) {
    return
  }

  try {
    await projectService.delete(projectId)
    toast.success('Deleted', `Project "${projectName}" has been deleted`)
    await loadProjects()
    if (route.params.projectId === projectId) {
      router.push('/projects')
    }
  } catch (error: any) {
    toast.error('Error', `Failed to delete project: ${error.message}`)
  }
}

const handleDuplicateProject = async (projectId: string, projectName: string) => {
  try {
    const project = projects.value.find(p => p.id === projectId)
    if (!project) return

    await projectService.create({
      name: `${projectName} (Copy)`,
      description: project.description
    })
    toast.success('Duplicated', `Project "${projectName}" has been duplicated`)
    await loadProjects()
  } catch (error: any) {
    toast.error('Error', `Failed to duplicate project: ${error.message}`)
  }
}

const handleCreateTable = (projectId: string) => {
  router.push(`/projects/${projectId}/tables?action=create`)
}

const handleEditTable = (projectId: string, tableId: string) => {
  router.push(`/projects/${projectId}/tables/${tableId}/builder`)
}

const handleDeleteTable = async (projectId: string, tableId: string, tableName: string) => {
  if (!confirm(`Are you sure you want to delete table "${tableName}"?`)) {
    return
  }

  try {
    await tableService.delete(tableId)
    toast.success('Deleted', `Table "${tableName}" has been deleted`)
    await loadProjectData(projectId)
  } catch (error: any) {
    toast.error('Error', `Failed to delete table: ${error.message}`)
  }
}

const handleCreateWorkflow = (projectId: string) => {
  router.push(`/projects/${projectId}/workflows?action=create`)
}

const handleEditWorkflow = (projectId: string, workflowId: string) => {
  router.push(`/projects/${projectId}/workflows/${workflowId}/builder`)
}

const handleToggleWorkflow = async (projectId: string, workflowId: string, isActive: boolean) => {
  try {
    await workflowService.toggle(projectId, workflowId)
    toast.success(
      isActive ? 'Deactivated' : 'Activated',
      `Workflow has been ${isActive ? 'deactivated' : 'activated'}`
    )
    await loadProjectData(projectId)
  } catch (error: any) {
    toast.error('Error', `Failed to toggle workflow: ${error.message}`)
  }
}

const handleDeleteWorkflow = async (projectId: string, workflowId: string, workflowName: string) => {
  if (!confirm(`Are you sure you want to delete workflow "${workflowName}"?`)) {
    return
  }

  try {
    await workflowService.delete(projectId, workflowId)
    toast.success('Deleted', `Workflow "${workflowName}" has been deleted`)
    await loadProjectData(projectId)
  } catch (error: any) {
    toast.error('Error', `Failed to delete workflow: ${error.message}`)
  }
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

// Lifecycle
onMounted(async () => {
  await loadProjects()

  // Auto-expand current project if in a project route
  const currentProjectId = route.params.projectId as string
  if (currentProjectId) {
    expandedProjects.value.add(currentProjectId)
    await loadProjectData(currentProjectId)
  }
})

// Watch for route changes to expand relevant project
watch(() => route.params.projectId, async (newProjectId) => {
  if (newProjectId && typeof newProjectId === 'string') {
    expandedProjects.value.add(newProjectId)
    if (!projectTables.value.has(newProjectId)) {
      await loadProjectData(newProjectId)
    }
  }
})
</script>

<template>
  <Sidebar collapsible="icon">
    <SidebarHeader class="p-4">
      <div class="flex items-center space-x-2">
        <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
          <span class="text-primary-foreground font-bold text-sm">F</span>
        </div>
        <span class="font-semibold text-foreground">Flow</span>
      </div>
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
              <SidebarMenuButton @click="router.push('/settings')" :isActive="isActive('/settings')">
                <SettingsIcon class="w-4 h-4" />
                <span>Settings</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- Projects Section -->
      <SidebarGroup>
        <SidebarGroupLabel class="flex items-center justify-between">
          <span>Projects</span>
          <Button
            variant="ghost"
            size="sm"
            class="h-6 w-6 p-0"
            @click="handleCreateProject"
          >
            <Plus class="w-4 h-4" />
          </Button>
        </SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="project in projects" :key="project.id">
              <!-- Project Item -->
              <div class="group relative">
                <SidebarMenuButton
                  @click="toggleProject(project.id)"
                  :isActive="isActive(`/projects/${project.id}`)"
                  class="pr-8"
                >
                  <component
                    :is="isProjectExpanded(project.id) ? ChevronDown : ChevronRight"
                    class="w-4 h-4"
                  />
                  <FolderOpen class="w-4 h-4" />
                  <span class="flex-1 truncate">{{ project.name }}</span>
                </SidebarMenuButton>

                <!-- Project Actions Menu -->
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
                    <DropdownMenuItem @click="router.push(`/projects/${project.id}`)">
                      <FolderOpen class="w-4 h-4 mr-2" />
                      Open
                    </DropdownMenuItem>
                    <DropdownMenuItem @click="handleEditProject(project.id)">
                      <Edit class="w-4 h-4 mr-2" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem @click="handleDuplicateProject(project.id, project.name)">
                      <Copy class="w-4 h-4 mr-2" />
                      Duplicate
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      @click="handleDeleteProject(project.id, project.name)"
                      class="text-destructive focus:text-destructive"
                    >
                      <Trash2 class="w-4 h-4 mr-2" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>

              <!-- Project Sub-items (Tables & Workflows) -->
              <SidebarMenuSub v-if="isProjectExpanded(project.id)">
                <!-- Tables Section -->
                <div class="space-y-1">
                  <div class="flex items-center justify-between px-2 py-1">
                    <span class="text-xs text-muted-foreground">Tables</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      class="h-4 w-4 p-0"
                      @click="handleCreateTable(project.id)"
                    >
                      <Plus class="w-3 h-3" />
                    </Button>
                  </div>
                  <SidebarMenuSubItem
                    v-for="table in projectTables.get(project.id) || []"
                    :key="table.id"
                  >
                    <div class="group relative w-full">
                      <SidebarMenuSubButton
                        @click="router.push(`/projects/${project.id}/tables/${table.id}/data`)"
                        :isActive="isActive(`/projects/${project.id}/tables/${table.id}`)"
                        class="pr-6"
                      >
                        <TableIcon class="w-4 h-4" />
                        <span class="flex-1 truncate">{{ table.name }}</span>
                      </SidebarMenuSubButton>

                      <!-- Table Actions -->
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            variant="ghost"
                            size="sm"
                            class="absolute right-0 top-0 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                          >
                            <MoreVertical class="w-3 h-3" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem @click="router.push(`/projects/${project.id}/tables/${table.id}/data`)">
                            <TableIcon class="w-4 h-4 mr-2" />
                            View Data
                          </DropdownMenuItem>
                          <DropdownMenuItem @click="handleEditTable(project.id, table.id)">
                            <Edit class="w-4 h-4 mr-2" />
                            Edit Schema
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            @click="handleDeleteTable(project.id, table.id, table.name)"
                            class="text-destructive focus:text-destructive"
                          >
                            <Trash2 class="w-4 h-4 mr-2" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </SidebarMenuSubItem>
                  <div v-if="!(projectTables.get(project.id) || []).length" class="px-2 py-2">
                    <p class="text-xs text-muted-foreground italic">No tables yet</p>
                  </div>
                </div>

                <!-- Workflows Section -->
                <div class="space-y-1 mt-3">
                  <div class="flex items-center justify-between px-2 py-1">
                    <span class="text-xs text-muted-foreground">Workflows</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      class="h-4 w-4 p-0"
                      @click="handleCreateWorkflow(project.id)"
                    >
                      <Plus class="w-3 h-3" />
                    </Button>
                  </div>
                  <SidebarMenuSubItem
                    v-for="workflow in projectWorkflows.get(project.id) || []"
                    :key="workflow.id"
                  >
                    <div class="group relative w-full">
                      <SidebarMenuSubButton
                        @click="router.push(`/projects/${project.id}/workflows/${workflow.id}/builder`)"
                        :isActive="isActive(`/projects/${project.id}/workflows/${workflow.id}`)"
                        class="pr-6"
                      >
                        <WorkflowIcon class="w-4 h-4" :class="workflow.is_active ? 'text-success' : 'text-muted-foreground'" />
                        <span class="flex-1 truncate">{{ workflow.name }}</span>
                      </SidebarMenuSubButton>

                      <!-- Workflow Actions -->
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            variant="ghost"
                            size="sm"
                            class="absolute right-0 top-0 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                          >
                            <MoreVertical class="w-3 h-3" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem @click="handleEditWorkflow(project.id, workflow.id)">
                            <Edit class="w-4 h-4 mr-2" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem @click="handleToggleWorkflow(project.id, workflow.id, workflow.is_active)">
                            <component :is="workflow.is_active ? Pause : Play" class="w-4 h-4 mr-2" />
                            {{ workflow.is_active ? 'Deactivate' : 'Activate' }}
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            @click="handleDeleteWorkflow(project.id, workflow.id, workflow.name)"
                            class="text-destructive focus:text-destructive"
                          >
                            <Trash2 class="w-4 h-4 mr-2" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </SidebarMenuSubItem>
                  <div v-if="!(projectWorkflows.get(project.id) || []).length" class="px-2 py-2">
                    <p class="text-xs text-muted-foreground italic">No workflows yet</p>
                  </div>
                </div>
              </SidebarMenuSub>
            </SidebarMenuItem>

            <!-- Empty State -->
            <div v-if="!projects.length && !isLoading" class="px-2 py-4">
              <p class="text-sm text-muted-foreground text-center">No projects yet</p>
              <Button variant="outline" size="sm" class="w-full mt-2" @click="handleCreateProject">
                <Plus class="w-4 h-4 mr-2" />
                Create Project
              </Button>
            </div>
          </SidebarMenu>
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
