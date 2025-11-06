<script setup lang="ts">
import LoadingWrapper from '@/components/common/LoadingWrapper.vue'
import {
  CheckCircle,
  Plus,
  Search,
  Settings2,
  Workflow as WorkflowIcon,
  Zap,
} from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'

import { useRouteContext } from '@/composables/routing'

const router = useRouter()

const { projectId, workflowId } = useRouteContext()
const searchQuery = ref('')
const isCreateDialogOpen = ref(false)
const isLoading = ref(false)

const workflows = ref([
  {
    id: '1',
    name: 'New User Welcome',
    description: 'Send welcome email when a new user is created',
    project_id: projectId.value,
    trigger: {
      type: 'on_row_created' as const,
      table_name: 'users',
      conditions: {},
    },
    actions: [
      {
        id: '1',
        type: 'send_email' as const,
        config: {
          template: 'welcome',
          to: '{{user.email}}',
        },
        order: 1,
      },
    ],
    is_active: true,
    last_run: '2024-01-20T10:30:00Z',
    created_at: '2024-01-15T10:30:00Z',
    updated_at: '2024-01-20T10:30:00Z',
  },
  {
    id: '2',
    name: 'Order Processing',
    description: 'Process orders and update inventory',
    project_id: projectId.value,
    trigger: {
      type: 'on_row_created' as const,
      table_name: 'orders',
      conditions: {},
    },
    actions: [
      {
        id: '1',
        type: 'update_row' as const,
        config: {
          table: 'inventory',
          condition: 'product_id = {{order.product_id}}',
          updates: {
            quantity: 'quantity - {{order.quantity}}',
          },
        },
        order: 1,
      },
      {
        id: '2',
        type: 'send_webhook' as const,
        config: {
          url: 'https://api.example.com/order-processed',
          method: 'POST',
          payload: {
            order_id: '{{order.id}}',
            status: 'processed',
          },
        },
        order: 2,
      },
    ],
    is_active: true,
    last_run: '2024-01-21T14:22:00Z',
    created_at: '2024-01-16T09:15:00Z',
    updated_at: '2024-01-21T14:22:00Z',
  },
  {
    id: '3',
    name: 'Data Backup',
    description: 'Daily backup of critical data',
    project_id: projectId.value,
    trigger: {
      type: 'scheduled' as const,
      schedule: '0 2 * * *',
    },
    actions: [
      {
        id: '1',
        type: 'send_webhook' as const,
        config: {
          url: 'https://backup.example.com/backup',
          method: 'POST',
        },
        order: 1,
      },
    ],
    is_active: false,
    last_run: '2024-01-21T02:00:00Z',
    created_at: '2024-01-17T11:20:00Z',
    updated_at: '2024-01-21T02:00:00Z',
  },
])

const newWorkflow = ref({
  name: '',
  description: '',
  trigger: {
    type: 'on_row_created' as const,
    table_name: '',
    conditions: {},
  },
  actions: [] as any[],
})

const triggerTypes = [
  { value: 'on_row_created', label: 'On Row Created' },
  { value: 'on_row_updated', label: 'On Row Updated' },
  { value: 'on_row_deleted', label: 'On Row Deleted' },
  { value: 'scheduled', label: 'Scheduled' },
  { value: 'webhook', label: 'Webhook' },
]

const filteredWorkflows = computed(() => {
  if (!searchQuery.value) return workflows.value
  return workflows.value.filter(
    (workflow) =>
      workflow.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (workflow.description &&
        workflow.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )
})

const selectedWorkflow = computed(() => {
  if (!workflowId.value) return null
  return workflows.value.find((w) => w.id === workflowId.value)
})

const handleCreateWorkflow = async () => {
  if (!newWorkflow.value.name.trim()) return

  try {
    const workflow = {
      id: Date.now().toString(),
      name: newWorkflow.value.name,
      description: newWorkflow.value.description,
      project_id: projectId.value,
      trigger: newWorkflow.value.trigger,
      actions: newWorkflow.value.actions,
      is_active: false,
      last_run: '2024-01-01T00:00:00Z',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    workflows.value.push(workflow)
    isCreateDialogOpen.value = false
    newWorkflow.value = {
      name: '',
      description: '',
      trigger: {
        type: 'on_row_created',
        table_name: '',
        conditions: {},
      },
      actions: [],
    }

    router.push(`/projects/${projectId.value}/workflows/${workflow.id}/builder`)
  } catch (error) {
    console.error('Failed to create workflow:', error)
  }
}

const handleSelectWorkflow = (id: string) => {
  router.push(`/projects/${projectId.value}/workflows/${id}`)
}
</script>

<template>
  <div class="flex h-[calc(100vh-3.5rem)]">
    <div class="w-80 border-r flex flex-col bg-muted/10">
      <div class="p-4 border-b space-y-3">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Workflows</h2>
          <Button size="sm" @click="isCreateDialogOpen = true">
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

      <LoadingWrapper :is-loading="isLoading" loading-text="Loading workflows...">
        <div class="flex-1 overflow-y-auto">
          <div class="p-2 space-y-1">
            <button
              v-for="workflow in filteredWorkflows"
              :key="workflow.id"
              @click="handleSelectWorkflow(workflow.id)"
              :class="[
                'w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors',
                workflowId === workflow.id
                  ? 'bg-primary/10 text-primary font-medium'
                  : 'hover:bg-muted text-muted-foreground hover:text-foreground',
              ]"
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

    <div class="flex-1 flex flex-col">
      <div v-if="!workflowId" class="flex-1 flex items-center justify-center">
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
          <Button @click="isCreateDialogOpen = true">
            <Plus class="w-4 h-4 mr-2" />
            Create New Workflow
          </Button>
        </div>
      </div>

      <div v-else class="flex-1 flex flex-col">
        <div class="border-b px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div>
                <h1 class="text-2xl font-bold">{{ selectedWorkflow?.name }}</h1>
                <p class="text-sm text-muted-foreground mt-1">
                  {{ selectedWorkflow?.description || 'No description' }}
                </p>
              </div>
              <Badge v-if="selectedWorkflow?.is_active" variant="default">Active</Badge>
              <Badge v-else variant="secondary">Inactive</Badge>
            </div>
            <div class="flex items-center gap-2">
              <Button variant="outline" size="sm">
                <Settings2 class="w-4 h-4 mr-2" />
                Settings
              </Button>
            </div>
          </div>
        </div>

        <div class="flex-1 overflow-auto p-6">
          <div class="space-y-4">
            <div class="text-sm text-muted-foreground">Workflow builder view would go here</div>
          </div>
        </div>
      </div>
    </div>

    <Dialog v-model:open="isCreateDialogOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New Workflow</DialogTitle>
          <DialogDescription>
            Create a new workflow to automate your database operations.
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div class="space-y-2">
            <Label for="workflow-name">Workflow Name</Label>
            <Input
              id="workflow-name"
              v-model="newWorkflow.name"
              placeholder="Enter workflow name"
              required
            />
          </div>
          <div class="space-y-2">
            <Label for="workflow-description">Description</Label>
            <Textarea
              id="workflow-description"
              v-model="newWorkflow.description"
              placeholder="Enter workflow description"
              rows="3"
            />
          </div>
          <div class="space-y-2">
            <Label for="trigger-type">Trigger Type</Label>
            <Select v-model="newWorkflow.trigger.type">
              <SelectTrigger>
                <SelectValue placeholder="Select trigger type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="trigger in triggerTypes"
                  :key="trigger.value"
                  :value="trigger.value"
                >
                  {{ trigger.label }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="isCreateDialogOpen = false"> Cancel </Button>
          <Button @click="handleCreateWorkflow" :disabled="!newWorkflow.name.trim()">
            Create Workflow
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
