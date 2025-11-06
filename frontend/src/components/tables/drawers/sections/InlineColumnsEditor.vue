<script setup lang="ts">
import { Plus, Trash2 } from 'lucide-vue-next'

import ColumnFormField from '@/components/tables/ColumnFormField.vue'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'

export interface Column {
  name: string
  type: string
  nullable: boolean
  primary_key?: boolean
  is_primary_key?: boolean
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

const addColumn = () => {
  const newColumn: Column = {
    name: 'column_name',
    type: 'TEXT',
    nullable: true,
    default_value: '',
    is_identity: false,
    unique: false,
    is_primary_key: false,
  }
  emit('update:columns', [...props.columns, newColumn])
}

const removeColumn = (index: number) => {
  const updated = props.columns.filter((_, i) => i !== index)
  emit('update:columns', updated)
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
</script>

<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between">
      <Label class="text-xs font-medium">Columns</Label>
      <div class="flex gap-1.5">
        <Button type="button" variant="outline" size="sm" @click="emit('add-default-id')" class="h-7 text-xs px-2">
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
      class="text-center py-8 text-xs text-muted-foreground border-2 border-dashed rounded-md"
    >
      <p>No columns added yet.</p>
      <p class="text-[11px] mt-1">Click "Add Column" or "Add ID Column" to get started.</p>
    </div>

    <div v-else class="space-y-2">
      <div v-for="(column, index) in columns" :key="index" class="border rounded-md p-3 bg-muted/20 hover:bg-muted/30 transition-colors">
        <div class="flex items-start gap-2 mb-3">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="text-xs font-medium truncate">{{ column.name || 'Unnamed' }}</span>
              <span class="text-[11px] text-muted-foreground">{{ column.type }}</span>
              <span v-if="column.primary_key || column.is_primary_key" class="text-[11px] px-1.5 py-0.5 rounded bg-primary/10 text-primary">Primary</span>
            </div>
            <div v-if="column.default_value" class="text-[11px] text-muted-foreground truncate">
              Default: {{ column.default_value }}
            </div>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            @click="removeColumn(index)"
            class="h-6 w-6 p-0 text-destructive hover:text-destructive hover:bg-destructive/10"
          >
            <Trash2 class="w-3.5 h-3.5" />
          </Button>
        </div>

        <ColumnFormField
          :model-value="column"
          @update:model-value="updateColumn(index, $event)"
          mode="inline"
          :existing-columns="existingColumnNames(index)"
        />
      </div>
    </div>
  </div>
</template>
