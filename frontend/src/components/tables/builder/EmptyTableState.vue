<script setup lang="ts">
import { Database, Plus } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

interface EmptyTableStateEmits {
  (e: 'create'): void
}

const emit = defineEmits<EmptyTableStateEmits>()
</script>

<template>
  <div class="flex-1 flex items-center justify-center">
    <div class="text-center space-y-4">
      <div class="flex justify-center">
        <div class="w-16 h-16 rounded-full bg-muted flex items-center justify-center">
          <Database class="w-8 h-8 text-muted-foreground" />
        </div>
      </div>
      <div>
        <h3 class="text-lg font-semibold mb-1">Select a table</h3>
        <p class="text-sm text-muted-foreground">
          Choose a table from the list to view and manage its data
        </p>
      </div>
      <Button @click="emit('create')">
        <Plus class="w-4 h-4 mr-2" />
        Create New Table
      </Button>
    </div>
  </div>
</template>
