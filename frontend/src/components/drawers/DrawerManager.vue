<script setup lang="ts">
import { Sheet, SheetContent } from '@/components/ui/sheet'
import { useDrawerService } from '@/composables/useDrawerService'

const { drawerStack, closeDrawer } = useDrawerService()
</script>

<template>
  <Sheet
    v-for="(layer, index) in drawerStack"
    :key="layer.id"
    :open="true"
    @update:open="(open) => !open && closeDrawer(layer.id)"
  >
    <SheetContent :class="layer.width || 'w-[600px] sm:max-w-[600px]'" :style="{ zIndex: 50 + index }">
      <component
        :is="layer.component"
        v-bind="layer.props"
        @save="
          (data: any) => {
            layer.onSave?.(data)
            closeDrawer(layer.id)
          }
        "
        @close="() => closeDrawer(layer.id)"
      />
    </SheetContent>
  </Sheet>
</template>
