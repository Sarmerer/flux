<script setup lang="ts">
import { ChevronDown, Plus, Settings2, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuCheckboxItem,
} from '@/components/ui/dropdown-menu'

interface TableToolbarProps {
  totalRows: number
  selectedCount: number
  pageSize: number
  pageSizeOptions?: number[]
  columnVisibility?: any
  allColumns?: any[]
}

interface TableToolbarEmits {
  (e: 'insert-row'): void
  (e: 'delete-selected'): void
  (e: 'page-size-change', size: number): void
  (e: 'toggle-column', columnId: string, visible: boolean): void
}

const props = withDefaults(defineProps<TableToolbarProps>(), {
  pageSizeOptions: () => [10, 25, 50, 100, 250],
  allColumns: () => [],
})

const emit = defineEmits<TableToolbarEmits>()

const visibleColumns = () => {
  return props.allColumns?.filter((col) => col.getCanHide()) || []
}
</script>

<template>
  <div class="px-6 py-2.5 border-b bg-background">
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <div class="text-sm text-muted-foreground">
          {{ totalRows }} {{ totalRows === 1 ? 'row' : 'rows' }}
        </div>
        <div v-if="selectedCount > 0" class="flex items-center gap-2">
          <div class="text-sm font-medium">{{ selectedCount }} selected</div>
          <Button variant="destructive" size="sm" @click="emit('delete-selected')" class="h-7">
            <Trash2 class="w-3.5 h-3.5 mr-1.5" />
            Delete
          </Button>
        </div>
      </div>
      <div class="flex items-center gap-2">
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
      </div>
    </div>
  </div>
</template>
