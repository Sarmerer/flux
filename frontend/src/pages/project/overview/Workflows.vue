<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import ResourceListItem from '@/components/project/ResourceListItem.vue'
import { Plus, Workflow } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

import { useWorkflowStore } from '@/stores/workflows'
import type { Project } from '@/types/api'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const router = useRouter()
const props = defineProps<{ activeProject: Project }>()

const workflowsStore = useWorkflowStore()

const onWorkflowCreateClick = () => {
  router.push(`/projects/${props.activeProject.id}/workflows`)
}
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-center justify-between">
        <div>
          <CardTitle>Workflows</CardTitle>
          <CardDescription>Automate your database operations</CardDescription>
        </div>
        <Button @click="onWorkflowCreateClick">
          <Plus class="w-4 h-4 mr-2" />
          New Workflow
        </Button>
      </div>
    </CardHeader>
    <CardContent>
      <LoadingWrapper :is-loading="workflowsStore.isLoading" loading-text="Loading workflows...">
        <ErrorState
          v-if="workflowsStore.error"
          :error="workflowsStore.error"
          title="Failed to load workflows"
          :on-retry="() => workflowsStore.loadByProjectId(props.activeProject.id)"
        />
        <div v-else class="space-y-4">
          <ResourceListItem
            v-for="workflow in workflowsStore.workflows"
            :key="workflow.id"
            :name="workflow.name"
            :description="workflow.description"
            :icon="Workflow"
            icon-bg-color="bg-purple-100"
            icon-color="text-purple-600"
            @click="
              router.push(`/projects/${props.activeProject.id}/workflows/${workflow.id}/builder`)
            "
          />
        </div>
      </LoadingWrapper>
    </CardContent>
  </Card>
</template>
