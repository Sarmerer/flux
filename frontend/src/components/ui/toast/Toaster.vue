<script setup lang="ts">
import { computed } from 'vue'
import { useToast } from '@/composables/ui'
import Toast from './Toast.vue'

const { toasts, removeToast } = useToast()

const visibleToasts = computed(() => toasts.value)
</script>

<template>
  <div
    class="fixed top-0 right-0 z-[100] flex max-h-screen w-full flex-col-reverse p-4 sm:top-auto sm:bottom-0 sm:right-0 sm:flex-col md:max-w-[420px] pointer-events-none"
  >
    <TransitionGroup
      name="toast"
      tag="div"
      class="flex flex-col-reverse gap-2 sm:flex-col"
    >
      <Toast
        v-for="toast in visibleToasts"
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
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
