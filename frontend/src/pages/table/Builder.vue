<script setup lang="ts">
import LoadingWrapper from '@/components/common/LoadingWrapper.vue'
import type { TableFormData } from '@/components/tables/drawers/TableEditorContent.vue'
import {
  ChevronLeft,
  ChevronRight,
  Database,
  Key,
  Pencil,
  Plus,
  Search,
  Settings2,
  Sparkles,
  Table as TableIcon,
  Trash2,
} from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { tableDataService } from '@/api/services/table/data'
import { useSidebarItemsStore } from '@/stores/ui/sidebar-items'

import { Badge } from '@/components/ui/badge'
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
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { useTables } from '@/composables/api'
import { useTableDetail } from '@/composables/api/useTableDetail'
import { useDrawers } from '@/composables/drawerRegistry'
import { useRouteContext } from '@/composables/routing'
import { useToast } from '@/composables/ui'

const router = useRouter()
const toast = useToast()
const sidebarStore = useSidebarItemsStore()
const { openTableEditor: openTableEditorDrawer, openInsertRow, openEditRow, openDeleteRow } = useDrawers()

const { projectId, tableId } = useRouteContext()

const { tables, loading: isLoading, createTable, deleteTable } = useTables(projectId.value)

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
const activeTab = ref('data')

const isDeleteTableDialogOpen = ref(false)
const tableToDelete = ref<string | null>(null)

const currentPage = ref(1)
const pageLimit = ref(50)
const totalRows = ref(0)

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

const openTableEditor = () => {
  openTableEditorDrawer({
    props: { mode: 'create' },
    width: 'w-[600px] sm:max-w-[600px]',
    onSave: async (data: TableFormData) => {
      try {
        const schema = {
          columns: data.columns.map((col) => ({
            name: col.name,
            type: col.type,
            nullable: col.nullable,
            default_value: col.default_value || undefined,
          })),
          primary_key: data.columns
            .filter((col) => col.primary_key)
            .map((col) => col.name),
          indexes: [],
          foreign_keys: data.foreignKeys.map((fk) => ({
            column: fk.column,
            referenced_table: fk.referencedTable,
            referenced_column: fk.referencedColumn,
            on_update: fk.onUpdate,
            on_delete: fk.onDelete,
          })),
        }

        const createdTable = await createTable({
          name: data.name,
          description: data.description || undefined,
          schema,
        })

        toast.success('Success', `Table "${createdTable.name}" has been created`)
        router.push(`/projects/${projectId.value}/tables/${createdTable.id}`)
      } catch (error: any) {
        console.error('Failed to create table:', error)
        toast.error('Error', error.message || 'Failed to create table')
      }
    },
  })
}

const handleSelectTable = (id: string) => {
  router.push(`/projects/${projectId.value}/tables/${id}`)
}

const openAddColumnDrawer = () => {
  toast.info('Info', 'Please use the table editor to add columns')
}

const openInsertRowDrawer = () => {
  openInsertRow({
    props: {
      columns: tableColumns.value,
      onInsert: async (data: Record<string, any>) => {
        await tableDataService.insert(projectId.value, tableId.value!, data)
        await refreshTableDetail()
      },
    },
  })
}

const openEditRowDrawer = (row: Record<string, any>) => {
  const idColumn = tableColumns.value.find((col) => col.is_primary_key)
  if (!idColumn) {
    toast.error('Error', 'No primary key found for this table')
    return
  }

  const rowId = row[idColumn.name]
  openEditRow({
    props: {
      columns: tableColumns.value,
      rowData: row,
      rowId,
      onUpdate: async (id: string, data: Record<string, any>) => {
        await tableDataService.update(projectId.value, tableId.value!, id, data)
        await refreshTableDetail()
      },
    },
  })
}

const openDeleteRowDrawer = (row: Record<string, any>) => {
  const idColumn = tableColumns.value.find((col) => col.is_primary_key)
  if (!idColumn) {
    toast.error('Error', 'No primary key found for this table')
    return
  }

  const rowId = row[idColumn.name]
  openDeleteRow({
    props: {
      rowId,
      rowData: row,
      onDelete: async (id: string) => {
        await tableDataService.delete(projectId.value, tableId.value!, id)
        await refreshTableDetail()
      },
    },
  })
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


const handleDeleteTable = (id: string) => {
  tableToDelete.value = id
  isDeleteTableDialogOpen.value = true
}

const confirmDeleteTable = async () => {
  if (!tableToDelete.value) return

  try {
    await deleteTable(tableToDelete.value)
    toast.success('Success', 'Table deleted successfully')
    isDeleteTableDialogOpen.value = false
    tableToDelete.value = null

    if (tableId.value === tableToDelete.value) {
      router.push(`/projects/${projectId.value}/tables`)
    }
  } catch (error: any) {
    console.error('Failed to delete table:', error)
    toast.error('Error', error.message || 'Failed to delete table')
  }
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
          <Button size="sm" @click="openTableEditor">
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
            <div
              v-for="table in filteredTables"
              :key="table.id"
              :class="[
                'group relative w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors',
                tableId === table.id
                  ? 'bg-primary/10 text-primary font-medium'
                  : 'hover:bg-muted text-muted-foreground hover:text-foreground',
              ]"
            >
              <button
                @click="handleSelectTable(table.id)"
                class="flex items-center gap-3 flex-1 min-w-0"
              >
                <TableIcon class="w-4 h-4 flex-shrink-0" />
                <span class="flex-1 text-left truncate">{{ table.name }}</span>
              </button>
              <Button
                variant="ghost"
                size="sm"
                @click.stop="handleDeleteTable(table.id)"
                class="h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity text-destructive hover:text-destructive hover:bg-destructive/10"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </Button>
            </div>
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
          <Button @click="openTableEditor">
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
              <Button
                variant="outline"
                size="sm"
                @click="handleDeleteTable(selectedTable!.id)"
                class="text-destructive hover:text-destructive hover:bg-destructive/10"
              >
                <Trash2 class="w-4 h-4 mr-2" />
                Delete
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
                      <Button variant="outline" size="sm" @click="openInsertRowDrawer">
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
                                  @click="openEditRowDrawer(row)"
                                  class="h-8 w-8 p-0"
                                >
                                  <Pencil class="w-4 h-4" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  @click="openDeleteRowDrawer(row)"
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
                      <Button variant="outline" size="sm" @click="openAddColumnDrawer">
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
                              Constraints
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
                              <div class="flex items-center gap-2">
                                {{ column.name }}
                                <Key
                                  v-if="column.is_primary_key"
                                  class="w-3.5 h-3.5 text-amber-500"
                                  title="Primary Key"
                                />
                              </div>
                            </td>
                            <td class="px-4 py-3 text-sm text-muted-foreground">
                              {{ column.type }}
                            </td>
                            <td class="px-4 py-3 text-sm">
                              <div class="flex flex-wrap gap-1">
                                <Badge
                                  v-if="column.is_identity"
                                  variant="secondary"
                                  class="text-xs"
                                >
                                  <Sparkles class="w-3 h-3 mr-1" />
                                  Auto
                                </Badge>
                                <Badge v-if="!column.is_nullable" variant="outline" class="text-xs">
                                  NOT NULL
                                </Badge>
                                <Badge v-if="column.unique" variant="outline" class="text-xs">
                                  UNIQUE
                                </Badge>
                                <span
                                  v-if="column.is_nullable && !column.unique && !column.is_identity"
                                  class="text-muted-foreground"
                                  >-</span
                                >
                              </div>
                            </td>
                            <td class="px-4 py-3 text-sm text-muted-foreground">
                              <code
                                v-if="column.default_value"
                                class="text-xs bg-muted px-1.5 py-0.5 rounded"
                              >
                                {{ column.default_value }}
                              </code>
                              <span v-else>-</span>
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

    <Dialog v-model:open="isDeleteTableDialogOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Table</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this table? This action cannot be undone and all data
            will be permanently deleted.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="isDeleteTableDialogOpen = false">Cancel</Button>
          <Button variant="destructive" @click="confirmDeleteTable">Delete Table</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
