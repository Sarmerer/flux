<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import Toaster from '@/components/ui/toast/Toaster.vue'

import { AppSidebar } from '@/components/sidebar'

import { SidebarProvider } from '@/components/ui/sidebar'
</script>

<template>
  <SidebarProvider>
    <AppSidebar />

    <main class="flex-1 flex flex-col min-h-screen">
      <AppHeader />

      <div class="flex-1">
        <RouterView v-slot="{ Component }">
          <component :is="Component" />
        </RouterView>
      </div>
    </main>

    <Toaster />
  </SidebarProvider>
</template>
