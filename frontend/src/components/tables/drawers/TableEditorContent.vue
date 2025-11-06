<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import { Button } from '@/components/ui/button'
import { SheetDescription, SheetFooter, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import InlineColumnsEditor from './sections/InlineColumnsEditor.vue'
import TableMetadataSection from './sections/TableMetadataSection.vue'

export interface Column {
  name: string
  type: string
  nullable: boolean
  primary_key?: boolean
  is_primary_key?: boolean
  default_value: string
  is_identity: boolean
  unique: boolean
}

export interface ForeignKey {
  id?: string
  column: string
  referencedSchema: string
  referencedTable: string
  referencedColumn: string
  onUpdate: string
  onDelete: string
}

export interface TableFormData {
  name: string
  description: string
  columns: Column[]
  foreignKeys: ForeignKey[]
}

interface Props {
  mode?: 'create' | 'edit'
  tableId?: string
  initialData?: TableFormData
}

interface Emits {
  (e: 'save', data: TableFormData): void
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'create',
})

const emit = defineEmits<Emits>()

const formData = ref<TableFormData>({
  name: '',
  description: '',
  columns: [],
  foreignKeys: [],
})

const validationErrors = ref<Record<string, string>>({})

watch(
  () => props.initialData,
  (data) => {
    if (data) {
      formData.value = { ...data }
    }
  },
  { immediate: true }
)

const hasValidationErrors = computed(() => Object.keys(validationErrors.value).length > 0)

const validate = (): boolean => {
  const errors: Record<string, string> = {}

  if (!formData.value.name.trim()) {
    errors.name = 'Table name is required'
  } else if (!/^[a-z_][a-z0-9_]*$/.test(formData.value.name)) {
    errors.name =
      'Table name must start with a letter or underscore and contain only lowercase letters, numbers, and underscores'
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

    const duplicates = formData.value.columns.filter(
      (c) => c.name.toLowerCase() === col.name.toLowerCase()
    )
    if (duplicates.length > 1) {
      errors[`column_${index}_name`] = 'Duplicate column name'
    }
  })

  validationErrors.value = errors
  return Object.keys(errors).length === 0
}

const handleSave = () => {
  if (!validate()) return
  emit('save', formData.value)
}

const handleCancel = () => {
  emit('close')
}

const addDefaultIdColumn = () => {
  const hasIdColumn = formData.value.columns.some((col) => col.name === 'id')
  if (hasIdColumn) return

  formData.value.columns = [
    {
      name: 'id',
      type: 'UUID',
      nullable: false,
      primary_key: true,
      default_value: 'gen_random_uuid()',
      is_identity: false,
      unique: true,
    },
    ...formData.value.columns,
  ]
  validationErrors.value = {}
}

const updateColumns = (columns: Column[]) => {
  formData.value.columns = columns
  validationErrors.value = {}
}
</script>

<template>
  <div class="flex flex-col h-full">
    <SheetHeader class="space-y-1 pb-4">
      <SheetTitle class="text-base font-medium">{{ mode === 'create' ? 'Create New Table' : `Update table ${formData.name}` }}</SheetTitle>
      <SheetDescription class="text-xs">
        {{
          mode === 'create'
            ? 'Create a new table with custom columns and relationships.'
            : 'Update table structure, columns, and relationships.'
        }}
      </SheetDescription>
    </SheetHeader>

    <div class="flex-1 overflow-y-auto space-y-4 pb-4">
      <TableMetadataSection
        v-model="formData"
        :validation-errors="validationErrors"
        :disabled="mode === 'edit'"
      />

      <div class="border-t pt-4">
        <InlineColumnsEditor
          :columns="formData.columns"
          @update:columns="updateColumns"
          @add-default-id="addDefaultIdColumn"
        />
        <p v-if="validationErrors.columns" class="text-xs text-red-500 mt-2">
          {{ validationErrors.columns }}
        </p>
      </div>

    </div>

    <SheetFooter class="flex-row gap-2 pt-4 border-t">
      <Button variant="outline" size="sm" @click="handleCancel" class="flex-1">Cancel</Button>
      <Button size="sm" @click="handleSave" :disabled="hasValidationErrors && formData.columns.length > 0" class="flex-1">
        {{ mode === 'create' ? 'Create Table' : 'Save Changes' }}
      </Button>
    </SheetFooter>
  </div>
</template>
