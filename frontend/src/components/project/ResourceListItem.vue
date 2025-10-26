<script setup lang="ts">
import type { Component } from 'vue'
import { ExternalLink } from 'lucide-vue-next'

defineProps<{
  name: string
  description: string | null | undefined
  icon: Component
  iconBgColor: string
  iconColor: string
}>()

defineEmits<{
  click: []
}>()
</script>

<template>
  <div
    class="flex items-center justify-between p-4 border rounded-lg hover:bg-accent cursor-pointer"
    @click="$emit('click')"
  >
    <div class="flex items-center space-x-3">
      <div :class="`p-2 rounded-lg ${iconBgColor}`">
        <component :is="icon" :class="`h-5 w-5 ${iconColor}`" />
      </div>
      <div>
        <h4 class="font-medium">{{ name }}</h4>
        <p class="text-sm text-muted-foreground">
          {{ description || 'No description' }}
        </p>
      </div>
    </div>
    <ExternalLink class="h-4 w-4 text-muted-foreground" />
  </div>
</template>
