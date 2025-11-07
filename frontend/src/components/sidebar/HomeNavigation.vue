<script setup lang="ts">
import { Home, LayoutGrid } from 'lucide-vue-next'

import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
</script>

<template>
  <SidebarGroup>
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton :to="{ name: 'Dashboard' }" :active-routes="['Dashboard']">
            <Home class="w-4 h-4" />
            <span>Home</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton :to="{ name: 'Projects' }" :active-routes="['Projects', 'ProjectsNew']">
            <LayoutGrid class="w-4 h-4" />
            <span>All Projects</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
