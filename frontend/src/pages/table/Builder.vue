<script setup lang="ts">
import LoadingWrapper from '@/components/common/LoadingWrapper.vue'
import TableDataView from '@/components/tables/TableDataView.vue'
import type { TableFormData } from '@/components/tables/drawers/TableEditorContent.vue'
import { Database, Pencil, Plus, Search, Table as TableIcon, Trash2 } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

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

import { useTableData, useTables } from '@/composables/api'
import { useTableDetail } from '@/composables/api/useTableDetail'
import { useDrawers } from '@/composables/drawerRegistry'
import { useRouteContext } from '@/composables/routing'
import { useToast } from '@/composables/ui'

const router = useRouter()
const toast = useToast()
const sidebarStore = useSidebarItemsStore()
const {
  openTableEditor: openTableEditorDrawer,
  openInsertRow,
  openEditRow,
  openDeleteRow,
} = useDrawers()

const { projectId, tableId } = useRouteContext()

const { tables, loading: isLoading, createTable, deleteTable } = useTables(projectId.value)

const {
  data: tableDetail,
  columns: tableColumns,
  rows: tableRows,
  totalRows,
  loading: tableDetailLoading,
  currentPage,
  pageLimit,
  refresh: refreshTableDetail,
  goToPage,
} = useTableDetail(projectId.value, tableId.value ?? '')

const { insertRow, updateRow, deleteRow } = useTableData(projectId.value, tableId.value ?? '')

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
const isDeleteTableDialogOpen = ref(false)
const tableToDelete = ref<string | null>(null)

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
          primary_key: data.columns.filter((col) => col.primary_key).map((col) => col.name),
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

const openEditTableSchema = () => {
  if (!selectedTable.value || !tableDetail.value) return

  const initialData: TableFormData = {
    name: selectedTable.value.name,
    description: selectedTable.value.description || '',
    columns: tableColumns.value.map((col) => ({
      name: col.name,
      type: col.type,
      nullable: col.is_nullable,
      primary_key: col.is_primary_key || false,
      default_value: col.default_value || '',
      is_identity: col.is_identity || false,
      unique: col.unique || false,
    })),
    foreignKeys: [],
  }

  openTableEditorDrawer({
    props: {
      mode: 'edit',
      tableId: selectedTable.value.id,
      initialData,
    },
    width: 'w-[600px] sm:max-w-[600px]',
    onSave: async () => {
      toast.info('Info', 'Schema editing is not yet implemented')
    },
  })
}

const handleSelectTable = (id: string) => {
  router.push(`/projects/${projectId.value}/tables/${id}`)
}

const openInsertRowDrawer = () => {
  openInsertRow({
    props: {
      columns: tableColumns.value,
      onInsert: async (data: Record<string, any>) => {
        await insertRow(data)
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
        await updateRow(id, data)
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
        await deleteRow(id)
        await refreshTableDetail()
      },
    },
  })
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
              Choose a table from the list to view and manage its data
            </p>
          </div>
          <Button @click="openTableEditor">
            <Plus class="w-4 h-4 mr-2" />
            Create New Table
          </Button>
        </div>
      </div>

      <div v-else class="flex-1 flex flex-col overflow-hidden">
        <div class="border-b px-6 py-3 bg-background">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-xl font-semibold">{{ selectedTable?.name }}</h1>
              <p v-if="selectedTable?.description" class="text-xs text-muted-foreground mt-0.5">
                {{ selectedTable.description }}
              </p>
            </div>
            <div class="flex items-center gap-2">
              <Button variant="outline" size="sm" @click="openEditTableSchema" class="h-8">
                <Pencil class="w-3.5 h-3.5 mr-2" />
                Edit Schema
              </Button>
              <Button
                variant="outline"
                size="sm"
                @click="handleDeleteTable(selectedTable!.id)"
                class="text-destructive hover:text-destructive hover:bg-destructive/10 h-8"
              >
                <Trash2 class="w-3.5 h-3.5 mr-2" />
                Delete
              </Button>
            </div>
          </div>
        </div>

        <div class="flex-1 flex flex-col overflow-hidden">
          <div class="px-6 py-3 border-b bg-background flex items-center justify-between">
            <div class="text-sm text-muted-foreground">
              {{ totalRows }} {{ totalRows === 1 ? 'row' : 'rows' }}
            </div>
            <Button variant="default" size="sm" @click="openInsertRowDrawer" class="h-8">
              <Plus class="w-4 h-4 mr-2" />
              Insert Row
            </Button>
          </div>

          <div class="flex-1 overflow-auto px-6 py-4">
            <LoadingWrapper :is-loading="tableDetailLoading" loading-text="Loading table data...">
              <TableDataView
                :columns="tableColumns"
                :rows="tableRows"
                :total-rows="totalRows"
                :current-page="currentPage"
                :page-size="pageLimit"
                :on-page-change="goToPage"
                @page-change="goToPage"
                @edit-row="openEditRowDrawer"
                @delete-row="openDeleteRowDrawer"
              />
            </LoadingWrapper>
          </div>
        </div>
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
