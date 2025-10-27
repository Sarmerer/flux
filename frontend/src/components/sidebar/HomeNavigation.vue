<script setup lang="ts">
import { FolderOpen, Home, LayoutGrid } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

import { Button } from '@/components/ui/button'
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

import { useRouteState } from './useRouteState'

const router = useRouter()
const { isExactActive } = useRouteState()
</script>

<template>
  <SidebarGroup>
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton @click="router.push('/')" :isActive="isExactActive('/')">
            <Home class="w-4 h-4" />
            <span>Home</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            @click="router.push('/projects')"
            :isActive="isExactActive('/projects')"
          >
            <LayoutGrid class="w-4 h-4" />
            <span>All Projects</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>

  <SidebarGroup>
    <SidebarGroupContent>
      <div class="px-4 py-8 text-center space-y-3 group-data-[collapsible=icon]:hidden">
        <FolderOpen class="w-12 h-12 mx-auto text-muted-foreground/40" />
        <div>
          <p class="text-sm font-medium text-foreground mb-1">No project selected</p>
          <p class="text-xs text-muted-foreground leading-relaxed">
            Select a project to view tables and workflows
          </p>
        </div>
        <Button variant="default" size="sm" @click="router.push('/projects')" class="w-full shadow-sm">
          <LayoutGrid class="w-4 h-4 mr-2" />
          Browse Projects
        </Button>
      </div>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
