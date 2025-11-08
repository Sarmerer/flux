import { computed, ref, watch } from 'vue'

export interface PaginationState {
  currentPage: number
  pageSize: number
  totalItems: number
}

export function useTablePagination(initialPageSize = 50) {
  const currentPage = ref(1)
  const pageSize = ref(initialPageSize)
  const totalItems = ref(0)

  const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value))

  const canGoToPreviousPage = computed(() => currentPage.value > 1)
  const canGoToNextPage = computed(() => currentPage.value < totalPages.value)

  const hasMultiplePages = computed(() => totalPages.value > 1)

  const startIndex = computed(() => (currentPage.value - 1) * pageSize.value + 1)
  const endIndex = computed(() => Math.min(currentPage.value * pageSize.value, totalItems.value))

  const goToPage = (page: number) => {
    if (page < 1 || page > totalPages.value) return
    currentPage.value = page
  }

  const goToPreviousPage = () => {
    if (canGoToPreviousPage.value) {
      currentPage.value--
    }
  }

  const goToNextPage = () => {
    if (canGoToNextPage.value) {
      currentPage.value++
    }
  }

  const setPageSize = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
  }

  const setTotalItems = (total: number) => {
    totalItems.value = total
  }

  const reset = () => {
    currentPage.value = 1
  }

  watch(pageSize, reset)

  return {
    currentPage: computed(() => currentPage.value),
    pageSize: computed(() => pageSize.value),
    totalItems: computed(() => totalItems.value),
    totalPages,
    canGoToPreviousPage,
    canGoToNextPage,
    hasMultiplePages,
    startIndex,
    endIndex,
    goToPage,
    goToPreviousPage,
    goToNextPage,
    setPageSize,
    setTotalItems,
    reset,
  }
}
