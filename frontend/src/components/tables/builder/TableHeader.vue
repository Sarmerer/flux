<script setup lang="ts">
import {
  ChevronDown,
  MoreVertical,
  Pencil,
  Plus,
  RefreshCw,
  Settings2,
  Trash2,
  Trash2 as Trash2Icon,
} from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

interface TableHeaderProps {
  tableName: string
  tableDescription?: string
  loading?: boolean
  totalRows: number
  selectedCount: number
  pageSize: number
  pageSizeOptions?: number[]
  allColumns?: any[]
}

interface TableHeaderEmits {
  (e: 'edit-schema'): void
  (e: 'refresh'): void
  (e: 'delete'): void
  (e: 'insert-row'): void
  (e: 'delete-selected'): void
  (e: 'page-size-change', size: number): void
  (e: 'toggle-column', columnId: string, visible: boolean): void
}

const props = withDefaults(defineProps<TableHeaderProps>(), {
  pageSizeOptions: () => [10, 25, 50, 100, 250],
  allColumns: () => [],
})

const emit = defineEmits<TableHeaderEmits>()

const visibleColumns = () => {
  return props.allColumns?.filter((col) => col.getCanHide()) || []
}
</script>

<template>
  <div class="px-6 py-3 bg-background">
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <div>
          <div class="flex items-baseline gap-2">
            <h1 class="text-xl font-semibold">{{ tableName }}</h1>
            <span class="text-sm text-muted-foreground">
              {{ totalRows }} {{ totalRows === 1 ? 'row' : 'rows' }}
            </span>
          </div>
          <p v-if="tableDescription" class="text-xs text-muted-foreground mt-0.5">
            {{ tableDescription }}
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div v-if="selectedCount > 0" class="flex items-center gap-2 mr-2">
          <div class="text-sm font-medium">{{ selectedCount }} selected</div>
          <Button variant="destructive" size="sm" @click="emit('delete-selected')" class="h-8">
            <Trash2Icon class="w-3.5 h-3.5 mr-1.5" />
            Delete
          </Button>
        </div>

        <DropdownMenu v-if="allColumns && allColumns.length > 0">
          <DropdownMenuTrigger as-child>
            <Button variant="outline" size="sm" class="h-8">
              <Settings2 class="w-4 h-4 mr-2" />
              View
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-44">
            <DropdownMenuCheckboxItem
              v-for="column in visibleColumns()"
              :key="column.id"
              :checked="column.getIsVisible()"
              @update:checked="(value: boolean) => emit('toggle-column', column.id, value)"
              class="capitalize"
            >
              {{ column.id }}
            </DropdownMenuCheckboxItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="outline" size="sm" class="h-8">
              {{ pageSize }} rows
              <ChevronDown class="w-3.5 h-3.5 ml-1.5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem
              v-for="size in pageSizeOptions"
              :key="size"
              @click="emit('page-size-change', size)"
            >
              {{ size }} rows
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <Button variant="default" size="sm" @click="emit('insert-row')" class="h-8">
          <Plus class="w-4 h-4 mr-2" />
          Insert Row
        </Button>

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
            <DropdownMenuItem
              @click="emit('delete')"
              class="text-destructive focus:text-destructive"
            >
              <Trash2 class="w-4 h-4 mr-2" />
              Delete Table
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  </div>
</template>
