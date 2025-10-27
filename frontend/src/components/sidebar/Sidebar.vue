<script setup lang="ts">
import SidebarFooter from './Footer.vue'
import SidebarHeader from './Header.vue'
import HomeNavigation from './HomeNavigation.vue'
import ProjectNavigation from './ProjectNavigation.vue'
import { computed } from 'vue'

import { useActiveProjectStore } from '@/stores/activeProject'

import { Sidebar, SidebarContent } from '@/components/ui/sidebar'

const activeProjectStore = useActiveProjectStore()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)
</script>

<template>
  <Sidebar collapsible="icon">
    <SidebarHeader />

    <SidebarContent>
      <template v-if="!activeProjectId">
        <HomeNavigation />
      </template>

      <template v-else>
        <ProjectNavigation />
      </template>
    </SidebarContent>

    <SidebarFooter />
  </Sidebar>
</template>
