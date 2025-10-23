<script setup lang="ts">
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
</script>

<template>
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
</template>
