import { computed, ref } from 'vue'

import { tableService } from '@/api/services/table'
import type { Table } from '@/types/api'
import { defineStore } from 'pinia'

export const useTableStore = defineStore('tables', () => {
  const tables = ref<Table[]>([])
  const isLoading = ref(false)
  const error = ref<Error | null>(null)
  const length = computed(() => tables.value.length)

  const loadByProjectId = async (projectId: string) => {
    isLoading.value = true
    error.value = null
    try {
      tables.value = await tableService.getAll(projectId)
      return tables
    } catch (_error) {
      console.error('Failed to load tables:', _error)
      error.value = new Error('Failed to load tables')
      tables.value = []
      throw _error
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
    error,
    loadByProjectId,
    clear,
  }
})
