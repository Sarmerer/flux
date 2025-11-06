<script setup lang="ts">
import { ref, watch } from 'vue'
import { Sheet, SheetContent } from '@/components/ui/sheet'
import { useDrawerService } from '@/composables/useDrawerService'

const { drawerStack, closeDrawer } = useDrawerService()

const openStates = ref<Record<string, boolean>>({})

watch(
  drawerStack,
  (newStack) => {
    newStack.forEach((layer) => {
      if (!(layer.id in openStates.value)) {
        openStates.value[layer.id] = true
      }
    })
  },
  { immediate: true }
)

const handleOpenChange = (layerId: string, open: boolean) => {
  if (!open) {
    openStates.value[layerId] = false
    setTimeout(() => {
      closeDrawer(layerId)
      delete openStates.value[layerId]
    }, 200)
  }
}
</script>

<template>
  <Sheet
    v-for="(layer, index) in drawerStack"
    :key="layer.id"
    :open="openStates[layer.id] ?? true"
    @update:open="(open) => handleOpenChange(layer.id, open)"
  >
    <SheetContent :class="layer.width || 'w-[600px] sm:max-w-[600px]'" :style="{ zIndex: 50 + index }">
      <component
        :is="layer.component"
        v-bind="layer.props"
        @save="
          (data: any) => {
            layer.onSave?.(data)
            handleOpenChange(layer.id, false)
          }
        "
        @close="() => handleOpenChange(layer.id, false)"
      />
    </SheetContent>
  </Sheet>
</template>
