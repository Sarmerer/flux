import { computed, ref } from 'vue'

import { tableService } from '@/api/services/table'
import type { Table } from '@/types/api'
import { defineStore } from 'pinia'

export const useTableStore = defineStore('tables', () => {
  const tablesByProject = ref<Map<string, Table[]>>(new Map())
  const isLoading = ref(false)
  const currentProjectId = ref<string | null>(null)

  const currentTables = computed(() => {
    if (!currentProjectId.value) return []
    return tablesByProject.value.get(currentProjectId.value) || []
  })

  const tableCount = computed(() => currentTables.value.length)

  const setCurrentProject = (projectId: string | null) => {
    currentProjectId.value = projectId
  }

  const loadTables = async (projectId: string) => {
    isLoading.value = true
    try {
      const tables = await tableService.getAll(projectId)
      tablesByProject.value.set(projectId, Array.isArray(tables) ? tables : [])
      currentProjectId.value = projectId
      return tables
    } catch (error) {
      console.error('Failed to load tables:', error)
      tablesByProject.value.set(projectId, [])
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const addTable = (projectId: string, table: Table) => {
    const tables = tablesByProject.value.get(projectId) || []
    tablesByProject.value.set(projectId, [...tables, table])
  }

  const updateTable = (projectId: string, tableId: string, updates: Partial<Table>) => {
    const tables = tablesByProject.value.get(projectId)
    if (!tables) return

    const index = tables.findIndex((t) => t.id === tableId)
    if (index !== -1) {
      const updatedTable = { ...tables[index], ...updates } as Table
      tables[index] = updatedTable
      tablesByProject.value.set(projectId, [...tables])
    }
  }

  const removeTable = (projectId: string, tableId: string) => {
    const tables = tablesByProject.value.get(projectId)
    if (!tables) return

    tablesByProject.value.set(
      projectId,
      tables.filter((t) => t.id !== tableId)
    )
  }

  const clearProject = (projectId: string) => {
    tablesByProject.value.delete(projectId)
  }

  const clearAll = () => {
    tablesByProject.value.clear()
    currentProjectId.value = null
  }

  return {
    currentTables,
    tableCount,
    isLoading,
    currentProjectId,
    setCurrentProject,
    loadTables,
    addTable,
    updateTable,
    removeTable,
    clearProject,
    clearAll,
  }
})
