<script setup lang="ts">
import HomeNavigation from './HomeNavigation.vue'
import ProjectNavigation from './ProjectNavigation.vue'
import { computed } from 'vue'

import { useActiveProjectStore } from '@/stores/activeProject'

import { Sidebar, SidebarContent, SidebarFooter, SidebarTrigger } from '@/components/ui/sidebar'

const activeProjectStore = useActiveProjectStore()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)
</script>

<template>
  <Sidebar collapsible="icon" class="pt-14 border-r">
    <SidebarContent>
      <template v-if="!activeProjectId">
        <HomeNavigation />
      </template>

      <template v-else>
        <ProjectNavigation />
      </template>
    </SidebarContent>

    <SidebarFooter class="p-3 border-t">
      <SidebarTrigger />
    </SidebarFooter>
  </Sidebar>
</template>
