<script setup lang="ts">
import { Home, LayoutGrid } from 'lucide-vue-next'

import { ROUTE_NAMES } from '@/constants/routes'
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
          <SidebarMenuButton :to="{ name: ROUTE_NAMES.HOME }" :active-routes="[ROUTE_NAMES.HOME]">
            <Home class="w-4 h-4" />
            <span>Home</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton :to="{ name: ROUTE_NAMES.PROJECTS }" :active-routes="[ROUTE_NAMES.PROJECTS, ROUTE_NAMES.PROJECTS_NEW]">
            <LayoutGrid class="w-4 h-4" />
            <span>All Projects</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
