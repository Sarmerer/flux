<script setup lang="ts">
import AppSidebar from '@/components/AppSidebar.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'

const route = useRoute()

// Which pages should be cached
const cachedPages = ['Home', 'About']

const isCached = computed(() => cachedPages.includes(route.name as string))
</script>

<template>
  <SidebarProvider>
    <AppSidebar />

    <main class="flex-1">
      <header class="p-4 border-b">
        <SidebarTrigger />
      </header>

      <RouterView v-slot="{ Component }">
        <KeepAlive>
          <component :is="Component" />
        </KeepAlive>
      </RouterView>
    </main>
  </SidebarProvider>
</template>
