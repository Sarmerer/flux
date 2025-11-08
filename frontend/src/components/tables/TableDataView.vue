<script setup lang="ts">
import { computed, h } from 'vue'
import type { ColumnDef } from '@tanstack/vue-table'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Key, Pencil, Trash2, ArrowUpDown } from 'lucide-vue-next'
import DataTable from './DataTable.vue'
import type { TableColumn } from '@/types'
import { useTableFormatters } from '@/composables/table/useTableFormatters'

interface TableDataViewProps {
  columns: TableColumn[]
  rows: any[]
  totalRows: number
  currentPage: number
  pageSize: number
  loading?: boolean
  enableRowSelection?: boolean
  onPageChange?: (page: number) => void
}

interface TableDataViewEmits {
  (e: 'page-change', page: number): void
  (e: 'edit-row', row: any): void
  (e: 'delete-row', row: any): void
  (e: 'row-selection-change', rows: any[]): void
}

const props = withDefaults(defineProps<TableDataViewProps>(), {
  loading: false,
  enableRowSelection: false,
})

const emit = defineEmits<TableDataViewEmits>()

const { formatCellValue } = useTableFormatters()

const tableColumns = computed<ColumnDef<any>[]>(() => {
  const cols: ColumnDef<any>[] = []

  if (props.enableRowSelection) {
    cols.push({
      id: 'select',
      header: ({ table }) =>
        h(Checkbox, {
          checked: table.getIsAllPageRowsSelected(),
          'onUpdate:checked': (value: boolean) => table.toggleAllPageRowsSelected(!!value),
          ariaLabel: 'Select all',
        }),
      cell: ({ row }) =>
        h(Checkbox, {
          checked: row.getIsSelected(),
          'onUpdate:checked': (value: boolean) => row.toggleSelected(!!value),
          ariaLabel: 'Select row',
        }),
      enableSorting: false,
      enableHiding: false,
    })
  }

  props.columns.forEach((column) => {
    cols.push({
      accessorKey: column.name,
      id: column.name,
      header: ({ column: col }) => {
        return h(
          'div',
          { class: 'flex items-center gap-1' },
          [
            h(
              Button,
              {
                variant: 'ghost',
                size: 'sm',
                class: '-ml-3 h-8 data-[state=open]:bg-accent hover:bg-transparent',
                onClick: () => col.toggleSorting(col.getIsSorted() === 'asc'),
              },
              () => [
                h('span', { class: 'font-medium' }, column.name),
                h(ArrowUpDown, { class: 'ml-1 h-3.5 w-3.5 opacity-50' }),
              ]
            ),
            column.is_primary_key
              ? h(Key, { class: 'w-3 h-3 text-yellow-500 ml-1', title: 'Primary Key' })
              : null,
          ]
        )
      },
      cell: ({ row }) => {
        const value = row.getValue(column.name)
        const formatted = formatCellValue(value, column.type)

        return h(
          'div',
          {
            class: 'max-w-[300px] truncate',
            title: String(value ?? '-'),
          },
          formatted
        )
      },
    })
  })

  cols.push({
    id: 'actions',
    header: () => h('div', { class: 'text-right pr-2' }, 'Actions'),
    cell: ({ row }) => {
      return h(
        'div',
        { class: 'flex items-center justify-end gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity pr-2' },
        [
          h(
            Button,
            {
              variant: 'ghost',
              size: 'sm',
              class: 'h-7 w-7 p-0 hover:bg-muted',
              onClick: () => emit('edit-row', row.original),
            },
            () => h(Pencil, { class: 'w-3.5 h-3.5' })
          ),
          h(
            Button,
            {
              variant: 'ghost',
              size: 'sm',
              class: 'h-7 w-7 p-0 text-destructive hover:text-destructive hover:bg-destructive/10',
              onClick: () => emit('delete-row', row.original),
            },
            () => h(Trash2, { class: 'w-3.5 h-3.5' })
          ),
        ]
      )
    },
    enableSorting: false,
    enableHiding: false,
  })

  return cols
})

const handlePageChange = (page: number) => {
  emit('page-change', page)
}

const handleRowSelectionChange = (rows: any[]) => {
  emit('row-selection-change', rows)
}
</script>

<template>
  <DataTable
    :columns="tableColumns"
    :data="rows"
    :total-rows="totalRows"
    :current-page="currentPage"
    :page-size="pageSize"
    :enable-row-selection="enableRowSelection"
    :enable-column-visibility="true"
    :enable-filtering="false"
    :on-page-change="onPageChange"
    @page-change="handlePageChange"
    @row-selection-change="handleRowSelectionChange"
  />
</template>
