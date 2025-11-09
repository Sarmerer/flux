<script setup lang="ts">
import { ref, computed } from 'vue'
import draggable from 'vuedraggable'
import ColumnFormField from '@/components/tables/ColumnFormField.vue'
import { Plus, Trash2, GripVertical, ChevronDown, ChevronRight } from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'

export interface Column {
  name: string
  type: string
  nullable: boolean
  primary_key: boolean
  default_value: string
  is_identity: boolean
  unique: boolean
}

interface Props {
  columns: Column[]
}

interface Emits {
  (e: 'update:columns', columns: Column[]): void
  (e: 'add-default-id'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const expandedColumns = ref<Set<number>>(new Set())

const localColumns = computed({
  get: () => props.columns,
  set: (value) => emit('update:columns', value),
})

const toggleColumn = (index: number) => {
  if (expandedColumns.value.has(index)) {
    expandedColumns.value.delete(index)
  } else {
    expandedColumns.value.add(index)
  }
}

const isExpanded = (index: number) => expandedColumns.value.has(index)

const addColumn = () => {
  const newColumn: Column = {
    name: 'column_name',
    type: 'TEXT',
    nullable: true,
    default_value: '',
    is_identity: false,
    unique: false,
    primary_key: false,
  }
  const newColumns = [...props.columns, newColumn]
  emit('update:columns', newColumns)
  expandedColumns.value.add(newColumns.length - 1)
}

const removeColumn = (index: number) => {
  const updated = props.columns.filter((_, i) => i !== index)
  emit('update:columns', updated)
  expandedColumns.value.delete(index)
}

const updateColumn = (index: number, column: Column) => {
  const updated = [...props.columns]
  updated[index] = column
  emit('update:columns', updated)
}

const existingColumnNames = (excludeIndex?: number) => {
  return props.columns
    .map((c, i) => (i !== excludeIndex ? c.name.toLowerCase() : null))
    .filter(Boolean) as string[]
}

const getColumnBadges = (column: Column) => {
  const badges: string[] = []
  if (column.primary_key) badges.push('Primary')
  if (column.unique) badges.push('Unique')
  if (column.is_identity) badges.push('Auto')
  if (!column.nullable) badges.push('Not Null')
  return badges
}
</script>

<template>
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <Label class="text-xs font-medium">Columns</Label>
      <div class="flex gap-1.5">
        <Button
          type="button"
          variant="outline"
          size="sm"
          @click="emit('add-default-id')"
          class="h-7 text-xs px-2"
        >
          Add ID Column
        </Button>
        <Button type="button" size="sm" @click="addColumn" class="h-7 text-xs px-2">
          <Plus class="w-3 h-3 mr-1" />
          Add Column
        </Button>
      </div>
    </div>

    <div
      v-if="columns.length === 0"
      class="text-center py-6 text-xs text-muted-foreground border-2 border-dashed rounded-md"
    >
      <p>No columns added yet.</p>
      <p class="text-[11px] mt-1">Click "Add Column" or "Add ID Column" to get started.</p>
    </div>

    <draggable
      v-else
      v-model="localColumns"
      item-key="name"
      handle=".drag-handle"
      class="space-y-1.5"
    >
      <template #item="{ element: column, index }">
        <div class="border rounded-md bg-muted/10 hover:bg-muted/20 transition-colors">
          <div class="flex items-center gap-2 px-2.5 py-1.5">
            <GripVertical
              class="w-3.5 h-3.5 text-muted-foreground flex-shrink-0 cursor-grab active:cursor-grabbing drag-handle"
            />

            <component
              :is="isExpanded(index) ? ChevronDown : ChevronRight"
              class="w-3.5 h-3.5 text-muted-foreground flex-shrink-0 cursor-pointer"
              @click="toggleColumn(index)"
            />

            <div
              class="flex-1 min-w-0 flex items-center gap-1.5 text-xs cursor-pointer"
              @click="toggleColumn(index)"
            >
              <span class="font-medium truncate">{{ column.name || 'Unnamed' }}</span>
              <span class="text-muted-foreground">•</span>
              <span class="text-muted-foreground">{{ column.type }}</span>

              <div v-if="getColumnBadges(column).length > 0" class="flex gap-1 ml-1">
                <span
                  v-for="badge in getColumnBadges(column)"
                  :key="badge"
                  class="text-[10px] px-1.5 py-0.5 rounded bg-primary/10 text-primary"
                >
                  {{ badge }}
                </span>
              </div>
            </div>

            <Button
              type="button"
              variant="ghost"
              size="sm"
              @click.stop="removeColumn(index)"
              class="h-6 w-6 p-0 text-destructive hover:text-destructive hover:bg-destructive/10 flex-shrink-0"
            >
              <Trash2 class="w-3 h-3" />
            </Button>
          </div>

          <div v-if="isExpanded(index)" class="px-2.5 pb-2.5 pt-1">
            <ColumnFormField
              :model-value="column"
              @update:model-value="updateColumn(index, $event)"
              mode="inline"
              :existing-columns="existingColumnNames(index)"
              :show-labels="false"
            />
          </div>
        </div>
      </template>
    </draggable>
  </div>
</template>
