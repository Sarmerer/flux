<script setup lang="ts">
import { Edit, MoreVertical, Plus, Zap } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useWorkflowStore } from '@/stores/workflows'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

import { useToast } from '@/composables/ui'

import { useRouteState } from './useRouteState'

const router = useRouter()
const activeProjectStore = useActiveProjectStore()
const workflowStore = useWorkflowStore()
const toast = useToast()
const { isActive } = useRouteState()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)

const handleCreateWorkflow = () => {
  if (!activeProjectId.value) {
    toast.warning('No Project Selected', 'Please select a project first')
    return
  }
  router.push(`/projects/${activeProjectId.value}/workflows?action=create`)
}

const handleViewWorkflow = (workflowId: string) => {
  router.push(`/projects/${activeProjectId.value}/workflows/${workflowId}/builder`)
}
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel class="flex items-center justify-between">
      <span class="text-xs font-semibold uppercase tracking-wider">Workflows</span>
      <Button
        variant="ghost"
        size="sm"
        class="h-5 w-5 p-0 opacity-60 hover:opacity-100 transition-opacity group-data-[collapsible=icon]:hidden"
        @click="handleCreateWorkflow"
      >
        <Plus class="w-3.5 h-3.5" />
      </Button>
    </SidebarGroupLabel>
    <SidebarGroupContent>
      <SidebarMenu>
        <div v-if="workflowStore.isLoading" class="px-3 py-2 text-xs text-muted-foreground">
          Loading workflows...
        </div>
        <SidebarMenuItem v-for="workflow in workflowStore.workflows" :key="workflow.id">
          <div class="group relative w-full">
            <SidebarMenuButton
              @click="handleViewWorkflow(workflow.id)"
              :isActive="isActive(`/projects/${activeProjectId}/workflows/${workflow.id}`)"
              class="pr-8"
            >
              <Zap class="w-4 h-4 flex-shrink-0" />
              <span class="flex-1 truncate">{{ workflow.name }}</span>
            </SidebarMenuButton>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  class="absolute right-1 top-1 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity group-data-[collapsible=icon]:hidden"
                >
                  <MoreVertical class="w-3.5 h-3.5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem @click="handleViewWorkflow(workflow.id)">
                  <Edit class="w-4 h-4 mr-2" />
                  Edit Workflow
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </SidebarMenuItem>
        <SidebarMenuItem v-if="!workflowStore.isLoading && !workflowStore.workflows.length">
          <div class="px-3">
            <Button
              variant="ghost"
              size="sm"
              class="w-full justify-start h-8 text-xs"
              @click="handleCreateWorkflow"
            >
              <Plus class="w-3.5 h-3.5 mr-2" />
              Create your first workflow
            </Button>
          </div>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
