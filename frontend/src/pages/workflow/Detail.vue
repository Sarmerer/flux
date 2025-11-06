<script setup lang="ts">
import { ArrowLeft, CheckCircle, Database, Eye, Globe, Plus, Save, Trash2 } from 'lucide-vue-next'
import type { Component } from 'vue'
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'

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

import { useRouteContext } from '@/composables/routing'
import { useWorkflows } from '@/composables/api/useWorkflows'
import { useTables } from '@/composables/api/useTables'
import { useToast } from '@/composables/ui'
import { useWorkflowValidation } from '@/composables/useWorkflowValidation'
import { WORKFLOW_TRIGGER_TYPES, WORKFLOW_ACTION_TYPES } from '@/constants/workflows'

const router = useRouter()
const toast = useToast()

const { projectId, workflowId } = useRouteContext()
const { getWorkflowById, updateWorkflow } = useWorkflows(projectId.value)
const { tables, loading: loadingTables } = useTables(projectId.value)
const { validateWorkflow } = useWorkflowValidation()

const isSaving = ref(false)
const isLoading = ref(true)
const showPreview = ref(true)

interface WorkflowFormData {
  name: string
  description: string
  trigger: WorkflowTrigger
  actions: WorkflowAction[]
  is_active: boolean
}

const currentWorkflow = ref<WorkflowFormData>({
  name: '',
  description: '',
  trigger: {
    type: 'on_row_created',
    table_id: '',
    conditions: {},
  },
  actions: [],
  is_active: false,
})

const loadWorkflow = async () => {
  if (!workflowId.value) return

  try {
    isLoading.value = true
    const workflow = await getWorkflowById(workflowId.value)
    currentWorkflow.value = {
      name: workflow.name,
      description: workflow.description || '',
      trigger: workflow.trigger,
      actions: workflow.actions,
      is_active: workflow.is_active,
    }
  } catch (error) {
    console.error('Failed to load workflow:', error)
    toast.error('Load Failed', 'Failed to load workflow for editing')
    router.push(`/projects/${projectId.value}/workflows`)
  } finally {
    isLoading.value = false
  }
}

watch(workflowId, () => {
  if (workflowId.value) {
    loadWorkflow()
  }
}, { immediate: true })

const addAction = () => {
  const newAction: WorkflowAction = {
    id: crypto.randomUUID(),
    type: 'send_webhook',
    config: {
      url: '',
      method: 'POST',
    },
    order: currentWorkflow.value.actions.length + 1,
  }
  currentWorkflow.value.actions.push(newAction)
}

const removeAction = (actionId: string) => {
  currentWorkflow.value.actions = currentWorkflow.value.actions.filter(
    (action) => action.id !== actionId
  )

  currentWorkflow.value.actions.forEach((action, index) => {
    action.order = index + 1
  })
}

const saveWorkflow = async () => {
  const validation = validateWorkflow(
    currentWorkflow.value.name,
    currentWorkflow.value.trigger,
    currentWorkflow.value.actions
  )

  if (!validation.valid) {
    toast.error('Validation Error', validation.message || 'Invalid workflow data')
    return
  }

  isSaving.value = true
  try {
    if (workflowId.value) {
      await updateWorkflow(workflowId.value, {
        name: currentWorkflow.value.name,
        description: currentWorkflow.value.description,
        trigger: currentWorkflow.value.trigger,
        actions: currentWorkflow.value.actions,
        is_active: currentWorkflow.value.is_active,
      })
      toast.success('Workflow Updated', 'Your workflow has been updated successfully')
    }
  } catch (error) {
    console.error('Failed to save workflow:', error)
    toast.error('Save Failed', 'Failed to save workflow. Please try again.')
  } finally {
    isSaving.value = false
  }
}

const handleBackToList = () => {
  router.push(`/projects/${projectId.value}/workflows`)
}

const getActionIcon = (actionType: string): Component => {
  const action = WORKFLOW_ACTION_TYPES.find((a) => a.value === actionType)
  return action?.icon || Globe
}

const getTriggerIcon = (triggerType: string): Component => {
  const trigger = WORKFLOW_TRIGGER_TYPES.find((t) => t.value === triggerType)
  return trigger?.icon || Database
}
</script>

<template>
  <div v-if="isLoading" class="flex items-center justify-center h-[calc(100vh-3.5rem)]">
    <div class="text-center">
      <p class="text-muted-foreground">Loading workflow...</p>
    </div>
  </div>

  <div v-else class="p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <Button variant="ghost" size="sm" @click="handleBackToList">
          <ArrowLeft class="w-4 h-4 mr-2" />
          Back to Workflows
        </Button>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-3xl font-bold">{{ currentWorkflow.name }}</h1>
            <Badge v-if="currentWorkflow.is_active" variant="default">
              <CheckCircle class="w-3 h-3 mr-1" />
              Active
            </Badge>
            <Badge v-else variant="secondary">Inactive</Badge>
          </div>
          <p class="text-muted-foreground mt-1">
            {{ currentWorkflow.description || 'No description provided' }}
          </p>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <Button variant="outline" @click="showPreview = !showPreview">
          <Eye class="w-4 h-4 mr-2" />
          {{ showPreview ? 'Hide' : 'Show' }} Preview
        </Button>
        <Button @click="saveWorkflow" :disabled="isSaving">
          <Save class="w-4 h-4 mr-2" />
          {{ isSaving ? 'Saving...' : 'Save Changes' }}
        </Button>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Workflow Configuration</CardTitle>
          <CardDescription>Basic workflow settings and metadata</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="workflow-name">Workflow Name</Label>
            <Input
              id="workflow-name"
              v-model="currentWorkflow.name"
              placeholder="Enter workflow name"
            />
          </div>
          <div class="space-y-2">
            <Label for="workflow-description">Description</Label>
            <Textarea
              id="workflow-description"
              v-model="currentWorkflow.description"
              placeholder="Enter workflow description"
              rows="3"
            />
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox v-model:checked="currentWorkflow.is_active" id="is-active" />
            <Label for="is-active">Activate workflow</Label>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Trigger Configuration</CardTitle>
          <CardDescription>Define what triggers this workflow</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="trigger-type">Trigger Type</Label>
            <Select v-model="currentWorkflow.trigger.type">
              <SelectTrigger>
                <SelectValue placeholder="Select trigger type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="trigger in WORKFLOW_TRIGGER_TYPES"
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

          <div
            v-if="currentWorkflow.trigger.type.startsWith('on_row_')"
            class="space-y-2"
          >
            <Label for="table-id">Table</Label>
            <Select
              :model-value="
                'table_id' in currentWorkflow.trigger ? currentWorkflow.trigger.table_id : ''
              "
              @update:model-value="
                (value) => {
                  if ('table_id' in currentWorkflow.trigger && typeof value === 'string') {
                    currentWorkflow.trigger.table_id = value
                  }
                }
              "
              :disabled="loadingTables"
            >
              <SelectTrigger id="table-id">
                <SelectValue
                  :placeholder="loadingTables ? 'Loading tables...' : 'Select a table'"
                />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="table in tables" :key="table.id" :value="table.id">
                  {{ table.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div v-if="currentWorkflow.trigger.type === 'scheduled'" class="space-y-2">
            <Label for="schedule">Schedule (Cron)</Label>
            <Input
              id="schedule"
              :model-value="'schedule' in currentWorkflow.trigger ? currentWorkflow.trigger.schedule : ''"
              @update:model-value="(value) => { if ('schedule' in currentWorkflow.trigger && typeof value === 'string') { currentWorkflow.trigger.schedule = value } }"
              placeholder="0 2 * * * (daily at 2 AM)"
            />
          </div>

          <div v-if="currentWorkflow.trigger.type === 'webhook'" class="space-y-2">
            <Label for="webhook-url">Webhook URL</Label>
            <Input
              id="webhook-url"
              :model-value="'webhook_url' in currentWorkflow.trigger ? currentWorkflow.trigger.webhook_url : ''"
              @update:model-value="(value) => { if ('webhook_url' in currentWorkflow.trigger && typeof value === 'string') { currentWorkflow.trigger.webhook_url = value } }"
              placeholder="https://api.example.com/webhook"
            />
          </div>
        </CardContent>
      </Card>
    </div>

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
            v-for="action in currentWorkflow.actions"
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
                      v-for="actionType in WORKFLOW_ACTION_TYPES"
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
                  WORKFLOW_ACTION_TYPES.find((a) => a.value === action.type)?.label
                }}</Badge>
              </div>

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
                    :model-value="typeof action.config.payload === 'string' ? action.config.payload : JSON.stringify(action.config.payload || {})"
                    @update:model-value="(val: string | number) => {
                      if (typeof val === 'string') {
                        try {
                          action.config.payload = JSON.parse(val)
                        } catch {
                          action.config.payload = val as any
                        }
                      }
                    }"
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
                  <div v-if="'condition' in action.config">
                    <Label class="text-xs">Condition</Label>
                    <Input :model-value="action.config.condition" @update:model-value="(val: string | number) => 'condition' in action.config && (action.config.condition = String(val))" placeholder="id = {{row.id}}" />
                  </div>
                </div>
                <div v-if="'updates' in action.config">
                  <Label class="text-xs">Updates (JSON)</Label>
                  <Textarea
                    :model-value="typeof action.config.updates === 'string' ? action.config.updates : JSON.stringify(action.config.updates || {})"
                    @update:model-value="(val: string | number) => {
                      if ('updates' in action.config && typeof val === 'string') {
                        try {
                          action.config.updates = JSON.parse(val)
                        } catch {
                          action.config.updates = val as any
                        }
                      }
                    }"
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

    <Card v-if="showPreview">
      <CardHeader>
        <CardTitle>Workflow Preview</CardTitle>
        <CardDescription>Visual representation of your workflow</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="space-y-4">
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2 p-3 bg-blue-100 rounded-lg">
              <component
                :is="getTriggerIcon(currentWorkflow.trigger.type)"
                class="w-5 h-5 text-blue-600"
              />
              <span class="font-medium">{{
                WORKFLOW_TRIGGER_TYPES.find((t) => t.value === currentWorkflow.trigger.type)
                  ?.label
              }}</span>
            </div>
            <div class="text-gray-400">→</div>
          </div>

          <div class="space-y-2">
            <div
              v-for="(action, index) in currentWorkflow.actions"
              :key="action.id"
              class="flex items-center space-x-4"
            >
              <div v-if="index > 0" class="text-gray-400 ml-6">↓</div>
              <div class="flex items-center space-x-2 p-3 bg-green-100 rounded-lg">
                <component :is="getActionIcon(action.type)" class="w-5 h-5 text-green-600" />
                <span class="font-medium">{{
                  WORKFLOW_ACTION_TYPES.find((a) => a.value === action.type)?.label
                }}</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
