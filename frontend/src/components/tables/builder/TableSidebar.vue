<script setup lang="ts">
import { computed, ref } from 'vue'
import { Plus, Search, Table as TableIcon, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import LoadingWrapper from '@/components/common/LoadingWrapper.vue'

interface Table {
  id: string
  name: string
  description?: string
}

interface TableSidebarProps {
  tables: Table[]
  loading: boolean
  selectedTableId?: string
}

interface TableSidebarEmits {
  (e: 'create'): void
  (e: 'select', tableId: string): void
  (e: 'delete', tableId: string): void
}

const props = defineProps<TableSidebarProps>()
const emit = defineEmits<TableSidebarEmits>()

const searchQuery = ref('')

const filteredTables = computed(() => {
  if (!props.tables || !Array.isArray(props.tables)) return []
  if (!searchQuery.value) return props.tables
  return props.tables.filter(
    (table) =>
      table.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (table.description && table.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )
})
</script>

<template>
  <div class="w-80 border-r flex flex-col bg-muted/10">
    <div class="p-4 border-b space-y-3">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">Tables</h2>
        <Button size="sm" @click="emit('create')">
          <Plus class="w-4 h-4" />
        </Button>
      </div>
      <div class="relative">
        <Search
          class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4"
        />
        <Input v-model="searchQuery" placeholder="Search tables..." class="pl-9 h-9" />
      </div>
    </div>

    <LoadingWrapper :is-loading="loading" loading-text="Loading tables...">
      <div class="flex-1 overflow-y-auto">
        <div class="p-2 space-y-1">
          <div
            v-for="table in filteredTables"
            :key="table.id"
            :class="[
              'group relative w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors',
              selectedTableId === table.id
                ? 'bg-primary/10 text-primary font-medium'
                : 'hover:bg-muted text-muted-foreground hover:text-foreground',
            ]"
          >
            <button
              @click="emit('select', table.id)"
              class="flex items-center gap-3 flex-1 min-w-0"
            >
              <TableIcon class="w-4 h-4 flex-shrink-0" />
              <span class="flex-1 text-left truncate">{{ table.name }}</span>
            </button>
            <Button
              variant="ghost"
              size="sm"
              @click.stop="emit('delete', table.id)"
              class="h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity text-destructive hover:text-destructive hover:bg-destructive/10"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>

        <div
          v-if="!loading && filteredTables.length === 0"
          class="p-4 text-center text-sm text-muted-foreground"
        >
          <p>No tables found</p>
        </div>
      </div>
    </LoadingWrapper>
  </div>
</template>
