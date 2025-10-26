<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import {
  ArrowLeft,
  Download,
  Edit,
  Plus,
  RefreshCw,
  Save,
  Search,
  Trash2,
  Upload,
} from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

import { useTableDetail } from '@/composables/api/useTableDetail'
import { useFormatting } from '@/composables/formatting'

const route = useRoute()
const router = useRouter()
const { formatDateTime } = useFormatting()

const projectId = computed(() => route.params.projectId as string)
const tableId = computed(() => route.params.tableId as string)

const {
  data: tableDetail,
  loading,
  error,
  refresh,
} = useTableDetail(projectId.value, tableId.value)

const tableName = computed(() => tableDetail.value?.table.name ?? '')
const searchQuery = ref('')
const isAddDialogOpen = ref(false)
const isEditDialogOpen = ref(false)
const selectedRow = ref<Record<string, any> | null>(null)

interface TableColumn {
  key: string
  label: string
  type: string
  sortable: boolean
}

const columns = computed<TableColumn[]>(() => {
  if (!tableDetail.value) return []
  return tableDetail.value.columns.map((col) => ({
    key: col.name,
    label: col.name.charAt(0).toUpperCase() + col.name.slice(1).replace(/_/g, ' '),
    type: col.type,
    sortable: true,
  }))
})

const data = computed(() => tableDetail.value?.rows ?? [])

const newRow = ref<Record<string, any>>({})
const editRow = ref<Record<string, any>>({})
const sortColumn = ref('')
const sortDirection = ref<'asc' | 'desc'>('asc')
const currentPage = ref(1)
const pageSize = ref(10)

watch([projectId, tableId], () => {
  searchQuery.value = ''
  sortColumn.value = ''
  currentPage.value = 1
})

const filteredData = computed(() => {
  let filtered = data.value

  if (searchQuery.value) {
    filtered = filtered.filter((row) =>
      Object.values(row).some((value) =>
        String(value).toLowerCase().includes(searchQuery.value.toLowerCase())
      )
    )
  }

  if (sortColumn.value) {
    filtered = [...filtered].sort((a, b) => {
      const aVal = a[sortColumn.value] ?? ''
      const bVal = b[sortColumn.value] ?? ''

      if (sortDirection.value === 'asc') {
        return aVal > bVal ? 1 : -1
      } else {
        return aVal < bVal ? 1 : -1
      }
    })
  }

  return filtered
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(filteredData.value.length / pageSize.value)
})

const handleSort = (column: string) => {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = column
    sortDirection.value = 'asc'
  }
}

const handleAddRow = () => {
  const defaultRow: Record<string, any> = {}
  columns.value.forEach((col) => {
    defaultRow[col.key] = ''
  })
  newRow.value = defaultRow
  isAddDialogOpen.value = true
}

const handleEditRow = (row: Record<string, any>) => {
  editRow.value = { ...row }
  selectedRow.value = row
  isEditDialogOpen.value = true
}

const handleDeleteRow = async (rowId: string) => {
  if (confirm('Are you sure you want to delete this row?')) {
    try {
      console.warn('Delete functionality not yet implemented')
    } catch (error) {
      console.error('Failed to delete row:', error)
    }
  }
}

const saveNewRow = async () => {
  try {
    console.warn('Create functionality not yet implemented')
    isAddDialogOpen.value = false
    newRow.value = {}
  } catch (error) {
    console.error('Failed to create row:', error)
  }
}

const saveEditRow = async () => {
  try {
    console.warn('Update functionality not yet implemented')
    isEditDialogOpen.value = false
    editRow.value = {}
    selectedRow.value = null
  } catch (error) {
    console.error('Failed to update row:', error)
  }
}

const getStatusBadgeVariant = (status: string) => {
  switch (status) {
    case 'active':
      return 'default'
    case 'inactive':
      return 'secondary'
    case 'pending':
      return 'outline'
    default:
      return 'secondary'
  }
}
</script>

<template>
  <div class="p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <Button variant="ghost" size="sm" @click="router.push(`/projects/${projectId}/tables`)">
          <ArrowLeft class="w-4 h-4 mr-2" />
          Back to Tables
        </Button>
      </div>
    </div>

    <LoadingWrapper :is-loading="loading" loading-text="Loading table data...">
      <ErrorState v-if="error" :error="error" title="Failed to load table" :on-retry="refresh">
        <Button variant="default" size="sm" @click="router.push(`/projects/${projectId}/tables`)">
          Back to Tables
        </Button>
      </ErrorState>

      <template v-else>
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-foreground">{{ tableName }}</h1>
            <p class="text-muted-foreground">Manage table data and records</p>
          </div>
          <div class="flex items-center space-x-2">
            <Button variant="outline" size="sm">
              <Download class="w-4 h-4 mr-2" />
              Export
            </Button>
            <Button variant="outline" size="sm">
              <Upload class="w-4 h-4 mr-2" />
              Import
            </Button>
            <Button variant="outline" size="sm" @click="refresh">
              <RefreshCw class="w-4 h-4 mr-2" />
              Refresh
            </Button>
            <Button @click="handleAddRow">
              <Plus class="w-4 h-4 mr-2" />
              Add Row
            </Button>
          </div>
        </div>

        <Card>
          <CardContent class="p-4">
            <div class="flex items-center space-x-4">
              <div class="relative flex-1">
                <Search
                  class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4"
                />
                <Input v-model="searchQuery" placeholder="Search all columns..." class="pl-10" />
              </div>
              <Select v-model="pageSize">
                <SelectTrigger class="w-32">
                  <SelectValue placeholder="Page size" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="10">10 rows</SelectItem>
                  <SelectItem value="25">25 rows</SelectItem>
                  <SelectItem value="50">50 rows</SelectItem>
                  <SelectItem value="100">100 rows</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent class="p-0">
            <div class="overflow-x-auto">
              <table class="w-full">
                <thead class="border-b">
                  <tr>
                    <th
                      v-for="column in columns"
                      :key="column.key"
                      class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer hover:text-foreground transition-colors"
                      @click="column.sortable ? handleSort(column.key) : null"
                    >
                      <div class="flex items-center space-x-1">
                        <span>{{ column.label }}</span>
                        <span v-if="column.sortable" class="text-muted-foreground">
                          {{
                            sortColumn === column.key ? (sortDirection === 'asc' ? '↑' : '↓') : '↕'
                          }}
                        </span>
                      </div>
                    </th>
                    <th
                      class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                    >
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-border">
                  <tr
                    v-for="row in paginatedData"
                    :key="row.id"
                    class="hover:bg-muted/50 transition-colors"
                  >
                    <td
                      v-for="column in columns"
                      :key="column.key"
                      class="px-6 py-4 whitespace-nowrap text-sm text-foreground"
                    >
                      <span v-if="column.key === 'status'">
                        <Badge :variant="getStatusBadgeVariant(row[column.key])">
                          {{ row[column.key] }}
                        </Badge>
                      </span>
                      <span v-else-if="column.key === 'created_at'">
                        {{ formatDateTime(row[column.key]) }}
                      </span>
                      <span v-else>
                        {{ row[column.key] }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div class="flex items-center space-x-2">
                        <Button variant="ghost" size="sm" @click="handleEditRow(row)">
                          <Edit class="w-4 h-4" />
                        </Button>
                        <Button variant="ghost" size="sm" @click="handleDeleteRow(row.id)">
                          <Trash2 class="w-4 h-4" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>

        <div class="flex items-center justify-between">
          <div class="text-sm text-muted-foreground">
            Showing {{ (currentPage - 1) * pageSize + 1 }} to
            {{ Math.min(currentPage * pageSize, filteredData.length) }} of
            {{ filteredData.length }} results
          </div>
          <div class="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage === 1"
              @click="currentPage--"
            >
              Previous
            </Button>
            <span class="text-sm text-muted-foreground">
              Page {{ currentPage }} of {{ totalPages }}
            </span>
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage === totalPages"
              @click="currentPage++"
            >
              Next
            </Button>
          </div>
        </div>

        <Dialog v-model:open="isAddDialogOpen">
          <DialogContent class="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Add New Row</DialogTitle>
              <DialogDescription> Add a new record to the {{ tableName }} table </DialogDescription>
            </DialogHeader>
            <div class="space-y-4">
              <div v-for="column in columns" :key="column.key" class="space-y-2">
                <Label :for="column.key">{{ column.label }}</Label>
                <Input
                  :id="column.key"
                  v-model="newRow[column.key]"
                  :placeholder="`Enter ${column.label.toLowerCase()}`"
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" @click="isAddDialogOpen = false"> Cancel </Button>
              <Button @click="saveNewRow">
                <Save class="w-4 h-4 mr-2" />
                Save Row
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog v-model:open="isEditDialogOpen">
          <DialogContent class="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Edit Row</DialogTitle>
              <DialogDescription>
                Update the record in the {{ tableName }} table
              </DialogDescription>
            </DialogHeader>
            <div class="space-y-4">
              <div v-for="column in columns" :key="column.key" class="space-y-2">
                <Label :for="`edit-${column.key}`">{{ column.label }}</Label>
                <Input
                  :id="`edit-${column.key}`"
                  v-model="editRow[column.key]"
                  :placeholder="`Enter ${column.label.toLowerCase()}`"
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" @click="isEditDialogOpen = false"> Cancel </Button>
              <Button @click="saveEditRow">
                <Save class="w-4 h-4 mr-2" />
                Save Changes
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </template>
    </LoadingWrapper>
  </div>
</template>
