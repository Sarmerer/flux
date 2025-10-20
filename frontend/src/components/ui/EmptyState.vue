<script setup lang="ts">
import { FileQuestion } from 'lucide-vue-next'
import { type Component } from 'vue'

import { Button } from '@/components/ui/button'

export interface EmptyStateProps {
  icon?: Component
  title?: string
  description?: string
  actionLabel?: string
  actionIcon?: Component
}

const props = withDefaults(defineProps<EmptyStateProps>(), {
  icon: FileQuestion,
  title: 'No data found',
  description: 'Get started by creating your first item',
})

const emit = defineEmits<{
  action: []
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center p-12 text-center">
    <div class="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
      <component :is="icon" class="w-8 h-8 text-muted-foreground" />
    </div>
    <h3 class="text-lg font-semibold text-foreground mb-2">{{ title }}</h3>
    <p class="text-sm text-muted-foreground mb-6 max-w-sm">{{ description }}</p>
    <Button v-if="actionLabel" @click="emit('action')">
      <component :is="actionIcon" v-if="actionIcon" class="w-4 h-4 mr-2" />
      {{ actionLabel }}
    </Button>
  </div>
</template>
