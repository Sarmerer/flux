<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import DrawerManager from '@/components/drawers/DrawerManager.vue'
import Toaster from '@/components/ui/toast/Toaster.vue'

import { AppSidebar } from '@/components/sidebar'

import { SidebarProvider } from '@/components/ui/sidebar'
</script>

<template>
  <SidebarProvider>
    <div class="flex flex-col min-h-screen w-full">
      <AppHeader />

      <div class="flex flex-1 pt-14">
        <AppSidebar />

        <main class="flex-1">
          <RouterView v-slot="{ Component }">
            <component :is="Component" />
          </RouterView>
        </main>
      </div>
    </div>

    <Toaster />
    <DrawerManager />
  </SidebarProvider>
</template>
