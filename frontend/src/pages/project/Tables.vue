<script setup lang="ts">
import DeleteRowDialog from '@/components/tables/DeleteRowDialog.vue'
import EditRowDialog from '@/components/tables/EditRowDialog.vue'
import InsertRowDialog from '@/components/tables/InsertRowDialog.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import {
  ChevronLeft,
  ChevronRight,
  Database,
  Pencil,
  Plus,
  Search,
  Settings2,
  Table as TableIcon,
  Trash2,
} from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { tableDataService } from '@/api/services/table/data'
import { useSidebarItemsStore } from '@/stores/ui/sidebar-items'

import { Button } from '@/components/ui/button'
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
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { useTables } from '@/composables/api'
import { useTableDetail } from '@/composables/api/useTableDetail'
import { useToast } from '@/composables/ui'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const sidebarStore = useSidebarItemsStore()

const projectId = computed(() => route.params.projectId as string)
const tableId = computed(() => route.params.tableId as string | undefined)

const { tables, loading: isLoading, createTable } = useTables(projectId.value)

const {
  data: tableDetail,
  loading: tableDetailLoading,
  refresh: refreshTableDetail,
} = useTableDetail(projectId.value, tableId.value ?? '')

watch(
  tables,
  (newTables) => {
    if (projectId.value && sidebarStore.currentProjectId === projectId.value) {
      sidebarStore.updateProjectCounts(newTables.length, sidebarStore.workflowCount)
    }
  },
  { immediate: true }
)

const searchQuery = ref('')
const isCreateDialogOpen = ref(false)
const isCreating = ref(false)
const activeTab = ref('data')

const isInsertDialogOpen = ref(false)
const isEditDialogOpen = ref(false)
const isDeleteDialogOpen = ref(false)
const selectedRow = ref<Record<string, any> | null>(null)
const selectedRowId = ref<string | null>(null)

const currentPage = ref(1)
const pageLimit = ref(50)
const totalRows = ref(0)

const newTable = ref<{
  name: string
  description: string
  columns: Array<{
    name: string
    type: string
    nullable: boolean
    primary_key: boolean
    default_value: string
  }>
}>({
  name: '',
  description: '',
  columns: [],
})

const columnTypes = [
  'VARCHAR',
  'INTEGER',
  'BIGINT',
  'DECIMAL',
  'BOOLEAN',
  'DATE',
  'TIMESTAMP',
  'TEXT',
  'JSON',
  'UUID',
]

const filteredTables = computed(() => {
  if (!tables.value || !Array.isArray(tables.value)) return []
  if (!searchQuery.value) return tables.value
  return tables.value.filter(
    (table) =>
      table.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (table.description &&
        table.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )
})

const selectedTable = computed(() => {
  if (!tableId.value || !tables.value) return null
  return tables.value.find((t: any) => t.id === tableId.value)
})

const tableColumns = computed(() => {
  if (!tableDetail.value) return []
  return tableDetail.value.columns
})

const tableRows = computed(() => {
  if (!tableDetail.value) return []
  return tableDetail.value.rows ?? []
})

const handleCreateTable = async () => {
  if (!newTable.value.name.trim()) {
    toast.error('Validation Error', 'Table name is required')
    return
  }

  if (newTable.value.columns.length === 0) {
    toast.error('Validation Error', 'At least one column is required')
    return
  }

  const hasInvalidColumns = newTable.value.columns.some((col) => !col.name.trim() || !col.type)
  if (hasInvalidColumns) {
    toast.error('Validation Error', 'All columns must have a name and type')
    return
  }

  isCreating.value = true
  try {
    const schema = {
      columns: newTable.value.columns.map((col) => ({
        name: col.name,
        type: col.type,
        nullable: col.nullable,
        default: col.default_value || undefined,
      })),
      primary_key: newTable.value.columns.filter((col) => col.primary_key).map((col) => col.name),
      indexes: [],
      foreign_keys: [],
    }

    const createdTable = await createTable({
      name: newTable.value.name,
      description: newTable.value.description || undefined,
      schema,
    })

    toast.success('Success', `Table "${createdTable.name}" has been created`)
    isCreateDialogOpen.value = false
    newTable.value = { name: '', description: '', columns: [] }

    router.push(`/projects/${projectId.value}/tables/${createdTable.id}`)
  } catch (error: any) {
    console.error('Failed to create table:', error)
    toast.error('Error', error.message || 'Failed to create table')
  } finally {
    isCreating.value = false
  }
}

const handleSelectTable = (id: string) => {
  router.push(`/projects/${projectId.value}/tables/${id}`)
}

const addColumn = () => {
  newTable.value.columns.push({
    name: '',
    type: 'VARCHAR',
    nullable: true,
    primary_key: false,
    default_value: '',
  })
}

const removeColumn = (index: number) => {
  newTable.value.columns.splice(index, 1)
}

const handleInsertRow = async (data: Record<string, any>) => {
  await tableDataService.insert(projectId.value, tableId.value!, data)
  await refreshTableDetail()
}

const handleEditRow = (row: Record<string, any>) => {
  const idColumn = tableColumns.value.find((col) => col.is_primary_key)
  if (idColumn) {
    selectedRowId.value = row[idColumn.name]
    selectedRow.value = row
    isEditDialogOpen.value = true
  } else {
    toast.error('Error', 'No primary key found for this table')
  }
}

const handleUpdateRow = async (rowId: string, data: Record<string, any>) => {
  await tableDataService.update(projectId.value, tableId.value!, rowId, data)
  await refreshTableDetail()
}

const handleDeleteRow = (row: Record<string, any>) => {
  const idColumn = tableColumns.value.find((col) => col.is_primary_key)
  if (idColumn) {
    selectedRowId.value = row[idColumn.name]
    selectedRow.value = row
    isDeleteDialogOpen.value = true
  } else {
    toast.error('Error', 'No primary key found for this table')
  }
}

const handleDeleteConfirm = async (rowId: string) => {
  await tableDataService.delete(projectId.value, tableId.value!, rowId)
  await refreshTableDetail()
}

const totalPages = computed(() => Math.ceil(totalRows.value / pageLimit.value))

const canGoToPreviousPage = computed(() => currentPage.value > 1)
const canGoToNextPage = computed(() => currentPage.value < totalPages.value)

const goToPreviousPage = () => {
  if (canGoToPreviousPage.value) {
    currentPage.value--
    refreshTableDetail()
  }
}

const goToNextPage = () => {
  if (canGoToNextPage.value) {
    currentPage.value++
    refreshTableDetail()
  }
}

const formatCellValue = (value: any, columnType: string): string => {
  if (value === null || value === undefined) return '-'

  const upperType = columnType.toUpperCase()

  if (upperType === 'JSON' || upperType === 'JSONB') {
    if (typeof value === 'object') {
      const str = JSON.stringify(value)
      return str.length > 50 ? str.substring(0, 50) + '...' : str
    }
  }

  if (upperType === 'BOOLEAN' || upperType === 'BOOL') {
    return value ? 'true' : 'false'
  }

  if (upperType.includes('TIMESTAMP') || upperType === 'DATE') {
    try {
      const date = new Date(value)
      return date.toLocaleString()
    } catch {
      return String(value)
    }
  }

  const str = String(value)
  return str.length > 100 ? str.substring(0, 100) + '...' : str
}

watch(tableId, () => {
  if (tableId.value) {
    currentPage.value = 1
    refreshTableDetail()
  }
})
</script>

<template>
  <div class="flex h-[calc(100vh-3.5rem)]">
    <div class="w-80 border-r flex flex-col bg-muted/10">
      <div class="p-4 border-b space-y-3">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Tables</h2>
          <Button size="sm" @click="isCreateDialogOpen = true">
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

      <LoadingWrapper :is-loading="isLoading" loading-text="Loading tables...">
        <div class="flex-1 overflow-y-auto">
          <div class="p-2 space-y-1">
            <button
              v-for="table in filteredTables"
              :key="table.id"
              @click="handleSelectTable(table.id)"
              :class="[
                'w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors',
                tableId === table.id
                  ? 'bg-primary/10 text-primary font-medium'
                  : 'hover:bg-muted text-muted-foreground hover:text-foreground',
              ]"
            >
              <TableIcon class="w-4 h-4 flex-shrink-0" />
              <span class="flex-1 text-left truncate">{{ table.name }}</span>
            </button>
          </div>

          <div
            v-if="!isLoading && filteredTables.length === 0"
            class="p-4 text-center text-sm text-muted-foreground"
          >
            <p>No tables found</p>
          </div>
        </div>
      </LoadingWrapper>
    </div>

    <div class="flex-1 flex flex-col">
      <div v-if="!tableId" class="flex-1 flex items-center justify-center">
        <div class="text-center space-y-4">
          <div class="flex justify-center">
            <div class="w-16 h-16 rounded-full bg-muted flex items-center justify-center">
              <Database class="w-8 h-8 text-muted-foreground" />
            </div>
          </div>
          <div>
            <h3 class="text-lg font-semibold mb-1">Select a table</h3>
            <p class="text-sm text-muted-foreground">
              Choose a table from the list to view its data and schema
            </p>
          </div>
          <Button @click="isCreateDialogOpen = true">
            <Plus class="w-4 h-4 mr-2" />
            Create New Table
          </Button>
        </div>
      </div>

      <div v-else class="flex-1 flex flex-col">
        <div class="border-b px-6 py-4">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold">{{ selectedTable?.name }}</h1>
              <p class="text-sm text-muted-foreground mt-1">
                {{ selectedTable?.description || 'No description' }}
              </p>
            </div>
            <div class="flex items-center gap-2">
              <Button variant="outline" size="sm">
                <Settings2 class="w-4 h-4 mr-2" />
                Settings
              </Button>
            </div>
          </div>
        </div>

        <Tabs v-model="activeTab" class="flex-1 flex flex-col">
          <div class="border-b px-6">
            <TabsList class="h-11 bg-transparent p-0">
              <TabsTrigger value="data" class="h-11">Data</TabsTrigger>
              <TabsTrigger value="schema" class="h-11">Schema</TabsTrigger>
            </TabsList>
          </div>

          <div class="flex-1 overflow-auto">
            <TabsContent value="data" class="m-0 p-6 h-full">
              <LoadingWrapper :is-loading="tableDetailLoading" loading-text="Loading table data...">
                <div class="space-y-4">
                  <div class="flex items-center justify-between">
                    <div class="text-sm text-muted-foreground">{{ tableRows.length }} rows</div>
                    <div class="flex items-center gap-2">
                      <Button variant="outline" size="sm" @click="isInsertDialogOpen = true">
                        <Plus class="w-4 h-4 mr-2" />
                        Insert Row
                      </Button>
                    </div>
                  </div>

                  <div class="border rounded-lg overflow-hidden">
                    <div class="overflow-x-auto">
                      <table class="w-full">
                        <thead class="bg-muted/50">
                          <tr>
                            <th
                              v-for="column in tableColumns"
                              :key="column.name"
                              class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                            >
                              {{ column.name }}
                            </th>
                            <th
                              class="px-4 py-3 text-right text-xs font-medium text-muted-foreground uppercase tracking-wider w-24"
                            >
                              Actions
                            </th>
                          </tr>
                        </thead>
                        <tbody class="divide-y divide-border bg-background">
                          <tr
                            v-for="(row, idx) in tableRows"
                            :key="idx"
                            class="hover:bg-muted/30 transition-colors group"
                          >
                            <td
                              v-for="column in tableColumns"
                              :key="column.name"
                              class="px-4 py-3 text-sm text-foreground"
                            >
                              <div
                                class="max-w-xs truncate"
                                :title="String(row[column.name] ?? '-')"
                              >
                                {{ formatCellValue(row[column.name], column.type) }}
                              </div>
                            </td>
                            <td class="px-4 py-3 text-right">
                              <div
                                class="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
                              >
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  @click="handleEditRow(row)"
                                  class="h-8 w-8 p-0"
                                >
                                  <Pencil class="w-4 h-4" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  @click="handleDeleteRow(row)"
                                  class="h-8 w-8 p-0 text-destructive hover:text-destructive"
                                >
                                  <Trash2 class="w-4 h-4" />
                                </Button>
                              </div>
                            </td>
                          </tr>
                          <tr v-if="tableRows.length === 0">
                            <td
                              :colspan="tableColumns.length + 1"
                              class="px-4 py-8 text-center text-sm text-muted-foreground"
                            >
                              No data yet
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>

                  <div v-if="tableRows.length > 0" class="flex items-center justify-between pt-4">
                    <div class="text-sm text-muted-foreground">
                      Page {{ currentPage }} of {{ totalPages || 1 }}
                    </div>
                    <div class="flex items-center gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        @click="goToPreviousPage"
                        :disabled="!canGoToPreviousPage"
                      >
                        <ChevronLeft class="w-4 h-4" />
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        @click="goToNextPage"
                        :disabled="!canGoToNextPage"
                      >
                        <ChevronRight class="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              </LoadingWrapper>
            </TabsContent>

            <TabsContent value="schema" class="m-0 p-6 h-full">
              <LoadingWrapper :is-loading="tableDetailLoading" loading-text="Loading schema...">
                <div class="space-y-4">
                  <div class="flex items-center justify-between">
                    <div class="text-sm text-muted-foreground">
                      {{ tableColumns.length }} columns
                    </div>
                    <div class="flex items-center gap-2">
                      <Button variant="outline" size="sm">
                        <Plus class="w-4 h-4 mr-2" />
                        Add Column
                      </Button>
                    </div>
                  </div>

                  <div class="border rounded-lg overflow-hidden">
                    <div class="overflow-x-auto">
                      <table class="w-full">
                        <thead class="bg-muted/50">
                          <tr>
                            <th
                              class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                            >
                              Name
                            </th>
                            <th
                              class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                            >
                              Type
                            </th>
                            <th
                              class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                            >
                              Nullable
                            </th>
                            <th
                              class="px-4 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider"
                            >
                              Default
                            </th>
                          </tr>
                        </thead>
                        <tbody class="divide-y divide-border bg-background">
                          <tr
                            v-for="column in tableColumns"
                            :key="column.name"
                            class="hover:bg-muted/30 transition-colors"
                          >
                            <td class="px-4 py-3 text-sm font-medium text-foreground">
                              {{ column.name }}
                            </td>
                            <td class="px-4 py-3 text-sm text-muted-foreground">
                              {{ column.type }}
                            </td>
                            <td class="px-4 py-3 text-sm text-muted-foreground">
                              {{ column.is_nullable ? 'Yes' : 'No' }}
                            </td>
                            <td class="px-4 py-3 text-sm text-muted-foreground">
                              {{ column.default_value ?? '-' }}
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>
                </div>
              </LoadingWrapper>
            </TabsContent>
          </div>
        </Tabs>
      </div>
    </div>

    <Dialog v-model:open="isCreateDialogOpen">
      <DialogContent class="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Create New Table</DialogTitle>
          <DialogDescription>
            Create a new table with custom columns and schema.
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="table-name">Table Name</Label>
              <Input
                id="table-name"
                v-model="newTable.name"
                placeholder="e.g., users, products"
                required
              />
            </div>
            <div class="space-y-2">
              <Label for="table-description">Description</Label>
              <Input
                id="table-description"
                v-model="newTable.description"
                placeholder="Brief description"
              />
            </div>
          </div>

          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <Label>Columns</Label>
                <p class="text-xs text-muted-foreground mt-1">At least one column is required</p>
              </div>
              <Button type="button" variant="outline" size="sm" @click="addColumn">
                <Plus class="w-4 h-4 mr-1" />
                Add Column
              </Button>
            </div>

            <div class="space-y-3 max-h-60 overflow-y-auto">
              <div
                v-for="(column, index) in newTable.columns"
                :key="index"
                class="flex items-center space-x-2 p-3 border rounded-lg"
              >
                <div class="flex-1 grid grid-cols-4 gap-2">
                  <Input v-model="column.name" placeholder="Column name" />
                  <Select v-model="column.type">
                    <SelectTrigger>
                      <SelectValue placeholder="Type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="type in columnTypes" :key="type" :value="type">
                        {{ type }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <Input v-model="column.default_value" placeholder="Default value" />
                  <div class="flex items-center space-x-2">
                    <input v-model="column.nullable" type="checkbox" class="rounded" />
                    <span class="text-sm">Nullable</span>
                  </div>
                </div>
                <Button type="button" variant="ghost" size="sm" @click="removeColumn(index)">
                  <Plus class="w-4 h-4 rotate-45" />
                </Button>
              </div>
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="isCreateDialogOpen = false" :disabled="isCreating">
            Cancel
          </Button>
          <Button
            @click="handleCreateTable"
            :disabled="!newTable.name.trim() || newTable.columns.length === 0 || isCreating"
          >
            {{ isCreating ? 'Creating...' : 'Create Table' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <InsertRowDialog
      :open="isInsertDialogOpen"
      @update:open="isInsertDialogOpen = $event"
      :columns="tableColumns"
      :on-insert="handleInsertRow"
    />

    <EditRowDialog
      :open="isEditDialogOpen"
      @update:open="isEditDialogOpen = $event"
      :columns="tableColumns"
      :row-data="selectedRow"
      :row-id="selectedRowId"
      :on-update="handleUpdateRow"
    />

    <DeleteRowDialog
      :open="isDeleteDialogOpen"
      @update:open="isDeleteDialogOpen = $event"
      :row-data="selectedRow"
      :row-id="selectedRowId"
      :on-delete="handleDeleteConfirm"
    />
  </div>
</template>
