<script setup lang="ts">
import { MoreVertical, Pencil, RefreshCw, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'

interface TableHeaderProps {
  tableName: string
  tableDescription?: string
  loading?: boolean
}

interface TableHeaderEmits {
  (e: 'edit-schema'): void
  (e: 'refresh'): void
  (e: 'delete'): void
}

defineProps<TableHeaderProps>()
const emit = defineEmits<TableHeaderEmits>()
</script>

<template>
  <div class="border-b px-6 py-3 bg-background">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold">{{ tableName }}</h1>
        <p v-if="tableDescription" class="text-xs text-muted-foreground mt-0.5">
          {{ tableDescription }}
        </p>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" size="sm" class="h-8 w-8 p-0">
            <MoreVertical class="w-4 h-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem @click="emit('edit-schema')">
            <Pencil class="w-4 h-4 mr-2" />
            Edit Schema
          </DropdownMenuItem>
          <DropdownMenuItem @click="emit('refresh')" :disabled="loading">
            <RefreshCw class="w-4 h-4 mr-2" />
            Refresh Data
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="emit('delete')" class="text-destructive focus:text-destructive">
            <Trash2 class="w-4 h-4 mr-2" />
            Delete Table
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </div>
</template>
