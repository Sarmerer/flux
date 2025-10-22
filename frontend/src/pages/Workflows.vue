<script setup lang="ts">
import {
  CheckCircle,
  Clock,
  Edit,
  Play,
  Plus,
  Search,
  Trash2,
  Workflow,
  XCircle,
} from 'lucide-vue-next'
import { computed, ref } from 'vue'
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
  DialogTrigger,
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

const route = useRoute()
const router = useRouter()

const projectId = computed(() => route.params.projectId as string)
const searchQuery = ref('')
const isCreateDialogOpen = ref(false)
const isEditDialogOpen = ref(false)
const selectedWorkflow = ref(null)

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

const handleEditWorkflow = (workflow: any) => {
  selectedWorkflow.value = workflow
  isEditDialogOpen.value = true
}

const handleDeleteWorkflow = async (workflowId: string) => {
  if (confirm('Are you sure you want to delete this workflow? This action cannot be undone.')) {
    try {
      workflows.value = workflows.value.filter((w) => w.id !== workflowId)
    } catch (error) {
      console.error('Failed to delete workflow:', error)
    }
  }
}

const handleToggleWorkflow = async (workflow: any) => {
  try {
    workflow.is_active = !workflow.is_active
  } catch (error) {
    console.error('Failed to toggle workflow:', error)
  }
}

const handleRunWorkflow = async (workflow: any) => {
  try {
    console.log('Running workflow:', workflow.id)
  } catch (error) {
    console.error('Failed to run workflow:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const getStatusIcon = (workflow: any) => {
  if (workflow.is_active) {
    return CheckCircle
  } else {
    return XCircle
  }
}

const getStatusColor = (workflow: any) => {
  if (workflow.is_active) {
    return 'text-green-600'
  } else {
    return 'text-gray-400'
  }
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Workflows</h1>
        <p class="text-gray-600">Automate your database operations with workflows</p>
      </div>
      <Dialog v-model:open="isCreateDialogOpen">
        <DialogTrigger asChild>
          <Button class="flex items-center space-x-2">
            <Plus class="w-4 h-4" />
            <span>New Workflow</span>
          </Button>
        </DialogTrigger>
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

    <!-- Search -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
      <Input v-model="searchQuery" placeholder="Search workflows..." class="pl-10" />
    </div>

    <!-- Workflows Grid -->
    <div v-if="filteredWorkflows.length === 0" class="text-center py-12">
      <Workflow class="mx-auto h-12 w-12 text-gray-400" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">
        {{ searchQuery ? 'No workflows found' : 'No workflows yet' }}
      </h3>
      <p class="mt-1 text-sm text-gray-500">
        {{
          searchQuery
            ? 'Try adjusting your search terms.'
            : 'Get started by creating your first workflow.'
        }}
      </p>
      <div v-if="!searchQuery" class="mt-6">
        <Button @click="isCreateDialogOpen = true">Create Workflow</Button>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card
        v-for="workflow in filteredWorkflows"
        :key="workflow.id"
        class="hover:shadow-lg transition-shadow duration-200 flex flex-col justify-between"
      >
        <CardHeader class="pb-3">
          <div class="flex items-start justify-between">
            <div class="flex items-center space-x-3">
              <div class="p-2 bg-purple-100 rounded-lg">
                <Workflow class="h-5 w-5 text-purple-600" />
              </div>
              <div>
                <CardTitle class="text-lg">{{ workflow.name }}</CardTitle>
                <CardDescription class="mt-1">
                  {{ workflow.description || 'No description' }}
                </CardDescription>
              </div>
            </div>
            <div class="flex items-center space-x-1">
              <Button
                variant="ghost"
                size="sm"
                @click="handleRunWorkflow(workflow)"
                title="Run Workflow"
              >
                <Play class="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="handleToggleWorkflow(workflow)"
                :title="workflow.is_active ? 'Pause Workflow' : 'Activate Workflow'"
              >
                <component
                  :is="getStatusIcon(workflow)"
                  :class="getStatusColor(workflow)"
                  class="h-4 w-4"
                />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="handleEditWorkflow(workflow)"
                title="Edit Workflow"
              >
                <Edit class="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="handleDeleteWorkflow(workflow.id)"
                title="Delete Workflow"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent class="pt-0">
          <div class="space-y-3">
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Trigger</span>
              <Badge variant="outline">
                {{ triggerTypes.find((t) => t.value === workflow.trigger.type)?.label }}
              </Badge>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Actions</span>
              <span class="font-medium">{{ workflow.actions.length }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Status</span>
              <Badge :variant="workflow.is_active ? 'default' : 'secondary'">
                {{ workflow.is_active ? 'Active' : 'Inactive' }}
              </Badge>
            </div>
            <div
              v-if="workflow.last_run"
              class="flex items-center justify-between text-sm text-gray-500"
            >
              <div class="flex items-center space-x-1">
                <Clock class="h-4 w-4" />
                <span>Last run {{ formatDate(workflow.last_run) }}</span>
              </div>
            </div>
          </div>
          <div class="mt-4 flex space-x-2">
            <Button
              variant="outline"
              size="sm"
              class="flex-1"
              @click="router.push(`/projects/${projectId}/workflows/${workflow.id}/builder`)"
            >
              <Edit class="w-4 h-4 mr-1" />
              Edit
            </Button>
            <Button variant="outline" size="sm" class="flex-1" @click="handleRunWorkflow(workflow)">
              <Play class="w-4 h-4 mr-1" />
              Run
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
