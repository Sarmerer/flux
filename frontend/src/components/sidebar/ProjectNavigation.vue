<script setup lang="ts">
import { Activity, LayoutGrid, Settings as SettingsIcon, Table, Users, Workflow } from 'lucide-vue-next'
import { computed } from 'vue'

import { ROUTE_NAMES } from '@/constants/routes'
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
            :to="{ name: ROUTE_NAMES.PROJECT, params: { projectId: activeProjectId } }"
            :active-routes="[ROUTE_NAMES.PROJECT]"
          >
            <LayoutGrid class="w-4 h-4" />
            <span>Overview</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            :to="{ name: ROUTE_NAMES.PROJECT_TABLES, params: { projectId: activeProjectId } }"
            :active-routes="[ROUTE_NAMES.PROJECT_TABLES, ROUTE_NAMES.PROJECT_TABLE_DETAIL]"
          >
            <Table class="w-4 h-4" />
            <span>Tables</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            :to="{ name: ROUTE_NAMES.PROJECT_WORKFLOWS, params: { projectId: activeProjectId } }"
            :active-routes="[ROUTE_NAMES.PROJECT_WORKFLOWS, ROUTE_NAMES.PROJECT_WORKFLOW_DETAIL]"
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
            :to="{ name: ROUTE_NAMES.PROJECT_MEMBERS, params: { projectId: activeProjectId } }"
            :active-routes="[ROUTE_NAMES.PROJECT_MEMBERS]"
          >
            <Users class="w-4 h-4" />
            <span>Members</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton
            :to="{ name: ROUTE_NAMES.PROJECT_ACTIVITY, params: { projectId: activeProjectId } }"
            :active-routes="[ROUTE_NAMES.PROJECT_ACTIVITY]"
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
