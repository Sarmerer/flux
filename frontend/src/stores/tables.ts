import { computed, ref } from 'vue'

import { tableService } from '@/api/services/table'
import type { Table } from '@/types/api'
import { defineStore } from 'pinia'

export const useTableStore = defineStore('tables', () => {
  const tables = ref<Table[]>([])
  const isLoading = ref(false)
  const length = computed(() => tables.value.length)

  const loadByProjectId = async (projectId: string) => {
    isLoading.value = true
    try {
      tables.value = await tableService.getAll(projectId)
      return tables
    } catch (error) {
      console.error('Failed to load tables:', error)
      tables.value = []
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const clear = () => {
    tables.value = []
  }

  return {
    tables,
    length,
    isLoading,
    loadByProjectId,
    clear,
  }
})
