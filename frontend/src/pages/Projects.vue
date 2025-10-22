<script setup lang="ts">
import ApiErrorBoundary from '@/components/ApiErrorBoundary.vue'
import { Calendar, Edit, FolderOpen, Plus, Search, Settings, Trash2 } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

import { useProjects } from '@/composables/api'
import { usePermissions } from '@/composables/usePermissions'

const route = useRoute()
const router = useRouter()

const { projects, loading, error, refresh, createProject, deleteProject } = useProjects()
const { can } = usePermissions()

const searchQuery = ref('')
const isCreateDialogOpen = ref(false)
const isCreating = ref(false)
const newProject = ref({
  name: '',
  description: '',
})

const filteredProjects = computed(() => {
  if (!projects.value || !Array.isArray(projects.value)) return []
  if (!searchQuery.value) return projects.value
  return projects.value.filter(
    (project) =>
      project?.name?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (project?.description &&
        project.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )
})

const handleCreateProject = async () => {
  if (!newProject.value.name.trim()) return

  isCreating.value = true
  try {
    const project = await createProject({
      name: newProject.value.name,
      description: newProject.value.description || undefined,
    })
    isCreateDialogOpen.value = false
    newProject.value = { name: '', description: '' }

    router.push(`/projects/${project.id}`)
  } catch (error) {
    console.error('Failed to create project:', error)
    alert('Failed to create project. Please try again.')
  } finally {
    isCreating.value = false
  }
}

const handleProjectClick = (projectId: string) => {
  router.push(`/projects/${projectId}`)
}

const handleEditProject = (projectId: string, event: Event) => {
  event.stopPropagation()
}

const handleDeleteProject = async (projectId: string, event: Event) => {
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

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

// Watch route to open modal when navigating to /projects/new
watch(
  () => route.meta.openCreateModal,
  (shouldOpen) => {
    if (shouldOpen) {
      isCreateDialogOpen.value = true
    }
  },
  { immediate: true }
)

// Watch modal state to update route when closing
watch(isCreateDialogOpen, (isOpen) => {
  if (!isOpen && route.path === '/projects/new') {
    router.push('/projects')
  }
})
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Projects</h1>
        <p class="text-gray-600 dark:text-gray-400">Manage your database projects</p>
      </div>
      <!-- Only show create button if user has permission -->
      <Button
        v-if="can('projects.create')"
        class="flex items-center space-x-2"
        @click="router.push('/projects/new')"
      >
        <Plus class="w-4 h-4" />
        <span>New Project</span>
      </Button>
      <Dialog v-model:open="isCreateDialogOpen">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create New Project</DialogTitle>
            <DialogDescription>
              Create a new project to start building your database schema and workflows.
            </DialogDescription>
          </DialogHeader>
          <div class="space-y-4">
            <div class="space-y-2">
              <Label for="project-name">Project Name</Label>
              <Input
                id="project-name"
                v-model="newProject.name"
                placeholder="Enter project name"
                required
              />
            </div>
            <div class="space-y-2">
              <Label for="project-description">Description (Optional)</Label>
              <Textarea
                id="project-description"
                v-model="newProject.description"
                placeholder="Enter project description"
                rows="3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" @click="isCreateDialogOpen = false" :disabled="isCreating">
              Cancel
            </Button>
            <Button @click="handleCreateProject" :disabled="!newProject.name.trim() || isCreating">
              {{ isCreating ? 'Creating...' : 'Create Project' }}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>

    <!-- Search -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
      <Input v-model="searchQuery" placeholder="Search projects..." class="pl-10" />
    </div>

    <!-- Projects Grid -->
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
        <Button @click="isCreateDialogOpen = true">Create Project</Button>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card
        v-for="project in filteredProjects"
        :key="project.id"
        class="cursor-pointer hover:shadow-lg dark:hover:shadow-gray-800 transition-shadow duration-200"
        @click="handleProjectClick(project.id)"
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
              <!-- Only show edit button if user has permission -->
              <Button
                v-if="can('projects.edit')"
                variant="ghost"
                size="sm"
                @click="handleEditProject(project.id, $event)"
              >
                <Edit class="h-4 w-4" />
              </Button>
              <!-- Only show delete button if user has permission -->
              <Button
                v-if="can('projects.delete')"
                variant="ghost"
                size="sm"
                @click="handleDeleteProject(project.id, $event)"
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
            <Button variant="outline" size="sm">
              <Settings class="h-4 w-4 mr-1" />
              Settings
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
