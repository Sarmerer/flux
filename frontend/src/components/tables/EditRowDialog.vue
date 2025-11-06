<script setup lang="ts">
import { ref, watch } from 'vue'

import type { TableColumn } from '@/types/table'

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

import { useToast } from '@/composables/ui'

const props = defineProps<{
  open: boolean
  columns: TableColumn[]
  rowData: Record<string, any> | null
  rowId: string | null
  onUpdate: (rowId: string, data: Record<string, any>) => Promise<void>
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const toast = useToast()
const isSubmitting = ref(false)
const formData = ref<Record<string, any>>({})

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.rowData) {
      resetForm()
    }
  }
)

const resetForm = () => {
  formData.value = { ...props.rowData }
}

const handleClose = () => {
  emit('update:open', false)
}

const handleSubmit = async () => {
  if (!props.rowId) {
    toast.error('Error', 'No row ID provided')
    return
  }

  const data: Record<string, any> = {}

  for (const col of props.columns) {
    const value = formData.value[col.name]

    if (!col.is_nullable && (value === '' || value === null || value === undefined)) {
      toast.error('Validation Error', `${col.name} is required`)
      return
    }

    if (value === '') {
      data[col.name] = null
    } else {
      data[col.name] = convertValueByType(value, col.type)
    }
  }

  isSubmitting.value = true
  try {
    await props.onUpdate(props.rowId, data)
    toast.success('Success', 'Row updated successfully')
    handleClose()
  } catch (error: any) {
    console.error('Failed to update row:', error)
    toast.error('Error', error.message || 'Failed to update row')
  } finally {
    isSubmitting.value = false
  }
}

const convertValueByType = (value: any, type: string): any => {
  if (value === null || value === '') return null

  const upperType = type.toUpperCase()

  if (upperType.includes('INT') || upperType.includes('SERIAL')) {
    return parseInt(value, 10)
  }

  if (
    upperType.includes('DECIMAL') ||
    upperType.includes('NUMERIC') ||
    upperType.includes('REAL') ||
    upperType.includes('DOUBLE')
  ) {
    return parseFloat(value)
  }

  if (upperType === 'BOOLEAN' || upperType === 'BOOL') {
    if (typeof value === 'boolean') return value
    if (typeof value === 'string') {
      const lower = value.toLowerCase()
      return lower === 'true' || lower === '1' || lower === 'yes'
    }
    return Boolean(value)
  }

  if (upperType === 'JSON' || upperType === 'JSONB') {
    if (typeof value === 'string') {
      try {
        return JSON.parse(value)
      } catch {
        return value
      }
    }
    return value
  }

  return value
}

const getInputType = (columnType: string): string => {
  const upperType = columnType.toUpperCase()

  if (upperType.includes('INT') || upperType.includes('SERIAL')) {
    return 'number'
  }

  if (
    upperType.includes('DECIMAL') ||
    upperType.includes('NUMERIC') ||
    upperType.includes('REAL') ||
    upperType.includes('DOUBLE')
  ) {
    return 'number'
  }

  if (upperType === 'BOOLEAN' || upperType === 'BOOL') {
    return 'checkbox'
  }

  if (upperType === 'DATE') {
    return 'date'
  }

  if (upperType.includes('TIMESTAMP') || upperType.includes('TIME')) {
    return 'datetime-local'
  }

  return 'text'
}

const isTextArea = (columnType: string): boolean => {
  const upperType = columnType.toUpperCase()
  return upperType === 'TEXT' || upperType === 'JSON' || upperType === 'JSONB'
}

const formatValueForInput = (value: any, columnType: string): any => {
  if (value === null || value === undefined) return ''

  const upperType = columnType.toUpperCase()

  if (upperType === 'JSON' || upperType === 'JSONB') {
    if (typeof value === 'object') {
      return JSON.stringify(value, null, 2)
    }
    return value
  }

  if (upperType.includes('TIMESTAMP')) {
    if (typeof value === 'string') {
      try {
        const date = new Date(value)
        return date.toISOString().slice(0, 16)
      } catch {
        return value
      }
    }
  }

  if (upperType === 'DATE') {
    if (typeof value === 'string') {
      return value.split('T')[0]
    }
  }

  return value
}
</script>

<template>
  <Dialog :open="open" @update:open="handleClose">
    <DialogContent class="max-w-2xl max-h-[80vh] flex flex-col">
      <DialogHeader>
        <DialogTitle>Edit Row</DialogTitle>
        <DialogDescription>
          Update the values for this row. Changes will be saved immediately.
        </DialogDescription>
      </DialogHeader>

      <div class="flex-1 overflow-y-auto space-y-4 py-4">
        <div v-for="column in columns" :key="column.name" class="space-y-2">
          <Label :for="column.name" class="flex items-center gap-2">
            {{ column.name }}
            <span class="text-xs text-muted-foreground">{{ column.type }}</span>
            <span v-if="!column.is_nullable" class="text-xs text-destructive">*</span>
          </Label>

          <textarea
            v-if="isTextArea(column.type)"
            :id="column.name"
            v-model="formData[column.name]"
            :placeholder="`Enter ${column.name}`"
            class="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          />

          <div
            v-else-if="getInputType(column.type) === 'checkbox'"
            class="flex items-center space-x-2"
          >
            <input
              :id="column.name"
              v-model="formData[column.name]"
              type="checkbox"
              class="rounded"
            />
            <span class="text-sm text-muted-foreground">Check for true</span>
          </div>

          <Input
            v-else
            :id="column.name"
            v-model="formData[column.name]"
            :type="getInputType(column.type)"
            :placeholder="`Enter ${column.name}`"
            :step="column.type.toUpperCase().includes('DECIMAL') ? '0.01' : undefined"
          />
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="handleClose" :disabled="isSubmitting"> Cancel </Button>
        <Button @click="handleSubmit" :disabled="isSubmitting">
          {{ isSubmitting ? 'Updating...' : 'Update Row' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
