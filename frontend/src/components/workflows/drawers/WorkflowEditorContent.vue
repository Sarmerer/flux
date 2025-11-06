<script setup lang="ts">
import { Globe, Plus, Trash2 } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

import type { WorkflowAction, WorkflowTrigger } from '@/types/api'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { SheetDescription, SheetFooter, SheetHeader, SheetTitle } from '@/components/ui/sheet'
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

import { WORKFLOW_ACTION_TYPES, WORKFLOW_TRIGGER_TYPES } from '@/constants/workflows'

interface Table {
  id: string
  name: string
}

interface Props {
  tables: Table[]
  loadingTables?: boolean
}

interface Emits {
  (e: 'save', data: CreateWorkflowData): void
  (e: 'close'): void
}

export interface CreateWorkflowData {
  name: string
  description: string
  trigger: WorkflowTrigger
  actions: WorkflowAction[]
  is_active: boolean
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const formData = ref<CreateWorkflowData>({
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

const validationErrors = ref<Record<string, string>>({})

const triggerRequiresTable = computed(() => {
  return formData.value.trigger.type.startsWith('on_row_')
})

const triggerRequiresSchedule = computed(() => {
  return formData.value.trigger.type === 'scheduled'
})

const triggerRequiresWebhook = computed(() => {
  return formData.value.trigger.type === 'webhook'
})

watch(
  () => formData.value.trigger.type,
  (newType) => {
    if (newType.startsWith('on_row_')) {
      formData.value.trigger = {
        type: newType as 'on_row_created' | 'on_row_updated' | 'on_row_deleted',
        table_id: '',
        conditions: {},
      }
    } else if (newType === 'scheduled') {
      formData.value.trigger = {
        type: 'scheduled',
        schedule: '',
      }
    } else if (newType === 'webhook') {
      formData.value.trigger = {
        type: 'webhook',
        webhook_url: '',
      }
    }
    validationErrors.value = {}
  }
)

const addAction = () => {
  const newAction: WorkflowAction = {
    id: crypto.randomUUID(),
    type: 'send_webhook',
    config: {
      url: '',
      method: 'POST',
    },
    order: formData.value.actions.length + 1,
  }
  formData.value.actions.push(newAction)
  validationErrors.value = {}
}

const removeAction = (actionId: string) => {
  formData.value.actions = formData.value.actions.filter((action) => action.id !== actionId)
  formData.value.actions.forEach((action, index) => {
    action.order = index + 1
  })
}

const updateActionType = (actionId: string, newType: string) => {
  const action = formData.value.actions.find((a) => a.id === actionId)
  if (!action) return

  action.type = newType as WorkflowAction['type']

  if (newType === 'send_webhook') {
    action.config = { url: '', method: 'POST' }
  } else if (newType === 'send_email') {
    action.config = { to: '', subject: '', template: '' }
  } else if (newType === 'update_row') {
    action.config = { table: '', condition: '', updates: {} }
  } else if (newType === 'create_row') {
    action.config = { table: '', values: {} }
  } else if (newType === 'delete_row') {
    action.config = { table: '', condition: '' }
  }
}

const validate = (): boolean => {
  const errors: Record<string, string> = {}

  if (!formData.value.name.trim()) {
    errors.name = 'Workflow name is required'
  }

  if (triggerRequiresTable.value && !('table_id' in formData.value.trigger && formData.value.trigger.table_id)) {
    errors.trigger_table = 'Table selection is required for row-based triggers'
  }

  if (triggerRequiresSchedule.value && !('schedule' in formData.value.trigger && formData.value.trigger.schedule)) {
    errors.trigger_schedule = 'Schedule is required for scheduled triggers'
  }

  if (triggerRequiresWebhook.value && !('webhook_url' in formData.value.trigger && formData.value.trigger.webhook_url)) {
    errors.trigger_webhook = 'Webhook URL is required for webhook triggers'
  }

  if (formData.value.actions.length === 0) {
    errors.actions = 'At least one action is required'
  }

  formData.value.actions.forEach((action, index) => {
    if (action.type === 'send_webhook' && !action.config.url) {
      errors[`action_${index}_url`] = 'URL is required'
    }
    if (action.type === 'send_email' && !action.config.to) {
      errors[`action_${index}_to`] = 'Recipient is required'
    }
    if ((action.type === 'update_row' || action.type === 'create_row' || action.type === 'delete_row') && !action.config.table) {
      errors[`action_${index}_table`] = 'Table name is required'
    }
  })

  validationErrors.value = errors
  return Object.keys(errors).length === 0
}

const handleSave = () => {
  if (!validate()) return
  emit('save', formData.value)
}

const handleClose = () => {
  emit('close')
}

const resetForm = () => {
  formData.value = {
    name: '',
    description: '',
    trigger: {
      type: 'on_row_created',
      table_id: '',
      conditions: {},
    },
    actions: [],
    is_active: false,
  }
  validationErrors.value = {}
}

const getActionIcon = (actionType: string) => {
  const action = WORKFLOW_ACTION_TYPES.find((a) => a.value === actionType)
  return action?.icon || Globe
}
</script>

<template>
  <div class="flex flex-col h-full">
    <SheetHeader class="space-y-1 pb-4">
      <SheetTitle class="text-base font-medium">Create New Workflow</SheetTitle>
      <SheetDescription class="text-xs">
        Create a complete workflow with trigger and actions. All fields are validated before creation.
      </SheetDescription>
    </SheetHeader>

    <div class="flex-1 overflow-y-auto space-y-6 pb-4">
        <Card>
          <CardContent class="pt-6 space-y-4">
            <div class="space-y-2">
              <Label for="workflow-name">
                Workflow Name <span class="text-red-500">*</span>
              </Label>
              <Input
                id="workflow-name"
                v-model="formData.name"
                placeholder="Enter workflow name"
                :class="validationErrors.name && 'border-red-500'"
              />
              <p v-if="validationErrors.name" class="text-sm text-red-500">
                {{ validationErrors.name }}
              </p>
            </div>

            <div class="space-y-2">
              <Label for="workflow-description">Description</Label>
              <Textarea
                id="workflow-description"
                v-model="formData.description"
                placeholder="Describe what this workflow does"
                rows="2"
              />
            </div>

            <div class="flex items-center space-x-2">
              <Checkbox v-model:checked="formData.is_active" id="is-active" />
              <Label for="is-active" class="font-normal">Activate workflow immediately</Label>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="pt-6 space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-semibold">Trigger Configuration</h3>
              <Badge variant="outline">Required</Badge>
            </div>

            <div class="space-y-2">
              <Label for="trigger-type">
                Trigger Type <span class="text-red-500">*</span>
              </Label>
              <Select v-model="formData.trigger.type">
                <SelectTrigger id="trigger-type">
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

            <div v-if="triggerRequiresTable" class="space-y-2">
              <Label for="table-id">
                Table <span class="text-red-500">*</span>
              </Label>
              <Select
                :model-value="'table_id' in formData.trigger ? formData.trigger.table_id : ''"
                @update:model-value="
                  (value) => {
                    if ('table_id' in formData.trigger && typeof value === 'string') {
                      formData.trigger.table_id = value
                    }
                  }
                "
                :disabled="loadingTables"
              >
                <SelectTrigger
                  id="table-id"
                  :class="validationErrors.trigger_table && 'border-red-500'"
                >
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
              <p v-if="validationErrors.trigger_table" class="text-sm text-red-500">
                {{ validationErrors.trigger_table }}
              </p>
            </div>

            <div v-if="triggerRequiresSchedule" class="space-y-2">
              <Label for="schedule">
                Schedule (Cron) <span class="text-red-500">*</span>
              </Label>
              <Input
                id="schedule"
                :model-value="'schedule' in formData.trigger ? formData.trigger.schedule : ''"
                @update:model-value="
                  (value) => {
                    if ('schedule' in formData.trigger && typeof value === 'string') {
                      formData.trigger.schedule = value
                    }
                  }
                "
                placeholder="0 2 * * * (daily at 2 AM)"
                :class="validationErrors.trigger_schedule && 'border-red-500'"
              />
              <p v-if="validationErrors.trigger_schedule" class="text-sm text-red-500">
                {{ validationErrors.trigger_schedule }}
              </p>
              <p class="text-xs text-muted-foreground">
                Format: minute hour day month weekday
              </p>
            </div>

            <div v-if="triggerRequiresWebhook" class="space-y-2">
              <Label for="webhook-url">
                Webhook URL <span class="text-red-500">*</span>
              </Label>
              <Input
                id="webhook-url"
                :model-value="'webhook_url' in formData.trigger ? formData.trigger.webhook_url : ''"
                @update:model-value="
                  (value) => {
                    if ('webhook_url' in formData.trigger && typeof value === 'string') {
                      formData.trigger.webhook_url = value
                    }
                  }
                "
                placeholder="https://api.example.com/webhook"
                :class="validationErrors.trigger_webhook && 'border-red-500'"
              />
              <p v-if="validationErrors.trigger_webhook" class="text-sm text-red-500">
                {{ validationErrors.trigger_webhook }}
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="pt-6 space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-sm font-semibold">Actions</h3>
                <p class="text-xs text-muted-foreground">
                  Define what happens when the workflow is triggered
                </p>
              </div>
              <Button variant="outline" size="sm" @click="addAction">
                <Plus class="w-4 h-4 mr-1" />
                Add Action
              </Button>
            </div>

            <p v-if="validationErrors.actions" class="text-sm text-red-500">
              {{ validationErrors.actions }}
            </p>

            <div v-if="formData.actions.length === 0" class="text-center py-6 text-sm text-muted-foreground border-2 border-dashed rounded-lg">
              No actions added yet. Click "Add Action" to get started.
            </div>

            <div class="space-y-3">
              <div
                v-for="(action, index) in formData.actions"
                :key="action.id"
                class="flex items-start space-x-3 p-4 border rounded-lg"
              >
                <div class="flex items-center justify-center w-8 h-8 bg-primary/10 rounded-full flex-shrink-0">
                  <span class="text-sm font-medium text-primary">{{ action.order }}</span>
                </div>

                <div class="flex-1 space-y-3">
                  <div class="flex items-center space-x-2">
                    <Select
                      :model-value="action.type"
                      @update:model-value="(val) => updateActionType(action.id, val as string)"
                    >
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
                    <component :is="getActionIcon(action.type)" class="w-4 h-4 text-muted-foreground" />
                  </div>

                  <div v-if="action.type === 'send_webhook'" class="space-y-2">
                    <div class="grid grid-cols-2 gap-2">
                      <div>
                        <Label class="text-xs">URL <span class="text-red-500">*</span></Label>
                        <Input
                          v-model="action.config.url"
                          placeholder="https://api.example.com/webhook"
                          :class="validationErrors[`action_${index}_url`] && 'border-red-500'"
                        />
                        <p v-if="validationErrors[`action_${index}_url`]" class="text-xs text-red-500">
                          {{ validationErrors[`action_${index}_url`] }}
                        </p>
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
                        rows="2"
                      />
                    </div>
                  </div>

                  <div v-else-if="action.type === 'send_email'" class="space-y-2">
                    <div class="grid grid-cols-2 gap-2">
                      <div>
                        <Label class="text-xs">To <span class="text-red-500">*</span></Label>
                        <Input
                          v-model="action.config.to"
                          placeholder="user@example.com"
                          :class="validationErrors[`action_${index}_to`] && 'border-red-500'"
                        />
                        <p v-if="validationErrors[`action_${index}_to`]" class="text-xs text-red-500">
                          {{ validationErrors[`action_${index}_to`] }}
                        </p>
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
                    v-else-if="action.type === 'update_row' || action.type === 'create_row' || action.type === 'delete_row'"
                    class="space-y-2"
                  >
                    <div class="grid grid-cols-2 gap-2">
                      <div>
                        <Label class="text-xs">Table <span class="text-red-500">*</span></Label>
                        <Input
                          v-model="action.config.table"
                          placeholder="table_name"
                          :class="validationErrors[`action_${index}_table`] && 'border-red-500'"
                        />
                        <p v-if="validationErrors[`action_${index}_table`]" class="text-xs text-red-500">
                          {{ validationErrors[`action_${index}_table`] }}
                        </p>
                      </div>
                      <div v-if="action.type !== 'create_row'">
                        <Label class="text-xs">Condition</Label>
                        <Input
                          :model-value="'condition' in action.config ? action.config.condition : ''"
                          @update:model-value="(val: string | number) => 'condition' in action.config && (action.config.condition = String(val))"
                          placeholder="id = {{row.id}}"
                        />
                      </div>
                    </div>
                    <div v-if="action.type === 'update_row'">
                      <Label class="text-xs">Updates (JSON)</Label>
                      <Textarea
                        :model-value="'updates' in action.config ? (typeof action.config.updates === 'string' ? action.config.updates : JSON.stringify(action.config.updates || {})) : ''"
                        @update:model-value="(val: string | number) => {
                          if ('updates' in action.config && typeof val === 'string') {
                            try {
                              action.config.updates = JSON.parse(val)
                            } catch {
                              action.config.updates = val as any
                            }
                          }
                        }"
                        placeholder='{"status": "processed"}'
                        rows="2"
                      />
                    </div>
                    <div v-if="action.type === 'create_row'">
                      <Label class="text-xs">Values (JSON)</Label>
                      <Textarea
                        :model-value="'values' in action.config ? (typeof action.config.values === 'string' ? action.config.values : JSON.stringify(action.config.values || {})) : ''"
                        @update:model-value="(val: string | number) => {
                          if ('values' in action.config && typeof val === 'string') {
                            try {
                              action.config.values = JSON.parse(val)
                            } catch {
                              action.config.values = val as any
                            }
                          }
                        }"
                        placeholder='{"name": "New Record"}'
                        rows="2"
                      />
                    </div>
                  </div>
                </div>

                <Button
                  variant="ghost"
                  size="sm"
                  @click="removeAction(action.id)"
                  class="text-red-600 hover:text-red-700 flex-shrink-0"
                >
                  <Trash2 class="w-4 h-4" />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
    </div>

    <SheetFooter class="flex-row gap-2 pt-4 border-t">
      <Button variant="outline" size="sm" @click="handleClose" class="flex-1">Cancel</Button>
      <Button size="sm" @click="handleSave" class="flex-1">Create Workflow</Button>
    </SheetFooter>
  </div>
</template>
