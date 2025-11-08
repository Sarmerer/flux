<script setup lang="ts" generic="TData">
import { computed, ref } from 'vue'
import {
  FlexRender,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  getFilteredRowModel,
  useVueTable,
  type ColumnDef,
  type SortingState,
  type VisibilityState,
  type RowSelectionState,
  type ColumnFiltersState,
} from '@tanstack/vue-table'
import { valueUpdater } from '@/lib/utils'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  TableEmpty,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { ChevronDown, ChevronLeft, ChevronRight, Settings2 } from 'lucide-vue-next'

interface DataTableProps {
  columns: ColumnDef<TData>[]
  data: TData[]
  pageSize?: number
  totalRows?: number
  currentPage?: number
  enablePagination?: boolean
  enableSorting?: boolean
  enableFiltering?: boolean
  enableColumnVisibility?: boolean
  enableRowSelection?: boolean
  searchPlaceholder?: string
  filterKey?: string
  onPageChange?: (page: number) => void
  onRowSelectionChange?: (selectedRows: TData[]) => void
}

const props = withDefaults(defineProps<DataTableProps>(), {
  pageSize: 50,
  totalRows: 0,
  currentPage: 1,
  enablePagination: true,
  enableSorting: true,
  enableFiltering: true,
  enableColumnVisibility: true,
  enableRowSelection: false,
  searchPlaceholder: 'Filter...',
  filterKey: '',
})

const emit = defineEmits<{
  pageChange: [page: number]
  rowSelectionChange: [selectedRows: TData[]]
}>()

const sorting = ref<SortingState>([])
const columnFilters = ref<ColumnFiltersState>([])
const columnVisibility = ref<VisibilityState>({})
const rowSelection = ref<RowSelectionState>({})

const isServerSidePagination = computed(() => !!props.onPageChange)

const table = useVueTable({
  get data() {
    return props.data
  },
  get columns() {
    return props.columns
  },
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: !isServerSidePagination.value && props.enablePagination ? getPaginationRowModel() : undefined,
  getSortedRowModel: props.enableSorting ? getSortedRowModel() : undefined,
  getFilteredRowModel: props.enableFiltering ? getFilteredRowModel() : undefined,
  onSortingChange: (updaterOrValue) => valueUpdater(updaterOrValue, sorting),
  onColumnFiltersChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnFilters),
  onColumnVisibilityChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnVisibility),
  onRowSelectionChange: (updaterOrValue) => {
    valueUpdater(updaterOrValue, rowSelection)
    const selectedRows = table.getFilteredSelectedRowModel().rows.map((row) => row.original)
    emit('rowSelectionChange', selectedRows)
  },
  state: {
    get sorting() {
      return sorting.value
    },
    get columnFilters() {
      return columnFilters.value
    },
    get columnVisibility() {
      return columnVisibility.value
    },
    get rowSelection() {
      return rowSelection.value
    },
  },
  manualPagination: isServerSidePagination.value,
  pageCount: isServerSidePagination.value
    ? Math.ceil(props.totalRows / props.pageSize)
    : undefined,
})

const goToPreviousPage = () => {
  if (props.onPageChange) {
    emit('pageChange', props.currentPage - 1)
  } else {
    table.previousPage()
  }
}

const goToNextPage = () => {
  if (props.onPageChange) {
    emit('pageChange', props.currentPage + 1)
  } else {
    table.nextPage()
  }
}

const canGoToPreviousPage = computed(() => {
  if (props.onPageChange) {
    return props.currentPage > 1
  }
  return table.getCanPreviousPage()
})

const canGoToNextPage = computed(() => {
  if (props.onPageChange) {
    const totalPages = Math.ceil(props.totalRows / props.pageSize)
    return props.currentPage < totalPages
  }
  return table.getCanNextPage()
})

const currentPageNumber = computed(() => {
  if (props.onPageChange) {
    return props.currentPage
  }
  return table.getState().pagination.pageIndex + 1
})

const totalPages = computed(() => {
  if (props.onPageChange) {
    return Math.ceil(props.totalRows / props.pageSize)
  }
  return table.getPageCount()
})

const startIndex = computed(() => {
  if (props.onPageChange) {
    return (props.currentPage - 1) * props.pageSize + 1
  }
  const state = table.getState().pagination
  return state.pageIndex * state.pageSize + 1
})

const endIndex = computed(() => {
  if (props.onPageChange) {
    return Math.min(props.currentPage * props.pageSize, props.totalRows)
  }
  const state = table.getState().pagination
  return Math.min((state.pageIndex + 1) * state.pageSize, props.data.length)
})

const totalRowsCount = computed(() => {
  return props.onPageChange ? props.totalRows : props.data.length
})
</script>

<template>
  <div class="space-y-4">
    <div v-if="enableFiltering || enableColumnVisibility" class="flex items-center justify-between gap-2 mb-4">
      <Input
        v-if="enableFiltering && filterKey"
        :placeholder="searchPlaceholder"
        :model-value="table.getColumn(filterKey)?.getFilterValue() as string"
        @update:model-value="table.getColumn(filterKey)?.setFilterValue($event)"
        class="max-w-sm h-9"
      />
      <div v-else />

      <DropdownMenu v-if="enableColumnVisibility">
        <DropdownMenuTrigger as-child>
          <Button variant="outline" size="sm" class="h-9">
            <Settings2 class="w-4 h-4 mr-2" />
            View
            <ChevronDown class="w-4 h-4 ml-2" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-44">
          <DropdownMenuCheckboxItem
            v-for="column in table.getAllColumns().filter((col) => col.getCanHide())"
            :key="column.id"
            :checked="column.getIsVisible()"
            @update:checked="(value: boolean) => column.toggleVisibility(!!value)"
            class="capitalize"
          >
            {{ column.id }}
          </DropdownMenuCheckboxItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <div class="border rounded-lg overflow-hidden">
      <div class="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
              <TableHead v-for="header in headerGroup.headers" :key="header.id">
                <FlexRender
                  v-if="!header.isPlaceholder"
                  :render="header.column.columnDef.header"
                  :props="header.getContext()"
                />
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <template v-if="table.getRowModel().rows?.length">
              <TableRow
                v-for="row in table.getRowModel().rows"
                :key="row.id"
                :data-state="row.getIsSelected() ? 'selected' : undefined"
                class="group"
              >
                <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                  <FlexRender
                    :render="cell.column.columnDef.cell"
                    :props="cell.getContext()"
                  />
                </TableCell>
              </TableRow>
            </template>
            <TableEmpty v-else :colspan="table.getAllColumns().length">
              No results found.
            </TableEmpty>
          </TableBody>
        </Table>
      </div>
    </div>

    <div v-if="enablePagination && totalRowsCount > 0" class="flex items-center justify-between pt-2">
      <div class="text-sm text-muted-foreground">
        Showing {{ startIndex }} to {{ endIndex }} of {{ totalRowsCount }} rows
      </div>
      <div class="flex items-center gap-2">
        <div class="text-sm text-muted-foreground">
          Page {{ currentPageNumber }} of {{ totalPages || 1 }}
        </div>
        <div class="flex items-center gap-1">
          <Button
            variant="outline"
            size="sm"
            class="h-8 w-8 p-0"
            @click="goToPreviousPage"
            :disabled="!canGoToPreviousPage"
          >
            <ChevronLeft class="w-4 h-4" />
          </Button>
          <Button
            variant="outline"
            size="sm"
            class="h-8 w-8 p-0"
            @click="goToNextPage"
            :disabled="!canGoToNextPage"
          >
            <ChevronRight class="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
