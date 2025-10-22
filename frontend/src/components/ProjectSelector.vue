<script setup lang="ts">
import { Check, ChevronsUpDown, FolderOpen, Plus } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { projectService } from '@/api/services/project'
import type { Project } from '@/types/api'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

import { useResourceCache } from '@/composables/data'

const router = useRouter()
const route = useRoute()

const { data: projects, loading } = useResourceCache<Project>('projects-list', {
  fetchFn: () => projectService.getAll(),
  subscribeToUpdates: true,
  events: ['project:created', 'project:updated', 'project:deleted'],
  ttlMs: 60000,
})

const currentProjectId = computed(() => route.params.projectId as string)

const currentProject = computed(() => {
  if (!currentProjectId.value || !projects.value) return null
  return projects.value.find((p: Project) => p.id === currentProjectId.value)
})

const handleSelectProject = (projectId: string) => {
  router.push(`/projects/${projectId}/tables`)
}

const handleCreateProject = () => {
  router.push('/projects?action=create')
}

const handleViewAllProjects = () => {
  router.push('/projects')
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger asChild>
      <Button variant="outline" class="w-full justify-between" :disabled="loading">
        <div class="flex items-center space-x-2 truncate">
          <FolderOpen class="w-4 h-4 flex-shrink-0" />
          <span class="truncate">
            {{ currentProject?.name || 'Select Project' }}
          </span>
        </div>
        <ChevronsUpDown class="w-4 h-4 ml-2 flex-shrink-0 opacity-50" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="start" class="w-64">
      <DropdownMenuLabel>Projects</DropdownMenuLabel>
      <DropdownMenuSeparator />

      <!-- Project List -->
      <div class="max-h-64 overflow-y-auto">
        <DropdownMenuItem
          v-for="project in projects"
          :key="project.id"
          @click="handleSelectProject(project.id)"
          class="cursor-pointer"
        >
          <div class="flex items-center justify-between w-full">
            <div class="flex items-center space-x-2 flex-1 truncate">
              <FolderOpen class="w-4 h-4 flex-shrink-0" />
              <span class="truncate">{{ project.name }}</span>
            </div>
            <Check v-if="currentProjectId === project.id" class="w-4 h-4 flex-shrink-0" />
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

      <!-- Actions -->
      <DropdownMenuItem @click="handleViewAllProjects" class="cursor-pointer">
        <FolderOpen class="w-4 h-4 mr-2" />
        View All Projects
      </DropdownMenuItem>
      <DropdownMenuItem @click="handleCreateProject" class="cursor-pointer">
        <Plus class="w-4 h-4 mr-2" />
        Create New Project
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
