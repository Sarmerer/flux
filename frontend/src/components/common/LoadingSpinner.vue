<script setup lang="ts">
export interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  text?: string
}

withDefaults(defineProps<LoadingSpinnerProps>(), {
  size: 'md',
})

const sizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-8 h-8',
  lg: 'w-12 h-12',
  xl: 'w-16 h-16',
}
</script>

<template>
  <div class="flex flex-col items-center justify-center p-8 space-y-4">
    <div
      :class="[sizeClasses[size], 'animate-spin rounded-full border-4 border-primary border-t-transparent']"
    ></div>
    <p v-if="text" class="text-sm text-muted-foreground">{{ text }}</p>
  </div>
</template>
