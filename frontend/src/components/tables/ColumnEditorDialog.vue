<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Info, Sparkles } from 'lucide-vue-next'

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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

interface ColumnData {
  name: string
  type: string
  nullable: boolean
  default_value: string
  is_identity: boolean
  unique: boolean
  is_primary_key: boolean
}

interface Props {
  open: boolean
  column?: ColumnData
  mode?: 'create' | 'edit'
  existingColumns?: string[]
  availableTables?: string[]
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'save', column: ColumnData): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'create',
  existingColumns: () => [],
  availableTables: () => [],
})

const emit = defineEmits<Emits>()

const formData = ref<ColumnData>({
  name: '',
  type: 'TEXT',
  nullable: true,
  default_value: '',
  is_identity: false,
  unique: false,
  is_primary_key: false,
})

const dataTypeCategories = {
  numeric: [
    { value: 'SMALLINT', label: 'Small Integer', description: '2 bytes, -32K to 32K' },
    { value: 'INTEGER', label: 'Integer', description: '4 bytes, -2B to 2B' },
    { value: 'BIGINT', label: 'Big Integer', description: '8 bytes, very large numbers' },
    { value: 'DECIMAL', label: 'Decimal', description: 'Exact numeric with precision' },
    { value: 'NUMERIC', label: 'Numeric', description: 'Same as DECIMAL' },
    { value: 'REAL', label: 'Real', description: '4 bytes, 6 decimal precision' },
    { value: 'DOUBLE PRECISION', label: 'Double', description: '8 bytes, 15 decimal precision' },
  ],
  text: [
    { value: 'TEXT', label: 'Text', description: 'Unlimited length text' },
    { value: 'VARCHAR', label: 'Varchar', description: 'Variable length with limit' },
    { value: 'CHAR', label: 'Char', description: 'Fixed length text' },
  ],
  datetime: [
    { value: 'TIMESTAMP', label: 'Timestamp', description: 'Date and time' },
    { value: 'TIMESTAMPTZ', label: 'Timestamp TZ', description: 'With timezone' },
    { value: 'DATE', label: 'Date', description: 'Date only (no time)' },
    { value: 'TIME', label: 'Time', description: 'Time only (no date)' },
  ],
  special: [
    { value: 'UUID', label: 'UUID', description: 'Universally unique identifier' },
    { value: 'BOOLEAN', label: 'Boolean', description: 'True or false' },
    { value: 'JSON', label: 'JSON', description: 'JSON data' },
    { value: 'JSONB', label: 'JSONB', description: 'Binary JSON (faster)' },
  ],
}

const defaultValueSuggestions = computed(() => {
  const type = formData.value.type.toUpperCase()

  if (type.includes('TIMESTAMP')) {
    return [
      { value: 'NOW()', label: 'NOW()', description: 'Current timestamp' },
      { value: 'CURRENT_TIMESTAMP', label: 'CURRENT_TIMESTAMP', description: 'Current timestamp' },
    ]
  }

  if (type === 'DATE') {
    return [{ value: 'CURRENT_DATE', label: 'CURRENT_DATE', description: 'Current date' }]
  }

  if (type === 'UUID') {
    return [
      { value: 'gen_random_uuid()', label: 'gen_random_uuid()', description: 'Random UUID' },
    ]
  }

  if (type === 'BOOLEAN') {
    return [
      { value: 'TRUE', label: 'TRUE', description: 'Default true' },
      { value: 'FALSE', label: 'FALSE', description: 'Default false' },
    ]
  }

  if (['SMALLINT', 'INTEGER', 'BIGINT', 'DECIMAL', 'NUMERIC'].includes(type)) {
    return [
      { value: '0', label: '0', description: 'Zero' },
      { value: '1', label: '1', description: 'One' },
    ]
  }

  return []
})

const isNumericType = computed(() => {
  const type = formData.value.type.toUpperCase()
  return ['SMALLINT', 'INTEGER', 'BIGINT'].includes(type)
})

const canHaveDefault = computed(() => {
  return !formData.value.is_identity
})

const validationErrors = computed(() => {
  const errors: string[] = []

  if (!formData.value.name.trim()) {
    errors.push('Column name is required')
  } else if (!/^[a-zA-Z_][a-zA-Z0-9_]*$/.test(formData.value.name)) {
    errors.push('Column name must start with letter/underscore and contain only alphanumeric/underscore')
  } else if (
    props.mode === 'create' &&
    props.existingColumns.includes(formData.value.name.toLowerCase())
  ) {
    errors.push('Column name already exists')
  }

  if (formData.value.is_identity && !isNumericType.value) {
    errors.push('Auto-increment requires SMALLINT, INTEGER, or BIGINT type')
  }

  if (formData.value.is_identity && formData.value.default_value) {
    errors.push('Auto-increment columns cannot have default values')
  }

  return errors
})

const isValid = computed(() => validationErrors.value.length === 0)

const dialogTitle = computed(() => {
  return props.mode === 'create' ? 'Add Column' : 'Edit Column'
})

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      if (props.column) {
        formData.value = { ...props.column }
      } else {
        formData.value = {
          name: '',
          type: 'TEXT',
          nullable: true,
          default_value: '',
          is_identity: false,
          unique: false,
          is_primary_key: false,
        }
      }
    }
  }
)

watch(
  () => formData.value.is_identity,
  (isIdentity) => {
    if (isIdentity) {
      formData.value.default_value = ''
      formData.value.nullable = false
    }
  }
)

const handleSave = () => {
  if (!isValid.value) return
  emit('save', { ...formData.value })
  emit('update:open', false)
}

const handleCancel = () => {
  emit('update:open', false)
}

const setDefaultValue = (value: string) => {
  formData.value.default_value = value
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="max-w-2xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>{{ dialogTitle }}</DialogTitle>
        <DialogDescription>
          Configure the column properties, data type, and constraints.
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-6">
        <div class="space-y-2">
          <Label for="column-name">Column Name *</Label>
          <Input
            id="column-name"
            v-model="formData.name"
            placeholder="e.g., user_id, email, created_at"
            :disabled="mode === 'edit'"
          />
          <p class="text-xs text-muted-foreground">
            Must start with a letter or underscore, alphanumeric characters only
          </p>
        </div>

        <div class="space-y-2">
          <Label>Data Type *</Label>
          <div class="grid grid-cols-2 gap-3">
            <div v-for="(types, category) in dataTypeCategories" :key="category" class="space-y-2">
              <div class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                {{ category }}
              </div>
              <Select v-model="formData.type">
                <SelectTrigger>
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="type in types" :key="type.value" :value="type.value">
                    <div class="flex flex-col">
                      <span>{{ type.label }}</span>
                      <span class="text-xs text-muted-foreground">{{ type.description }}</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>

        <div class="space-y-3">
          <Label>Constraints</Label>
          <div class="space-y-3">
            <div class="flex items-center space-x-2">
              <Checkbox id="nullable" v-model:checked="formData.nullable" :disabled="formData.is_identity" />
              <Label for="nullable" class="font-normal cursor-pointer">Allow NULL values</Label>
            </div>

            <div class="flex items-center space-x-2">
              <Checkbox id="unique" v-model:checked="formData.unique" />
              <Label for="unique" class="font-normal cursor-pointer">Unique constraint</Label>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <Info class="w-3.5 h-3.5 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent>
                    <p class="max-w-xs text-xs">All values in this column must be unique</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>

            <div class="flex items-center space-x-2">
              <Checkbox id="primary-key" v-model:checked="formData.is_primary_key" />
              <Label for="primary-key" class="font-normal cursor-pointer">Primary Key</Label>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <Info class="w-3.5 h-3.5 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent>
                    <p class="max-w-xs text-xs">Uniquely identifies each row in the table</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>

            <div v-if="isNumericType" class="flex items-center space-x-2">
              <Checkbox id="is-identity" v-model:checked="formData.is_identity" />
              <Label for="is-identity" class="font-normal cursor-pointer flex items-center gap-1">
                Auto-increment
                <Sparkles class="w-3.5 h-3.5 text-primary" />
              </Label>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <Info class="w-3.5 h-3.5 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent>
                    <p class="max-w-xs text-xs">
                      Automatically generates sequential numbers (1, 2, 3...) using PostgreSQL IDENTITY
                    </p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </div>
        </div>

        <div v-if="canHaveDefault" class="space-y-2">
          <Label for="default-value">Default Value</Label>
          <Input
            id="default-value"
            v-model="formData.default_value"
            placeholder="e.g., NOW(), 'active', 0"
          />

          <div v-if="defaultValueSuggestions.length > 0" class="flex flex-wrap gap-2 pt-1">
            <Button
              v-for="suggestion in defaultValueSuggestions"
              :key="suggestion.value"
              type="button"
              variant="outline"
              size="sm"
              @click="setDefaultValue(suggestion.value)"
              class="h-7 text-xs"
            >
              {{ suggestion.label }}
            </Button>
          </div>

          <p class="text-xs text-muted-foreground">
            PostgreSQL expressions: NOW(), gen_random_uuid(), CURRENT_DATE, literals ('text', 0, TRUE)
          </p>
        </div>

        <div v-if="validationErrors.length > 0" class="rounded-lg bg-destructive/10 p-3 space-y-1">
          <p class="text-sm font-medium text-destructive">Validation Errors:</p>
          <ul class="list-disc list-inside text-xs text-destructive space-y-0.5">
            <li v-for="error in validationErrors" :key="error">{{ error }}</li>
          </ul>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="handleCancel">Cancel</Button>
        <Button @click="handleSave" :disabled="!isValid">
          {{ mode === 'create' ? 'Add Column' : 'Save Changes' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
