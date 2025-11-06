<script setup lang="ts">
import {
  Check,
  ChevronDown,
  ChevronRight,
  FolderOpen,
  LogOut,
  Plus,
  Search,
  Settings as SettingsIcon,
  User,
  Workflow,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useAuthStore } from '@/stores/auth'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useProjects } from '@/composables/api'
import { useToast } from '@/composables/ui'
import { useBreadcrumbs } from '@/composables/useBreadcrumbs'
import { ROUTE_PATHS } from '@/constants/routes'

const router = useRouter()
const activeProjectStore = useActiveProjectStore()
const authStore = useAuthStore()
const { projects } = useProjects()
const toast = useToast()
const { breadcrumbs } = useBreadcrumbs()

const activeProject = computed(() => activeProjectStore.activeProject)
const showQuickActions = computed(() => !!activeProject.value)

const onProjectSelect = (projectId: string) => {
  router.push(ROUTE_PATHS.PROJECT_DETAIL(projectId))
}

const onQuickAction = (action: string) => {
  if (!activeProject.value) {
    toast.warning('No Project Selected', 'Please select a project first')
    return
  }

  switch (action) {
    case 'table':
      router.push(`${ROUTE_PATHS.PROJECT_TABLES(activeProject.value.id)}?action=create`)
      break
    case 'workflow':
      router.push(`${ROUTE_PATHS.PROJECT_WORKFLOWS(activeProject.value.id)}?action=create`)
      break
  }
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <header class="fixed top-0 left-0 right-0 z-50 flex h-14 items-center gap-4 border-b bg-background px-6">
    <div class="flex items-center gap-2.5">
      <div class="w-7 h-7 bg-gradient-to-br from-primary to-primary/80 rounded-lg flex items-center justify-center flex-shrink-0 shadow-sm">
        <Workflow class="w-3.5 h-3.5 text-primary-foreground" :stroke-width="2.5" />
      </div>
      <span class="font-semibold text-foreground tracking-tight">
        Flow
      </span>
    </div>

    <nav class="flex items-center gap-2 text-sm flex-1">
      <template v-for="(crumb, index) in breadcrumbs" :key="index">
        <ChevronRight v-if="index === 0 && breadcrumbs.length > 0" class="w-4 h-4 text-muted-foreground" />

        <DropdownMenu v-if="crumb.isProjectSelector">
          <DropdownMenuTrigger asChild>
            <button class="flex items-center gap-1 text-foreground font-medium hover:bg-muted px-2 py-1 rounded-md transition-colors">
              {{ crumb.label }}
              <ChevronDown class="w-3.5 h-3.5 text-muted-foreground" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" class="w-64">
            <DropdownMenuLabel class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Switch Project
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <div class="max-h-64 overflow-y-auto">
              <DropdownMenuItem
                v-for="project in projects"
                :key="project.id"
                @click="onProjectSelect(project.id)"
                class="cursor-pointer"
              >
                <div class="flex items-center justify-between w-full">
                  <div class="flex items-center space-x-2 flex-1 truncate">
                    <FolderOpen class="w-4 h-4 flex-shrink-0" />
                    <span class="truncate">{{ project.name }}</span>
                  </div>
                  <Check
                    v-if="activeProject?.id === project.id"
                    class="w-4 h-4 flex-shrink-0 text-primary"
                  />
                </div>
              </DropdownMenuItem>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="router.push('/projects')" class="cursor-pointer">
              <FolderOpen class="w-4 h-4 mr-2" />
              View All Projects
            </DropdownMenuItem>
            <DropdownMenuItem @click="router.push('/projects/new')" class="cursor-pointer font-medium">
              <Plus class="w-4 h-4 mr-2" />
              Create New Project
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <button
          v-else-if="crumb.path"
          @click="router.push(crumb.path)"
          class="text-muted-foreground hover:text-foreground transition-colors"
        >
          {{ crumb.label }}
        </button>
        <span v-else class="text-foreground font-medium">{{ crumb.label }}</span>
        <ChevronRight v-if="index < breadcrumbs.length - 1 && index > 0" class="w-4 h-4 text-muted-foreground" />
      </template>
    </nav>

    <div class="flex items-center gap-3">
      <Button
        variant="outline"
        size="sm"
        class="gap-2 w-[220px] justify-start text-muted-foreground hover:text-foreground transition-colors"
      >
        <Search class="w-3.5 h-3.5 flex-shrink-0" />
        <span class="flex-1 text-left text-sm">Search...</span>
        <kbd
          class="inline-flex h-5 select-none items-center gap-0.5 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium"
        >
          <span class="text-xs">⌘</span>K
        </kbd>
      </Button>

      <DropdownMenu v-if="showQuickActions">
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="sm" class="gap-2">
            <Plus class="w-4 h-4" />
            <span>Create</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>Quick Actions</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="onQuickAction('table')">
            <Plus class="w-4 h-4 mr-2" />
            New Table
          </DropdownMenuItem>
          <DropdownMenuItem @click="onQuickAction('workflow')">
            <Plus class="w-4 h-4 mr-2" />
            New Workflow
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" class="w-8 h-8 rounded-full">
            <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
              <User class="w-4 h-4 text-primary" />
            </div>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-56">
          <DropdownMenuItem @click="router.push('/settings')">
            <SettingsIcon class="w-4 h-4 mr-2" />
            Account Settings
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="handleLogout">
            <LogOut class="w-4 h-4 mr-2" />
            Logout
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </header>
</template>
