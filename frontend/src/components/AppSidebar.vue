<script setup lang="ts">
import ProjectSelector from '@/components/ProjectSelector.vue'
import {
  Edit,
  Eye,
  FolderOpen,
  Home,
  LogOut,
  MoreVertical,
  Plus,
  Settings as SettingsIcon,
  Table as TableIcon,
  User,
  Users,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useAuthStore } from '@/stores/auth'
import { useTableStore } from '@/stores/tables'
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
const authStore = useAuthStore()
const activeProjectStore = useActiveProjectStore()
const tableStore = useTableStore()
const workflowStore = useWorkflowStore()
const toast = useToast()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)

const isActive = (itemUrl: string) => {
  return route.path === itemUrl || route.path.startsWith(itemUrl + '/')
}

const onTableCreateClick = () => {
  if (!activeProjectId.value) {
    toast.warning('No Project Selected', 'Please select a project first')
    return
  }
  router.push(`/projects/${activeProjectId.value}/tables?action=create`)
}

const handleViewTable = (tableId: string) => {
  router.push(`/projects/${activeProjectId.value}/tables/${tableId}/data`)
}

const onTableEditClick = (tableId: string) => {
  router.push(`/projects/${activeProjectId.value}/tables/${tableId}/builder`)
}

const onLogoutClick = async () => {
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

      <ProjectSelector />
    </SidebarHeader>

    <SidebarContent>
      <!-- Main Navigation -->
      <SidebarGroup>
        <SidebarGroupLabel>Navigation</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton @click="router.push('/')" :isActive="isActive('/')">
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

      <SidebarGroup v-if="activeProjectId">
        <SidebarGroupLabel>Project</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                @click="router.push(`/projects/${activeProjectId}/members`)"
                :isActive="isActive(`/projects/${activeProjectId}/members`)"
              >
                <Users class="w-4 h-4" />
                <span>Members</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- Tables Section -->
      <SidebarGroup v-if="activeProjectId">
        <SidebarGroupLabel class="flex items-center justify-between">
          <span>Tables</span>
          <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="onTableCreateClick">
            <Plus class="w-4 h-4" />
          </Button>
        </SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <div v-if="tableStore.isLoading" class="px-2 py-2 text-sm text-muted-foreground">
              Loading...
            </div>
            <SidebarMenuItem v-for="table in tableStore.tables" :key="table.id">
              <div class="group relative w-full">
                <SidebarMenuButton
                  @click="handleViewTable(table.id)"
                  :isActive="isActive(`/projects/${activeProjectId}/tables/${table.id}`)"
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
                    <DropdownMenuItem @click="onTableEditClick(table.id)">
                      <Edit class="w-4 h-4 mr-2" />
                      Edit Schema
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </SidebarMenuItem>
            <div v-if="!tableStore.isLoading && !tableStore.tables.length" class="px-2 py-2">
              <p class="text-xs text-muted-foreground italic">No tables yet</p>
            </div>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- Workflows Section -->
      <SidebarGroup v-if="activeProjectId">
        <SidebarGroupLabel class="flex items-center justify-between">
          <span>Workflows</span>
        </SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <div v-if="workflowStore.isLoading" class="px-2 py-2 text-sm text-muted-foreground">
              Loading...
            </div>
            <SidebarMenuItem v-for="workflow in workflowStore.workflows" :key="workflow.id">
              <div class="group relative w-full"></div>
            </SidebarMenuItem>
            <div
              v-if="!workflowStore.isLoading && !workflowStore.workflows.length"
              class="px-2 py-2"
            >
              <p class="text-xs text-muted-foreground italic">No workflows yet</p>
            </div>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- No Project Selected State -->
      <SidebarGroup v-if="!activeProjectId">
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
        <Button variant="ghost" size="sm" class="w-full justify-start" @click="onLogoutClick">
          <LogOut class="w-4 h-4 mr-2" />
          Logout
        </Button>
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
