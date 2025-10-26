<script setup lang="ts">
import LoadingSpinner from './LoadingSpinner.vue'

export interface LoadingWrapperProps {
  isLoading: boolean
  loadingText?: string
}

withDefaults(defineProps<LoadingWrapperProps>(), {
  loadingText: 'Loading...',
})
</script>

<template>
  <LoadingSpinner v-if="isLoading" :text="loadingText" />
  <slot v-else />
</template>
