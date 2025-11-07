<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { computed } from 'vue'
import type { RouteLocationRaw } from 'vue-router'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { cn } from '@/lib/utils'

import type { SidebarMenuButtonVariants } from '.'
import { sidebarMenuButtonVariants } from '.'

export interface SidebarMenuButtonProps {
  variant?: SidebarMenuButtonVariants['variant']
  size?: SidebarMenuButtonVariants['size']
  isActive?: boolean
  class?: HTMLAttributes['class']
  to?: RouteLocationRaw
  activeRoutes?: string[]
}

const props = withDefaults(defineProps<SidebarMenuButtonProps>(), {
  as: 'button',
  variant: 'default',
  size: 'default',
})

const router = useRouter()
const route = useRoute()

const isRouteActive = computed(() => {
  if (props.isActive !== undefined) {
    return props.isActive
  }

  if (props.activeRoutes && props.activeRoutes.length > 0) {
    return props.activeRoutes.includes(route.name as string)
  }

  if (props.to) {
    const resolved = router.resolve(props.to)
    return resolved.fullPath === route.fullPath
  }

  return false
})

const computedClass = computed(() =>
  cn(sidebarMenuButtonVariants({ variant: props.variant, size: props.size }), props.class)
)
</script>

<template>
  <RouterLink
    v-if="to"
    :to="to"
    :class="computedClass"
    :data-active="isRouteActive"
    data-slot="sidebar-menu-button"
    data-sidebar="menu-button"
    :data-size="size"
    v-bind="$attrs"
  >
    <slot />
  </RouterLink>

  <button
    v-else
    :class="computedClass"
    :data-active="isActive"
    data-slot="sidebar-menu-button"
    data-sidebar="menu-button"
    :data-size="size"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>
