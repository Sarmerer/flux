<script setup lang="ts">
import { X, CheckCircle2, AlertCircle, Info, AlertTriangle } from 'lucide-vue-next'
import { computed } from 'vue'
import { Button } from '@/components/ui/button'

export interface ToastProps {
  id?: string
  title?: string
  description?: string
  variant?: 'default' | 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

const props = withDefaults(defineProps<ToastProps>(), {
  variant: 'default',
  duration: 5000,
})

const emit = defineEmits<{
  close: []
}>()

const icon = computed(() => {
  switch (props.variant) {
    case 'success':
      return CheckCircle2
    case 'error':
      return AlertCircle
    case 'warning':
      return AlertTriangle
    case 'info':
      return Info
    default:
      return null
  }
})

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'success':
      return 'bg-success/10 border-success text-success-foreground'
    case 'error':
      return 'bg-destructive/10 border-destructive text-destructive-foreground'
    case 'warning':
      return 'bg-warning/10 border-warning text-warning-foreground'
    case 'info':
      return 'bg-info/10 border-info text-info-foreground'
    default:
      return 'bg-card border-border text-card-foreground'
  }
})
</script>

<template>
  <div
    :class="[
      'group pointer-events-auto relative flex w-full items-center justify-between space-x-4 overflow-hidden rounded-md border p-4 pr-8 shadow-lg transition-all',
      variantClasses,
    ]"
  >
    <div class="flex items-start space-x-3 flex-1">
      <component :is="icon" v-if="icon" class="w-5 h-5 mt-0.5 flex-shrink-0" />
      <div class="flex-1 space-y-1">
        <div v-if="title" class="text-sm font-semibold">
          {{ title }}
        </div>
        <div v-if="description" class="text-sm opacity-90">
          {{ description }}
        </div>
      </div>
    </div>
    <Button
      variant="ghost"
      size="sm"
      class="absolute right-2 top-2 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
      @click="emit('close')"
    >
      <X class="h-4 w-4" />
    </Button>
  </div>
</template>
