<script setup lang="ts">
import { Activity, LayoutGrid, Settings as SettingsIcon, Table, Users, Workflow } from 'lucide-vue-next'
import { computed } from 'vue'

import { useActiveProjectStore } from '@/stores/activeProject'

import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

const activeProjectStore = useActiveProjectStore()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)
</script>

<template>
  <SidebarGroup>
    <SidebarGroupContent>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton
            :to="{ name: 'ProjectOverview', params: { projectId: activeProjectId } }"
            :active-routes="['ProjectDetail', 'ProjectOverview']"
          >
            <LayoutGrid class="w-4 h-4" />
            <span>Overview</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            :to="{ name: 'Tables', params: { projectId: activeProjectId } }"
            :active-routes="['Tables', 'TableDetail']"
          >
            <Table class="w-4 h-4" />
            <span>Tables</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            :to="{ name: 'Workflows', params: { projectId: activeProjectId } }"
            :active-routes="['Workflows', 'WorkflowDetail']"
          >
            <Workflow class="w-4 h-4" />
            <span>Workflows</span>
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
            :to="{ name: 'ProjectMembers', params: { projectId: activeProjectId } }"
            :active-routes="['ProjectMembers']"
          >
            <Users class="w-4 h-4" />
            <span>Members</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            :to="{ name: 'Activity', params: { projectId: activeProjectId } }"
            :active-routes="['Activity']"
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
