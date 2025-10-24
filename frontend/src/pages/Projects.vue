<script setup lang="ts">
import ApiErrorBoundary from '@/components/ApiErrorBoundary.vue'
import { Calendar, FolderOpen, Plus, Search, Settings, Trash2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { useProjects } from '@/composables/api'
import { useFormatting } from '@/composables/formatting'
import { usePermissions } from '@/composables/usePermissions'

const router = useRouter()
const { formatDate } = useFormatting()

const { projects, loading, error, refresh, deleteProject } = useProjects()
const { can } = usePermissions()

const searchQuery = ref('')

const filteredProjects = computed(() => {
  const list = projects.value
  const query = searchQuery.value?.toLowerCase().trim()

  if (!list?.length) return []
  if (!query) return list

  return list.filter(({ name = '', description = '' }) => {
    return name.toLowerCase().includes(query) || description.toLowerCase().includes(query)
  })
})

const onProjectClick = (projectId: string) => {
  router.push(`/projects/${projectId}`)
}

const onProjectDelete = async (projectId: string, event: Event) => {
  event.stopPropagation()
  if (confirm('Are you sure you want to delete this project? This action cannot be undone.')) {
    try {
      await deleteProject(projectId)
    } catch (error) {
      console.error('Failed to delete project:', error)
      alert('Failed to delete project. Please try again.')
    }
  }
}
</script>

<template>
  <div class="p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Projects</h1>
        <p class="text-gray-600 dark:text-gray-400">Manage your database projects</p>
      </div>
      <Button
        v-if="can('projects.create')"
        class="flex items-center space-x-2"
        @click="router.push('/projects/new')"
      >
        <Plus class="w-4 h-4" />
        <span>New Project</span>
      </Button>
    </div>

    <div class="relative">
      <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
      <Input v-model="searchQuery" placeholder="Search projects..." class="pl-10" />
    </div>

    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <ApiErrorBoundary v-else-if="error" :error="error" :retry="refresh" />

    <div v-else-if="filteredProjects.length === 0" class="text-center py-12">
      <FolderOpen class="mx-auto h-12 w-12 text-gray-400 dark:text-gray-600" />
      <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-gray-100">
        {{ searchQuery ? 'No projects found' : 'No projects yet' }}
      </h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        {{
          searchQuery
            ? 'Try adjusting your search terms.'
            : 'Get started by creating a new project.'
        }}
      </p>
      <div v-if="!searchQuery" class="mt-6">
        <Button @click="router.push('/projects/new')">Create Project</Button>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card
        v-for="project in filteredProjects"
        :key="project.id"
        class="cursor-pointer hover:shadow-lg dark:hover:shadow-gray-800 transition-shadow duration-200"
        @click="onProjectClick(project.id)"
      >
        <CardHeader class="pb-3">
          <div class="flex items-start justify-between">
            <div class="flex items-center space-x-3">
              <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                <FolderOpen class="h-5 w-5 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <CardTitle class="text-lg">{{ project.name }}</CardTitle>
                <CardDescription class="mt-1">
                  {{ project.description || 'No description' }}
                </CardDescription>
              </div>
            </div>
            <div class="flex items-center space-x-1">
              <Button v-if="can('projects.edit')" variant="ghost" size="sm">
                <Settings class="h-4 w-4 mr-1" />
              </Button>
              <Button
                v-if="can('projects.delete')"
                variant="ghost"
                size="sm"
                @click="onProjectDelete(project.id, $event)"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent class="pt-0">
          <div class="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
            <div class="flex items-center space-x-1">
              <Calendar class="h-4 w-4" />
              <span>Created {{ formatDate(project.created_at) }}</span>
            </div>
            <Badge variant="secondary">Active</Badge>
          </div>
          <div class="mt-4 flex items-center justify-between">
            <div class="flex space-x-4 text-sm text-gray-500 dark:text-gray-400">
              <span>0 Tables</span>
              <span>0 Workflows</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
