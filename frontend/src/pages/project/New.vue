<script setup lang="ts">
import { FolderOpen, Loader2 } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

import { useProjects } from '@/composables/api'

const router = useRouter()
const { createProject } = useProjects()

const newProject = ref({
  name: '',
  description: '',
})

const isCreating = ref(false)
const error = ref<string | null>(null)

const onCreateProjectClick = async () => {
  if (!newProject.value.name.trim()) return

  isCreating.value = true
  error.value = null

  try {
    const project = await createProject({
      name: newProject.value.name,
      description: newProject.value.description || undefined,
    })

    router.push(`/projects/${project.id}`)
  } catch (err) {
    console.error('Failed to create project:', err)
    error.value = err instanceof Error ? err.message : 'Failed to create project. Please try again.'
  } finally {
    isCreating.value = false
  }
}

const onCancelClick = () => {
  router.back()
}

const onKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && (event.metaKey || event.ctrlKey)) {
    onCreateProjectClick()
  }
}
</script>

<template>
  <div class="p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Create New Project</h1>
        <p class="text-gray-600 dark:text-gray-400">
          Set up a new project to organize your database schemas, tables, and workflows
        </p>
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Project Details</CardTitle>
        <CardDescription>
          Give your project a name and description to help you identify it later
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-6">
        <div class="space-y-2">
          <Label for="project-name">
            Project Name
            <span class="text-red-500">*</span>
          </Label>
          <Input
            id="project-name"
            v-model="newProject.name"
            placeholder="e.g., E-commerce Platform, Customer Portal"
            :disabled="isCreating"
            required
            autofocus
            @keydown="onKeyDown"
          />
        </div>

        <div class="space-y-2">
          <Label for="project-description">
            Description
            <span class="text-gray-400 font-normal">(Optional)</span>
          </Label>
          <Textarea
            id="project-description"
            v-model="newProject.description"
            placeholder="Describe what this project is for, its goals, or any other relevant information..."
            rows="4"
            class="resize-none"
            :disabled="isCreating"
            @keydown="onKeyDown"
          />
        </div>

        <div
          v-if="error"
          class="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
        >
          <p class="text-sm text-red-600 dark:text-red-400 font-medium">{{ error }}</p>
        </div>

        <div class="flex items-center justify-end space-x-3 pt-4">
          <Button variant="outline" @click="onCancelClick" :disabled="isCreating"> Cancel </Button>
          <Button @click="onCreateProjectClick" :disabled="!newProject.name.trim() || isCreating">
            <Loader2 v-if="isCreating" class="h-4 w-4 mr-2 animate-spin" />
            <FolderOpen v-else class="h-4 w-4 mr-2" />
            {{ isCreating ? 'Creating...' : 'Create Project' }}
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
