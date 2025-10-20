<script setup lang="ts">
import SimpleSidebar from '@/components/SimpleSidebar.vue'
import Toaster from '@/components/ui/toast/Toaster.vue'

import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
</script>

<template>
  <SidebarProvider>
    <SimpleSidebar />

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

    <!-- Global Toast Notifications -->
    <Toaster />
  </SidebarProvider>
</template>
