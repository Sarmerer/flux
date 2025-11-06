<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'

export interface TableSkeletonProps {
  rows?: number
  columns?: number
  showHeader?: boolean
}

const props = withDefaults(defineProps<TableSkeletonProps>(), {
  rows: 5,
  columns: 4,
  showHeader: true,
})
</script>

<template>
  <div class="w-full space-y-4">
    <!-- Header -->
    <div v-if="showHeader" class="flex items-center space-x-4">
      <Skeleton class="h-10 flex-1" v-for="n in columns" :key="`header-${n}`" />
    </div>

    <!-- Rows -->
    <div class="space-y-3">
      <div v-for="row in rows" :key="`row-${row}`" class="flex items-center space-x-4">
        <Skeleton class="h-12 flex-1" v-for="col in columns" :key="`row-${row}-col-${col}`" />
      </div>
    </div>
  </div>
</template>
