<script setup lang="ts">
import { useToast } from '@/composables/ui'
import Toast from './Toast.vue'

const { toasts, removeToast } = useToast()
</script>

<template>
  <div
    class="fixed top-4 right-4 z-[100] flex max-h-screen w-full max-w-sm flex-col gap-2 pointer-events-none"
  >
    <TransitionGroup name="toast" tag="div" class="flex flex-col gap-2">
      <Toast
        v-for="toast in toasts"
        :key="toast.id"
        :id="toast.id"
        :title="toast.title"
        :description="toast.description"
        :variant="toast.variant"
        @close="removeToast(toast.id!)"
      />
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-leave-active {
  transition: all 0.15s cubic-bezier(0.4, 0, 1, 1);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.toast-move {
  transition: transform 0.2s ease-out;
}
</style>
