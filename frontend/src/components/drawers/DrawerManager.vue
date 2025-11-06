<script setup lang="ts">
import { ref, watch } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import { Sheet, SheetContent } from '@/components/ui/sheet'
import { useDrawerService } from '@/composables/useDrawerService'

const { drawerStack, closeDrawer } = useDrawerService()

const openStates = ref<Record<string, boolean>>({})
const componentRefs = ref<Record<string, ComponentPublicInstance | null>>({})

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

const handleOpenChange = async (layerId: string, open: boolean) => {
  if (!open) {
    const componentRef = componentRefs.value[layerId]

    if (componentRef && 'confirmClose' in componentRef && typeof (componentRef as any).confirmClose === 'function') {
      const canClose = await (componentRef as any).confirmClose()
      if (!canClose) {
        return
      }
    }

    openStates.value[layerId] = false
    setTimeout(() => {
      closeDrawer(layerId)
      delete openStates.value[layerId]
      delete componentRefs.value[layerId]
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
        :ref="(el: any) => componentRefs[layer.id] = el"
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
