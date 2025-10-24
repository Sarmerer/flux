<script setup lang="ts">
import AppSidebar from '@/components/AppSidebar.vue'
import Toaster from '@/components/ui/toast/Toaster.vue'

import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
</script>

<template>
  <SidebarProvider>
    <AppSidebar />

    <main class="flex-1">
      <header class="p-4 border-b">
        <SidebarTrigger />
      </header>

      <RouterView v-slot="{ Component }">
        <component :is="Component" />
      </RouterView>
    </main>

    <!-- Global Toast Notifications -->
    <Toaster />
  </SidebarProvider>
</template>
