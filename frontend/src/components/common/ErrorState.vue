<script setup lang="ts">
import { AlertCircle, RefreshCw } from 'lucide-vue-next'
import { useSlots } from 'vue'

import { Button } from '@/components/ui/button'

export interface ErrorStateProps {
  error: Error | string | null
  title?: string
  description?: string
  handleRetry?: () => void
}

const props = withDefaults(defineProps<ErrorStateProps>(), {
  title: 'Failed to load',
})

const slots = useSlots()

const errorMessage =
  typeof props.error === 'string' ? props.error : props.error?.message || 'An error occurred'
</script>

<template>
  <div class="flex items-center justify-center py-12">
    <div class="text-center max-w-md px-4">
      <div
        class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-red-100 dark:bg-red-900/30 mb-4"
      >
        <AlertCircle class="w-6 h-6 text-red-600 dark:text-red-400" />
      </div>
      <h3 class="text-sm font-medium text-gray-900 dark:text-gray-100 mb-2">{{ title }}</h3>
      <p class="text-sm text-gray-500 dark:text-gray-400 mb-6">
        {{ description || errorMessage }}
      </p>
      <div class="flex items-center justify-center gap-3">
        <Button v-if="handleRetry" @click="handleRetry" variant="outline" size="sm">
          <RefreshCw class="w-4 h-4 mr-2" />
          Retry
        </Button>
        <slot v-if="slots.default" />
      </div>
    </div>
  </div>
</template>
