<script setup lang="ts">
import LoadingWrapper from '@/components/common/LoadingWrapper.vue'
import type { CreateWorkflowData } from '@/components/workflows/drawers/WorkflowEditorContent.vue'
import { CheckCircle, Plus, Search, Workflow as WorkflowIcon, Zap } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import type { Workflow } from '@/types'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import { useTables } from '@/composables/api/useTables'
import { useWorkflows } from '@/composables/api/useWorkflows'
import { useDrawers } from '@/composables/drawerRegistry'
import { useRouteContext } from '@/composables/routing'
import { useToast } from '@/composables/ui'

const router = useRouter()
const toast = useToast()
const { openWorkflowEditor } = useDrawers()

const { projectId } = useRouteContext()
const { workflows, loading, createWorkflow } = useWorkflows(projectId.value)
const { tables, loading: loadingTables } = useTables(projectId.value)

const searchQuery = ref('')

const filteredWorkflows = computed(() => {
  if (!searchQuery.value) return workflows.value
  return workflows.value.filter(
    (workflow: Workflow) =>
      workflow.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (workflow.description &&
        workflow.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )
})

const openCreateWorkflowDrawer = () => {
  openWorkflowEditor({
    props: {
      tables: tables.value,
      loadingTables: loadingTables.value,
    },
    onSave: async (data: CreateWorkflowData) => {
      try {
        const workflow = await createWorkflow(data)
        toast.success('Workflow Created', 'Your workflow has been created successfully')
        router.push(`/projects/${projectId.value}/workflows/${workflow.id}`)
      } catch (error) {
        console.error('Failed to create workflow:', error)
        toast.error('Creation Failed', 'Failed to create workflow. Please try again.')
      }
    },
  })
}

const handleWorkflowSelect = (id: string) => {
  router.push(`/projects/${projectId.value}/workflows/${id}`)
}
</script>

<template>
  <div class="flex h-[calc(100vh-3.5rem)]">
    <div class="w-80 border-r flex flex-col bg-muted/10">
      <div class="p-4 border-b space-y-3">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Workflows</h2>
          <Button size="sm" @click="openCreateWorkflowDrawer">
            <Plus class="w-4 h-4" />
          </Button>
        </div>
        <div class="relative">
          <Search
            class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4"
          />
          <Input v-model="searchQuery" placeholder="Search workflows..." class="pl-9 h-9" />
        </div>
      </div>

      <LoadingWrapper :is-loading="loading" loading-text="Loading workflows...">
        <div class="flex-1 overflow-y-auto">
          <div class="p-2 space-y-1">
            <button
              v-for="workflow in filteredWorkflows"
              :key="workflow.id"
              @click="handleWorkflowSelect(workflow.id)"
              class="w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors hover:bg-muted text-muted-foreground hover:text-foreground"
            >
              <Zap class="w-4 h-4 flex-shrink-0" />
              <span class="flex-1 text-left truncate">{{ workflow.name }}</span>
              <CheckCircle
                v-if="workflow.is_active"
                class="w-3.5 h-3.5 text-green-600 flex-shrink-0"
              />
            </button>
          </div>

          <div
            v-if="filteredWorkflows.length === 0"
            class="p-4 text-center text-sm text-muted-foreground"
          >
            <p>No workflows found</p>
          </div>
        </div>
      </LoadingWrapper>
    </div>

    <div class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-4">
        <div class="flex justify-center">
          <div class="w-16 h-16 rounded-full bg-muted flex items-center justify-center">
            <WorkflowIcon class="w-8 h-8 text-muted-foreground" />
          </div>
        </div>
        <div>
          <h3 class="text-lg font-semibold mb-1">Select a workflow</h3>
          <p class="text-sm text-muted-foreground">
            Choose a workflow from the list to view and edit it
          </p>
        </div>
        <Button @click="openCreateWorkflowDrawer">
          <Plus class="w-4 h-4 mr-2" />
          Create New Workflow
        </Button>
      </div>
    </div>
  </div>
</template>
