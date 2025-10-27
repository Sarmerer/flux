<script setup lang="ts">
import { Workflow } from 'lucide-vue-next'

import ProjectSelector from './ProjectSelector.vue'
import { SidebarHeader, SidebarTrigger } from '@/components/ui/sidebar'
</script>

<template>
  <SidebarHeader class="p-4 space-y-3 group-data-[collapsible=icon]:px-2">
    <div class="flex items-center justify-between group-data-[collapsible=icon]:flex-col group-data-[collapsible=icon]:gap-3">
      <div class="flex items-center space-x-2.5 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:w-full">
        <div class="w-8 h-8 bg-gradient-to-br from-primary to-primary/80 rounded-lg flex items-center justify-center flex-shrink-0 shadow-sm group-data-[collapsible=icon]:w-10 group-data-[collapsible=icon]:h-10 transition-all">
          <Workflow class="w-4.5 h-4.5 text-primary-foreground group-data-[collapsible=icon]:w-5 group-data-[collapsible=icon]:h-5 transition-all" :stroke-width="2.5" />
        </div>
        <span class="font-semibold text-foreground group-data-[collapsible=icon]:hidden tracking-tight">
          Flow
        </span>
      </div>
      <SidebarTrigger class="group-data-[collapsible=icon]:w-full" />
    </div>

    <ProjectSelector class="group-data-[collapsible=icon]:hidden" />
  </SidebarHeader>
</template>
