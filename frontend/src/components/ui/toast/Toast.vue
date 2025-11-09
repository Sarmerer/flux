<script setup lang="ts">
import { CircleCheck, CircleX, Info, TriangleAlert, X } from 'lucide-vue-next'
import { computed } from 'vue'

export interface ToastProps {
  id?: string
  title?: string
  description?: string
  variant?: 'default' | 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

const props = withDefaults(defineProps<ToastProps>(), {
  variant: 'default',
  duration: 4000,
})

const emit = defineEmits<{
  close: []
}>()

const icon = computed(() => {
  switch (props.variant) {
    case 'success':
      return CircleCheck
    case 'error':
      return CircleX
    case 'warning':
      return TriangleAlert
    case 'info':
      return Info
    default:
      return null
  }
})

const iconClasses = computed(() => {
  switch (props.variant) {
    case 'success':
      return 'text-green-500'
    case 'error':
      return 'text-red-500'
    case 'warning':
      return 'text-amber-500'
    case 'info':
      return 'text-blue-500'
    default:
      return 'text-foreground'
  }
})
</script>

<template>
  <div
    class="group pointer-events-auto relative flex w-full items-center gap-3 overflow-hidden rounded-lg border border-border bg-background px-4 py-3 shadow-lg transition-all hover:shadow-xl"
  >
    <component :is="icon" v-if="icon" :class="['h-4 w-4 flex-shrink-0', iconClasses]" />
    <div class="flex-1 min-w-0">
      <div v-if="title" class="text-sm font-medium text-foreground leading-tight">
        {{ title }}
      </div>
      <div
        v-if="description"
        class="text-sm text-muted-foreground leading-tight"
        :class="{ 'mt-1': title }"
      >
        {{ description }}
      </div>
    </div>
    <button
      class="flex-shrink-0 rounded-md p-1 text-foreground/50 opacity-0 transition-all hover:text-foreground hover:bg-accent group-hover:opacity-100 focus:opacity-100"
      @click="emit('close')"
    >
      <X class="h-3.5 w-3.5" />
    </button>
  </div>
</template>
