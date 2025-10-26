<script setup lang="ts">
import {
  ArrowLeft,
  Database,
  Eye,
  Globe,
  Mail,
  Play,
  Plus,
  Save,
  Settings,
  Trash2,
} from 'lucide-vue-next'
import type { Component } from 'vue'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { workflowService } from '@/api/services/workflow'
import type { WorkflowAction, WorkflowTrigger } from '@/types/api'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
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
const workflowId = computed(() => route.params.workflowId as string | undefined)
const isEditMode = computed(() => !!workflowId.value)

interface WorkflowState {
  name: string
  description: string
  trigger: WorkflowTrigger
  actions: WorkflowAction[]
  is_active: boolean
}

const workflow = ref<WorkflowState>({
  name: '',
  description: '',
  trigger: {
    type: 'on_row_created',
    table_name: '',
    conditions: {},
  },
  actions: [],
  is_active: false,
})

const isSaving = ref(false)
const isRunning = ref(false)
const showPreview = ref(true)

interface TypeOption {
  value: string
  label: string
  icon: Component
}

const triggerTypes: TypeOption[] = [
  { value: 'on_row_created', label: 'On Row Created', icon: Database },
  { value: 'on_row_updated', label: 'On Row Updated', icon: Database },
  { value: 'on_row_deleted', label: 'On Row Deleted', icon: Database },
  { value: 'scheduled', label: 'Scheduled', icon: Settings },
  { value: 'webhook', label: 'Webhook', icon: Globe },
]

const actionTypes: TypeOption[] = [
  { value: 'send_webhook', label: 'Send Webhook', icon: Globe },
  { value: 'send_email', label: 'Send Email', icon: Mail },
  { value: 'update_row', label: 'Update Row', icon: Database },
  { value: 'create_row', label: 'Create Row', icon: Database },
  { value: 'delete_row', label: 'Delete Row', icon: Database },
]

const addAction = () => {
  const newAction: WorkflowAction = {
    id: Date.now().toString(),
    type: 'send_webhook',
    config: {
      url: '',
      method: 'POST',
    },
    order: workflow.value.actions.length + 1,
  }
  workflow.value.actions.push(newAction)
}

const removeAction = (actionId: string) => {
  workflow.value.actions = workflow.value.actions.filter((action) => action.id !== actionId)

  workflow.value.actions.forEach((action, index) => {
    action.order = index + 1
  })
}

const moveAction = (fromIndex: number, toIndex: number) => {
  const action = workflow.value.actions.splice(fromIndex, 1)[0]
  if (action) {
    workflow.value.actions.splice(toIndex, 0, action)

    workflow.value.actions.forEach((act, index) => {
      act.order = index + 1
    })
  }
}

const loadWorkflowForEdit = async () => {
  if (!isEditMode.value || !workflowId.value) return

  try {
    const existingWorkflow = await workflowService.getById(projectId.value, workflowId.value)
    workflow.value = {
      name: existingWorkflow.name,
      description: existingWorkflow.description || '',
      trigger: existingWorkflow.trigger,
      actions: existingWorkflow.actions,
      is_active: existingWorkflow.is_active,
    }
  } catch (error) {
    console.error('Failed to load workflow:', error)
    alert('Failed to load workflow for editing')
    router.push(`/projects/${projectId.value}/workflows`)
  }
}

const saveWorkflow = async () => {
  if (!workflow.value.name.trim()) {
    alert('Please enter a workflow name')
    return
  }

  if (workflow.value.actions.length === 0) {
    alert('Please add at least one action')
    return
  }

  if (workflow.value.trigger.type.startsWith('on_row_') && !workflow.value.trigger.table_name) {
    alert('Please specify a table name for the trigger')
    return
  }

  isSaving.value = true
  try {
    if (isEditMode.value && workflowId.value) {
      await workflowService.update(projectId.value, workflowId.value, {
        name: workflow.value.name,
        description: workflow.value.description,
        trigger: workflow.value.trigger,
        actions: workflow.value.actions,
        is_active: workflow.value.is_active,
      })
    } else {
      await workflowService.create(projectId.value, {
        name: workflow.value.name,
        description: workflow.value.description,
        trigger: workflow.value.trigger,
        actions: workflow.value.actions,
        is_active: workflow.value.is_active,
      })
    }

    router.push(`/projects/${projectId.value}/workflows`)
  } catch (error) {
    console.error('Failed to save workflow:', error)
    alert('Failed to save workflow. Please check the console for details.')
  } finally {
    isSaving.value = false
  }
}

const runWorkflow = async () => {
  if (!workflowId.value) {
    alert('Please save the workflow before running it')
    return
  }

  isRunning.value = true
  try {
    const result = await workflowService.execute(projectId.value, workflowId.value)
    alert(`Workflow executed successfully!\n\n${result.message}`)
  } catch (error) {
    console.error('Failed to run workflow:', error)
    alert('Failed to run workflow. Please check the console for details.')
  } finally {
    isRunning.value = false
  }
}

onMounted(() => {
  loadWorkflowForEdit()
})

const getActionIcon = (actionType: string): Component => {
  const action = actionTypes.find((a) => a.value === actionType)
  return action?.icon || Globe
}

const getTriggerIcon = (triggerType: string): Component => {
  const trigger = triggerTypes.find((t) => t.value === triggerType)
  return trigger?.icon || Database
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <Button variant="ghost" size="sm" @click="router.push(`/projects/${projectId}/workflows`)">
          <ArrowLeft class="w-4 h-4 mr-2" />
          Back to Workflows
        </Button>
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Workflow Builder</h1>
          <p class="text-gray-600">Design and configure your automation workflow</p>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <Button variant="outline" @click="showPreview = !showPreview">
          <Eye v-if="!showPreview" class="w-4 h-4 mr-2" />
          <Eye v-else class="w-4 h-4 mr-2" />
          {{ showPreview ? 'Hide' : 'Show' }} Preview
        </Button>
        <Button variant="outline" @click="runWorkflow" :disabled="isRunning">
          <Play class="w-4 h-4 mr-2" />
          {{ isRunning ? 'Running...' : 'Test Run' }}
        </Button>
        <Button @click="saveWorkflow" :disabled="isSaving">
          <Save class="w-4 h-4 mr-2" />
          {{ isSaving ? 'Saving...' : 'Save Workflow' }}
        </Button>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Workflow Configuration -->
      <Card>
        <CardHeader>
          <CardTitle>Workflow Configuration</CardTitle>
          <CardDescription>Basic workflow settings and metadata</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="workflow-name">Workflow Name</Label>
            <Input id="workflow-name" v-model="workflow.name" placeholder="Enter workflow name" />
          </div>
          <div class="space-y-2">
            <Label for="workflow-description">Description</Label>
            <Textarea
              id="workflow-description"
              v-model="workflow.description"
              placeholder="Enter workflow description"
              rows="3"
            />
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox v-model="workflow.is_active" id="is-active" />
            <Label for="is-active">Activate workflow</Label>
          </div>
        </CardContent>
      </Card>

      <!-- Trigger Configuration -->
      <Card>
        <CardHeader>
          <CardTitle>Trigger Configuration</CardTitle>
          <CardDescription>Define what triggers this workflow</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="trigger-type">Trigger Type</Label>
            <Select v-model="workflow.trigger.type">
              <SelectTrigger>
                <SelectValue placeholder="Select trigger type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="trigger in triggerTypes"
                  :key="trigger.value"
                  :value="trigger.value"
                >
                  <div class="flex items-center space-x-2">
                    <component :is="trigger.icon" class="w-4 h-4" />
                    <span>{{ trigger.label }}</span>
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div v-if="workflow.trigger.type.startsWith('on_row_')" class="space-y-2">
            <Label for="table-name">Table Name</Label>
            <Input
              id="table-name"
              v-model="workflow.trigger.table_name"
              placeholder="Enter table name"
            />
          </div>

          <div v-if="workflow.trigger.type === 'scheduled'" class="space-y-2">
            <Label for="schedule">Schedule (Cron)</Label>
            <Input
              id="schedule"
              v-model="workflow.trigger.schedule"
              placeholder="0 2 * * * (daily at 2 AM)"
            />
          </div>

          <div v-if="workflow.trigger.type === 'webhook'" class="space-y-2">
            <Label for="webhook-url">Webhook URL</Label>
            <Input
              id="webhook-url"
              v-model="workflow.trigger.webhook_url"
              placeholder="https://api.example.com/webhook"
            />
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Actions Configuration -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>Actions</CardTitle>
            <CardDescription>Define what happens when the workflow is triggered</CardDescription>
          </div>
          <Button variant="outline" size="sm" @click="addAction">
            <Plus class="w-4 h-4 mr-1" />
            Add Action
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div class="space-y-4">
          <div
            v-for="(action, index) in workflow.actions"
            :key="action.id"
            class="flex items-start space-x-4 p-4 border rounded-lg"
          >
            <div class="flex items-center space-x-2">
              <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                <span class="text-sm font-medium text-blue-600">{{ action.order }}</span>
              </div>
              <component :is="getActionIcon(action.type)" class="w-5 h-5 text-blue-600" />
            </div>

            <div class="flex-1 space-y-3">
              <div class="flex items-center space-x-2">
                <Select v-model="action.type">
                  <SelectTrigger class="w-48">
                    <SelectValue placeholder="Select action type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="actionType in actionTypes"
                      :key="actionType.value"
                      :value="actionType.value"
                    >
                      <div class="flex items-center space-x-2">
                        <component :is="actionType.icon" class="w-4 h-4" />
                        <span>{{ actionType.label }}</span>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <Badge variant="outline">{{
                  actionTypes.find((a) => a.value === action.type)?.label
                }}</Badge>
              </div>

              <!-- Action-specific configuration -->
              <div v-if="action.type === 'send_webhook'" class="space-y-2">
                <div class="grid grid-cols-2 gap-2">
                  <div>
                    <Label class="text-xs">URL</Label>
                    <Input
                      v-model="action.config.url"
                      placeholder="https://api.example.com/webhook"
                    />
                  </div>
                  <div>
                    <Label class="text-xs">Method</Label>
                    <Select v-model="action.config.method">
                      <SelectTrigger>
                        <SelectValue placeholder="Method" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="POST">POST</SelectItem>
                        <SelectItem value="PUT">PUT</SelectItem>
                        <SelectItem value="PATCH">PATCH</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div>
                  <Label class="text-xs">Payload (JSON)</Label>
                  <Textarea
                    v-model="action.config.payload"
                    placeholder='{"message": "Hello from FlowDB!"}'
                    rows="3"
                  />
                </div>
              </div>

              <div v-else-if="action.type === 'send_email'" class="space-y-2">
                <div class="grid grid-cols-2 gap-2">
                  <div>
                    <Label class="text-xs">To</Label>
                    <Input v-model="action.config.to" placeholder="user@example.com" />
                  </div>
                  <div>
                    <Label class="text-xs">Subject</Label>
                    <Input v-model="action.config.subject" placeholder="Email subject" />
                  </div>
                </div>
                <div>
                  <Label class="text-xs">Template</Label>
                  <Input v-model="action.config.template" placeholder="welcome" />
                </div>
              </div>

              <div
                v-else-if="action.type.startsWith('update_') || action.type.startsWith('create_')"
                class="space-y-2"
              >
                <div class="grid grid-cols-2 gap-2">
                  <div>
                    <Label class="text-xs">Table</Label>
                    <Input v-model="action.config.table" placeholder="table_name" />
                  </div>
                  <div>
                    <Label class="text-xs">Condition</Label>
                    <Input v-model="action.config.condition" placeholder="id = {{row.id}}" />
                  </div>
                </div>
                <div>
                  <Label class="text-xs">Updates (JSON)</Label>
                  <Textarea
                    v-model="action.config.updates"
                    placeholder='{"status": "processed", "updated_at": "{{now}}"}'
                    rows="2"
                  />
                </div>
              </div>
            </div>

            <Button
              variant="ghost"
              size="sm"
              @click="removeAction(action.id)"
              class="text-red-600 hover:text-red-700"
            >
              <Trash2 class="w-4 h-4" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Workflow Preview -->
    <Card v-if="showPreview">
      <CardHeader>
        <CardTitle>Workflow Preview</CardTitle>
        <CardDescription>Visual representation of your workflow</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="space-y-4">
          <!-- Trigger -->
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2 p-3 bg-blue-100 rounded-lg">
              <component
                :is="getTriggerIcon(workflow.trigger.type)"
                class="w-5 h-5 text-blue-600"
              />
              <span class="font-medium">{{
                triggerTypes.find((t) => t.value === workflow.trigger.type)?.label
              }}</span>
            </div>
            <div class="text-gray-400">→</div>
          </div>

          <!-- Actions -->
          <div class="space-y-2">
            <div
              v-for="(action, index) in workflow.actions"
              :key="action.id"
              class="flex items-center space-x-4"
            >
              <div v-if="index > 0" class="text-gray-400 ml-6">↓</div>
              <div class="flex items-center space-x-2 p-3 bg-green-100 rounded-lg">
                <component :is="getActionIcon(action.type)" class="w-5 h-5 text-green-600" />
                <span class="font-medium">{{
                  actionTypes.find((a) => a.value === action.type)?.label
                }}</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
