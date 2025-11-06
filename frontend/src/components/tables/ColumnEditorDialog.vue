<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import ColumnFormField from '@/components/tables/ColumnFormField.vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

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

  const isNumericType = ['SMALLINT', 'INTEGER', 'BIGINT'].includes(formData.value.type.toUpperCase())

  if (formData.value.is_identity && !isNumericType) {
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

const handleSave = () => {
  if (!isValid.value) return
  emit('save', { ...formData.value })
  emit('update:open', false)
}

const handleCancel = () => {
  emit('update:open', false)
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

      <ColumnFormField
        v-model="formData"
        mode="full"
        :disabled="mode === 'edit'"
      />

      <div v-if="validationErrors.length > 0" class="rounded-lg bg-destructive/10 p-3 space-y-1">
        <p class="text-sm font-medium text-destructive">Validation Errors:</p>
        <ul class="list-disc list-inside text-xs text-destructive space-y-0.5">
          <li v-for="error in validationErrors" :key="error">{{ error }}</li>
        </ul>
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
