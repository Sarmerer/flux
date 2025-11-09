<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import LoadingWrapper from '@/components/common/LoadingWrapper.vue'
import TableDataView from '@/components/tables/TableDataView.vue'
import TableSidebar from '@/components/tables/builder/TableSidebar.vue'
import TableHeader from '@/components/tables/builder/TableHeader.vue'
import EmptyTableState from '@/components/tables/builder/EmptyTableState.vue'

import type { TableFormData } from '@/components/tables/drawers/TableEditorContent.vue'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

import { useSidebarItemsStore } from '@/stores/ui/sidebar-items'
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

const tableDetailData = ref<any>(null)
const tableColumns = ref<any[]>([])
const tableRows = ref<any[]>([])
const totalRows = ref(0)
const tableDetailLoading = ref(false)
const currentPage = ref(1)
const pageLimit = ref(50)
const selectedRows = ref<any[]>([])
const isDeleteTableDialogOpen = ref(false)
const tableToDelete = ref<string | null>(null)
const tableDataViewRef = ref<InstanceType<typeof TableDataView> | null>(null)

const tableInstance = computed(() => tableDataViewRef.value?.dataTableRef?.table)
const allColumns = computed(() => tableInstance.value?.getAllColumns() || [])

let currentTableDetail: ReturnType<typeof useTableDetail> | null = null

watch(
  tableId,
  (newTableId) => {
    currentPage.value = 1
    selectedRows.value = []

    if (newTableId) {
      currentTableDetail = useTableDetail(projectId.value, newTableId, {
        page: currentPage.value,
        limit: pageLimit.value,
        autoFetch: true,
      })

      watch(
        currentTableDetail.data,
        (newData) => {
          tableDetailData.value = newData
          tableColumns.value = currentTableDetail!.columns.value
          tableRows.value = currentTableDetail!.rows.value
          totalRows.value = currentTableDetail!.totalRows.value
        },
        { immediate: true }
      )

      watch(
        currentTableDetail.loading,
        (loading) => {
          tableDetailLoading.value = loading
        },
        { immediate: true }
      )

      watch(
        currentTableDetail.currentPage,
        (page) => {
          currentPage.value = page
        },
        { immediate: true }
      )

      watch(
        currentTableDetail.pageLimit,
        (limit) => {
          pageLimit.value = limit
        },
        { immediate: true }
      )
    } else {
      tableDetailData.value = null
      tableColumns.value = []
      tableRows.value = []
      totalRows.value = 0
    }
  },
  { immediate: true }
)

const refreshTableDetail = () => {
  if (currentTableDetail) {
    currentTableDetail.refresh()
  }
}

const goToPage = (page: number) => {
  if (currentTableDetail) {
    currentTableDetail.goToPage(page)
  }
}

const setPageLimit = (limit: number) => {
  if (currentTableDetail) {
    currentTableDetail.setPageLimit(limit)
  }
}

watch(
  tables,
  (newTables) => {
    if (projectId.value && sidebarStore.currentProjectId === projectId.value) {
      sidebarStore.updateProjectCounts(newTables.length, sidebarStore.workflowCount)
    }
  },
  { immediate: true }
)

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
            is_identity: col.is_identity,
            unique: col.unique,
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
  if (!selectedTable.value || !tableDetailData.value) return

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
  if (!tableId.value) return

  const tableDataService = useTableData(projectId.value, tableId.value)

  openInsertRow({
    props: {
      columns: tableColumns.value,
      onInsert: async (data: Record<string, any>) => {
        await tableDataService.insertRow(data)
        refreshTableDetail()
      },
    },
  })
}

const openEditRowDrawer = (row: Record<string, any>) => {
  if (!tableId.value) return

  const idColumn = tableColumns.value.find((col) => col.is_primary_key)
  if (!idColumn) {
    toast.error('Error', 'No primary key found for this table')
    return
  }

  const rowId = row[idColumn.name]
  const tableDataService = useTableData(projectId.value, tableId.value)

  openEditRow({
    props: {
      columns: tableColumns.value,
      rowData: row,
      rowId,
      onUpdate: async (id: string, data: Record<string, any>) => {
        await tableDataService.updateRow(id, data)
        refreshTableDetail()
      },
    },
  })
}

const openDeleteRowDrawer = (row: Record<string, any>) => {
  if (!tableId.value) return

  const idColumn = tableColumns.value.find((col) => col.is_primary_key)
  if (!idColumn) {
    toast.error('Error', 'No primary key found for this table')
    return
  }

  const rowId = row[idColumn.name]
  const tableDataService = useTableData(projectId.value, tableId.value)

  openDeleteRow({
    props: {
      rowId,
      rowData: row,
      onDelete: async (id: string) => {
        await tableDataService.deleteRow(id)
        refreshTableDetail()
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

const handleRowSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

const handleDeleteSelected = async () => {
  if (selectedRows.value.length === 0 || !tableId.value) return

  const idColumn = tableColumns.value.find((col) => col.is_primary_key)
  if (!idColumn) {
    toast.error('Error', 'No primary key found for this table')
    return
  }

  const tableDataService = useTableData(projectId.value, tableId.value)
  const count = selectedRows.value.length

  try {
    for (const row of selectedRows.value) {
      const rowId = row[idColumn.name]
      await tableDataService.deleteRow(rowId)
    }
    selectedRows.value = []
    refreshTableDetail()
    toast.success('Success', `Deleted ${count} row${count > 1 ? 's' : ''}`)
  } catch (error: any) {
    console.error('Failed to delete rows:', error)
    toast.error('Error', error.message || 'Failed to delete rows')
  }
}

const handlePageSizeChange = (size: number) => {
  setPageLimit(size)
}

const handleToggleColumn = (columnId: string, visible: boolean) => {
  const column = tableInstance.value?.getColumn(columnId)
  if (column) {
    column.toggleVisibility(visible)
  }
}
</script>

<template>
  <div class="flex h-[calc(100vh-3.5rem)]">
    <TableSidebar
      :tables="tables"
      :loading="isLoading"
      :selected-table-id="tableId"
      @create="openTableEditor"
      @select="handleSelectTable"
      @delete="handleDeleteTable"
    />

    <div class="flex-1 flex flex-col">
      <EmptyTableState v-if="!tableId" @create="openTableEditor" />

      <div v-else class="flex-1 flex flex-col overflow-hidden">
        <TableHeader
          v-if="selectedTable"
          :table-name="selectedTable.name"
          :table-description="selectedTable.description"
          :loading="tableDetailLoading"
          :total-rows="totalRows"
          :selected-count="selectedRows.length"
          :page-size="pageLimit"
          :all-columns="allColumns"
          @edit-schema="openEditTableSchema"
          @refresh="refreshTableDetail"
          @delete="handleDeleteTable(selectedTable.id)"
          @insert-row="openInsertRowDrawer"
          @delete-selected="handleDeleteSelected"
          @page-size-change="handlePageSizeChange"
          @toggle-column="handleToggleColumn"
        />

        <div class="flex-1 flex flex-col overflow-hidden">
          <LoadingWrapper :is-loading="tableDetailLoading" loading-text="Loading table data...">
            <TableDataView
              ref="tableDataViewRef"
              :columns="tableColumns"
              :rows="tableRows"
              :total-rows="totalRows"
              :current-page="currentPage"
              :page-size="pageLimit"
              :enable-row-selection="true"
              :on-page-change="goToPage"
              @page-change="goToPage"
              @edit-row="openEditRowDrawer"
              @delete-row="openDeleteRowDrawer"
              @row-selection-change="handleRowSelectionChange"
            />
          </LoadingWrapper>
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
