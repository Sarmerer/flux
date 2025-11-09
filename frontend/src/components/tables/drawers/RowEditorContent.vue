<script setup lang="ts">
import { ref, watch, computed } from 'vue'

import type { TableColumn } from '@/types'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { SheetDescription, SheetFooter, SheetHeader, SheetTitle } from '@/components/ui/sheet'

import { useToast } from '@/composables/ui'
import { useColumnTypeUtils } from '@/composables/tables/useColumnTypeUtils'

interface Props {
  mode: 'insert' | 'edit'
  columns: TableColumn[]
  rowData?: Record<string, any>
  rowId?: string
  onSave: (data: Record<string, any>, rowId?: string) => Promise<void>
}

interface Emits {
  (e: 'save', data: Record<string, any>): void
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const toast = useToast()
const { convertValueByType, getInputType, isTextArea, validateColumnValue } = useColumnTypeUtils()

const isSubmitting = ref(false)
const formData = ref<Record<string, any>>({})

const title = computed(() => (props.mode === 'insert' ? 'Insert Row' : 'Edit Row'))
const submitLabel = computed(() => (props.mode === 'insert' ? 'Insert Row' : 'Update Row'))
const submittingLabel = computed(() => (props.mode === 'insert' ? 'Inserting...' : 'Updating...'))

const resetForm = () => {
  if (props.mode === 'edit' && props.rowData) {
    formData.value = { ...props.rowData }
  } else {
    formData.value = {}
    props.columns.forEach((col) => {
      formData.value[col.name] = col.default_value || ''
    })
  }
}

watch([() => props.columns, () => props.rowData], resetForm, { immediate: true })

const handleClose = () => {
  emit('close')
}

const handleSubmit = async () => {
  const data: Record<string, any> = {}

  for (const col of props.columns) {
    const value = formData.value[col.name]

    const error = validateColumnValue(value, {
      type: col.type,
      isNullable: col.is_nullable,
      name: col.name,
    })

    if (error) {
      toast.error('Validation Error', error)
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
    await props.onSave(data, props.rowId)
    const successMessage =
      props.mode === 'insert' ? 'Row inserted successfully' : 'Row updated successfully'
    toast.success('Success', successMessage)
    emit('save', data)
    handleClose()
  } catch (error: any) {
    console.error(`Failed to ${props.mode} row:`, error)
    toast.error('Error', error.message || `Failed to ${props.mode} row`)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="flex flex-col h-full">
    <SheetHeader class="space-y-1 pb-4">
      <SheetTitle class="text-base font-medium">{{ title }}</SheetTitle>
      <SheetDescription class="text-xs">
        {{ mode === 'insert' ? 'Add a new row to the table.' : 'Update the values for this row.' }}
        Fill in the values for each column.
      </SheetDescription>
    </SheetHeader>

    <div class="flex-1 overflow-y-auto space-y-4 pb-4">
      <div v-for="column in columns" :key="column.name" class="space-y-2">
        <Label :for="column.name" class="flex items-center gap-2 text-xs">
          {{ column.name }}
          <span class="text-[11px] text-muted-foreground">{{ column.type }}</span>
          <span v-if="!column.is_nullable" class="text-[11px] text-destructive">*</span>
        </Label>

        <textarea
          v-if="isTextArea(column.type)"
          :id="column.name"
          v-model="formData[column.name]"
          :placeholder="column.default_value || `Enter ${column.name}`"
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
          <span class="text-xs text-muted-foreground">Check for true</span>
        </div>

        <Input
          v-else
          :id="column.name"
          v-model="formData[column.name]"
          :type="getInputType(column.type)"
          :placeholder="column.default_value || `Enter ${column.name}`"
          :step="column.type.toUpperCase().includes('DECIMAL') ? '0.01' : undefined"
          class="h-8 text-sm"
        />
      </div>
    </div>

    <SheetFooter class="flex-row gap-2 pt-4 border-t">
      <Button
        variant="outline"
        size="sm"
        @click="handleClose"
        :disabled="isSubmitting"
        class="flex-1"
      >
        Cancel
      </Button>
      <Button size="sm" @click="handleSubmit" :disabled="isSubmitting" class="flex-1">
        {{ isSubmitting ? submittingLabel : submitLabel }}
      </Button>
    </SheetFooter>
  </div>
</template>
