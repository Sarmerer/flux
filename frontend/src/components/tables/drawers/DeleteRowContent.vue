<script setup lang="ts">
import { ref } from 'vue'

import { Button } from '@/components/ui/button'
import { SheetDescription, SheetFooter, SheetHeader, SheetTitle } from '@/components/ui/sheet'

import { useToast } from '@/composables/ui'

interface Props {
  rowId: string
  rowData: Record<string, any>
  onDelete: (rowId: string) => Promise<void>
}

interface Emits {
  (e: 'save'): void
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const toast = useToast()
const isDeleting = ref(false)

const handleClose = () => {
  emit('close')
}

const handleDelete = async () => {
  if (!props.rowId) {
    toast.error('Error', 'No row ID provided')
    return
  }

  isDeleting.value = true
  try {
    await props.onDelete(props.rowId)
    toast.success('Success', 'Row deleted successfully')
    emit('save')
    handleClose()
  } catch (error: any) {
    console.error('Failed to delete row:', error)
    toast.error('Error', error.message || 'Failed to delete row')
  } finally {
    isDeleting.value = false
  }
}
</script>

<template>
  <div class="flex flex-col h-full">
    <SheetHeader class="space-y-1 pb-4">
      <SheetTitle class="text-base font-medium">Delete Row</SheetTitle>
      <SheetDescription class="text-xs">
        Are you sure you want to delete this row? This action cannot be undone.
      </SheetDescription>
    </SheetHeader>

    <div class="flex-1 overflow-y-auto pb-4">
      <div v-if="rowData" class="space-y-2">
        <p class="text-xs text-muted-foreground">Row data:</p>
        <div class="max-h-96 overflow-y-auto bg-muted/50 rounded-md p-3 text-xs font-mono">
          <div v-for="(value, key) in rowData" :key="key" class="flex gap-2">
            <span class="text-muted-foreground">{{ key }}:</span>
            <span class="text-foreground break-all">{{ value ?? 'null' }}</span>
          </div>
        </div>
      </div>
    </div>

    <SheetFooter class="flex-row gap-2 pt-4 border-t">
      <Button variant="outline" size="sm" @click="handleClose" :disabled="isDeleting" class="flex-1">
        Cancel
      </Button>
      <Button variant="destructive" size="sm" @click="handleDelete" :disabled="isDeleting" class="flex-1">
        {{ isDeleting ? 'Deleting...' : 'Delete Row' }}
      </Button>
    </SheetFooter>
  </div>
</template>
