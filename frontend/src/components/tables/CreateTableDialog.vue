<script setup lang="ts">
import { Plus, X } from 'lucide-vue-next'
import { computed, ref } from 'vue'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
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

interface Column {
  name: string
  type: string
  nullable: boolean
  primary_key: boolean
  default_value: string
  is_identity: boolean
  unique: boolean
}

interface Props {
  open: boolean
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'create', data: CreateTableData): void
}

export interface CreateTableData {
  name: string
  description: string
  columns: Column[]
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const formData = ref<CreateTableData>({
  name: '',
  description: '',
  columns: [],
})

const validationErrors = ref<Record<string, string>>({})

const columnTypes = [
  'TEXT',
  'VARCHAR',
  'CHAR',
  'INTEGER',
  'BIGINT',
  'SMALLINT',
  'DECIMAL',
  'REAL',
  'DOUBLE PRECISION',
  'BOOLEAN',
  'DATE',
  'TIMESTAMP',
  'TIMESTAMPTZ',
  'UUID',
  'JSON',
  'JSONB',
]

const hasValidationErrors = computed(() => Object.keys(validationErrors.value).length > 0)

const addColumn = () => {
  formData.value.columns.push({
    name: '',
    type: 'TEXT',
    nullable: true,
    primary_key: false,
    default_value: '',
    is_identity: false,
    unique: false,
  })
  validationErrors.value = {}
}

const removeColumn = (index: number) => {
  formData.value.columns.splice(index, 1)
  validationErrors.value = {}
}

const validate = (): boolean => {
  const errors: Record<string, string> = {}

  if (!formData.value.name.trim()) {
    errors.name = 'Table name is required'
  } else if (!/^[a-z_][a-z0-9_]*$/.test(formData.value.name)) {
    errors.name = 'Table name must start with a letter or underscore and contain only lowercase letters, numbers, and underscores'
  }

  if (formData.value.columns.length === 0) {
    errors.columns = 'At least one column is required'
  }

  formData.value.columns.forEach((col, index) => {
    if (!col.name.trim()) {
      errors[`column_${index}_name`] = 'Column name is required'
    } else if (!/^[a-z_][a-z0-9_]*$/.test(col.name)) {
      errors[`column_${index}_name`] = 'Invalid column name format'
    }

    if (!col.type) {
      errors[`column_${index}_type`] = 'Column type is required'
    }

    const duplicates = formData.value.columns.filter((c) => c.name.toLowerCase() === col.name.toLowerCase())
    if (duplicates.length > 1) {
      errors[`column_${index}_name`] = 'Duplicate column name'
    }
  })

  validationErrors.value = errors
  return Object.keys(errors).length === 0
}

const handleCreate = () => {
  if (!validate()) return
  emit('create', formData.value)
  resetForm()
}

const handleCancel = () => {
  resetForm()
  emit('update:open', false)
}

const resetForm = () => {
  formData.value = {
    name: '',
    description: '',
    columns: [],
  }
  validationErrors.value = {}
}

const addDefaultIdColumn = () => {
  const hasIdColumn = formData.value.columns.some((col) => col.name === 'id')
  if (hasIdColumn) return

  formData.value.columns.unshift({
    name: 'id',
    type: 'UUID',
    nullable: false,
    primary_key: true,
    default_value: 'gen_random_uuid()',
    is_identity: false,
    unique: true,
  })
}
</script>

<template>
  <Dialog :open="open" @update:open="(val) => !val && handleCancel()">
    <DialogContent class="max-w-3xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Create New Table</DialogTitle>
        <DialogDescription>
          Create a new table with custom columns. Table and column names must be lowercase with underscores.
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-6 py-4">
        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="table-name">
                Table Name <span class="text-red-500">*</span>
              </Label>
              <Input
                id="table-name"
                v-model="formData.name"
                placeholder="e.g., users, products, orders"
                :class="validationErrors.name && 'border-red-500'"
              />
              <p v-if="validationErrors.name" class="text-xs text-red-500">
                {{ validationErrors.name }}
              </p>
              <p v-else class="text-xs text-muted-foreground">
                Use lowercase letters, numbers, and underscores only
              </p>
            </div>
            <div class="space-y-2">
              <Label for="table-description">Description</Label>
              <Input
                id="table-description"
                v-model="formData.description"
                placeholder="Brief description of this table"
              />
            </div>
          </div>
        </div>

        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <Label class="text-sm font-semibold">Columns</Label>
              <p class="text-xs text-muted-foreground mt-1">
                Define the structure of your table
              </p>
            </div>
            <div class="flex gap-2">
              <Button type="button" variant="outline" size="sm" @click="addDefaultIdColumn">
                Add ID Column
              </Button>
              <Button type="button" variant="outline" size="sm" @click="addColumn">
                <Plus class="w-4 h-4 mr-1" />
                Add Column
              </Button>
            </div>
          </div>

          <p v-if="validationErrors.columns" class="text-sm text-red-500">
            {{ validationErrors.columns }}
          </p>

          <div v-if="formData.columns.length === 0" class="text-center py-8 text-sm text-muted-foreground border-2 border-dashed rounded-lg">
            No columns added yet. Click "Add Column" or "Add ID Column" to get started.
          </div>

          <div class="space-y-3">
            <div
              v-for="(column, index) in formData.columns"
              :key="index"
              class="p-4 border rounded-lg space-y-3"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="flex-1 grid grid-cols-2 gap-3">
                  <div class="space-y-1.5">
                    <Label class="text-xs">
                      Column Name <span class="text-red-500">*</span>
                    </Label>
                    <Input
                      v-model="column.name"
                      placeholder="e.g., email, created_at"
                      :class="validationErrors[`column_${index}_name`] && 'border-red-500'"
                    />
                    <p v-if="validationErrors[`column_${index}_name`]" class="text-xs text-red-500">
                      {{ validationErrors[`column_${index}_name`] }}
                    </p>
                  </div>

                  <div class="space-y-1.5">
                    <Label class="text-xs">
                      Data Type <span class="text-red-500">*</span>
                    </Label>
                    <Select v-model="column.type">
                      <SelectTrigger :class="validationErrors[`column_${index}_type`] && 'border-red-500'">
                        <SelectValue placeholder="Select type" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="type in columnTypes" :key="type" :value="type">
                          {{ type }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <p v-if="validationErrors[`column_${index}_type`]" class="text-xs text-red-500">
                      {{ validationErrors[`column_${index}_type`] }}
                    </p>
                  </div>
                </div>

                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  @click="removeColumn(index)"
                  class="text-red-600 hover:text-red-700 flex-shrink-0"
                >
                  <X class="w-4 h-4" />
                </Button>
              </div>

              <div class="space-y-1.5">
                <Label class="text-xs">Default Value</Label>
                <Input
                  v-model="column.default_value"
                  placeholder="e.g., '', 0, gen_random_uuid()"
                />
                <p class="text-xs text-muted-foreground">
                  Optional. Use SQL expressions like gen_random_uuid() or CURRENT_TIMESTAMP
                </p>
              </div>

              <div class="flex flex-wrap gap-4">
                <div class="flex items-center space-x-2">
                  <Checkbox v-model:checked="column.nullable" :id="`nullable-${index}`" />
                  <Label :for="`nullable-${index}`" class="text-xs font-normal cursor-pointer">
                    Nullable
                  </Label>
                </div>
                <div class="flex items-center space-x-2">
                  <Checkbox v-model:checked="column.primary_key" :id="`primary-${index}`" />
                  <Label :for="`primary-${index}`" class="text-xs font-normal cursor-pointer">
                    Primary Key
                  </Label>
                </div>
                <div class="flex items-center space-x-2">
                  <Checkbox v-model:checked="column.unique" :id="`unique-${index}`" />
                  <Label :for="`unique-${index}`" class="text-xs font-normal cursor-pointer">
                    Unique
                  </Label>
                </div>
                <div class="flex items-center space-x-2">
                  <Checkbox v-model:checked="column.is_identity" :id="`identity-${index}`" />
                  <Label :for="`identity-${index}`" class="text-xs font-normal cursor-pointer">
                    Identity
                  </Label>
                </div>
              </div>

              <div v-if="column.primary_key || column.unique || column.is_identity" class="flex gap-2">
                <Badge v-if="column.primary_key" variant="default" class="text-xs">Primary Key</Badge>
                <Badge v-if="column.unique" variant="secondary" class="text-xs">Unique</Badge>
                <Badge v-if="column.is_identity" variant="outline" class="text-xs">Identity</Badge>
              </div>
            </div>
          </div>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="handleCancel">Cancel</Button>
        <Button @click="handleCreate" :disabled="hasValidationErrors && formData.columns.length > 0">
          Create Table
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
