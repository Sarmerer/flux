<script setup lang="ts">
import { Check, ChevronsUpDown, FolderOpen, Plus } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'

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

const router = useRouter()
const { projects, loading } = useProjects()
const activeProjectStore = useActiveProjectStore()

const onProjectSelect = (projectId: string) => {
  router.push(`/projects/${projectId}`)
}

const onProjectCreateClick = () => {
  router.push('/projects/new')
}

const onViewAllProjectsClick = () => {
  router.push('/projects')
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger asChild>
      <Button variant="outline" class="w-full justify-between h-10" :disabled="loading">
        <div class="flex items-center space-x-2 truncate">
          <FolderOpen class="w-4 h-4 flex-shrink-0" />
          <span class="truncate text-sm">
            {{ activeProjectStore.activeProject?.name || 'Select a project' }}
          </span>
        </div>
        <ChevronsUpDown class="w-3.5 h-3.5 ml-2 flex-shrink-0 opacity-50" />
      </Button>
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
              v-if="activeProjectStore.activeProject?.id === project.id"
              class="w-4 h-4 flex-shrink-0 text-primary"
            />
          </div>
        </DropdownMenuItem>

        <div
          v-if="!projects?.length && !loading"
          class="px-2 py-3 text-sm text-muted-foreground text-center"
        >
          No projects yet
        </div>
      </div>

      <DropdownMenuSeparator />

      <DropdownMenuItem @click="onViewAllProjectsClick" class="cursor-pointer">
        <FolderOpen class="w-4 h-4 mr-2" />
        View All Projects
      </DropdownMenuItem>
      <DropdownMenuItem @click="onProjectCreateClick" class="cursor-pointer font-medium">
        <Plus class="w-4 h-4 mr-2" />
        Create New Project
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
