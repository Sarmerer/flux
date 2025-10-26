<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import { ExternalLink, Plus, Workflow } from 'lucide-vue-next'
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
          <div
            v-for="workflow in workflowsStore.workflows"
            :key="workflow.id"
            class="flex items-center justify-between p-4 border rounded-lg hover:bg-accent cursor-pointer"
            @click="
              router.push(`/projects/${props.activeProject.id}/workflows/${workflow.id}/builder`)
            "
          >
            <div class="flex items-center space-x-3">
              <div class="p-2 bg-purple-100 rounded-lg">
                <Workflow class="h-5 w-5 text-purple-600" />
              </div>
              <div>
                <h4 class="font-medium">{{ workflow.name }}</h4>
                <p class="text-sm text-muted-foreground">
                  {{ workflow.description || 'No description' }}
                </p>
              </div>
            </div>
            <ExternalLink class="h-4 w-4 text-muted-foreground" />
          </div>
        </div>
      </LoadingWrapper>
    </CardContent>
  </Card>
</template>
