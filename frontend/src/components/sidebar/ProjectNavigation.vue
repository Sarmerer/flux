<script setup lang="ts">
import { Activity, LayoutGrid, Settings as SettingsIcon, Users } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'

import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

import { useRouteState } from './useRouteState'

const router = useRouter()
const activeProjectStore = useActiveProjectStore()
const { isActive } = useRouteState()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)
</script>

<template>
  <SidebarGroup>
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton
            @click="router.push(`/projects/${activeProjectId}`)"
            :isActive="isActive(`/projects/${activeProjectId}/overview`)"
          >
            <LayoutGrid class="w-4 h-4" />
            <span>Overview</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>

  <SidebarGroup>
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton
            @click="router.push(`/projects/${activeProjectId}/members`)"
            :isActive="isActive(`/projects/${activeProjectId}/members`)"
          >
            <Users class="w-4 h-4" />
            <span>Members</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            @click="router.push(`/projects/${activeProjectId}/activity`)"
            :isActive="isActive(`/projects/${activeProjectId}/activity`)"
          >
            <Activity class="w-4 h-4" />
            <span>Activity</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton :isActive="false">
            <SettingsIcon class="w-4 h-4" />
            <span>Project Settings</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
