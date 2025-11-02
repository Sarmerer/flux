<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useToast } from '@/composables/ui'

const props = defineProps<{
  open: boolean
  rowId: string | null
  rowData: Record<string, any> | null
  onDelete: (rowId: string) => Promise<void>
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const toast = useToast()
const isDeleting = ref(false)

const handleClose = () => {
  emit('update:open', false)
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
  <Dialog :open="open" @update:open="handleClose">
    <DialogContent class="max-w-md">
      <DialogHeader>
        <DialogTitle>Delete Row</DialogTitle>
        <DialogDescription>
          Are you sure you want to delete this row? This action cannot be undone.
        </DialogDescription>
      </DialogHeader>

      <div v-if="rowData" class="py-4 space-y-2">
        <p class="text-sm text-muted-foreground">Row data:</p>
        <div class="max-h-40 overflow-y-auto bg-muted/50 rounded-md p-3 text-xs font-mono">
          <div v-for="(value, key) in rowData" :key="key" class="flex gap-2">
            <span class="text-muted-foreground">{{ key }}:</span>
            <span class="text-foreground">{{ value ?? 'null' }}</span>
          </div>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="handleClose" :disabled="isDeleting">
          Cancel
        </Button>
        <Button variant="destructive" @click="handleDelete" :disabled="isDeleting">
          {{ isDeleting ? 'Deleting...' : 'Delete Row' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
