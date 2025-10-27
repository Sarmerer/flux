<script setup lang="ts">
import { Table as TableIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useTableStore } from '@/stores/tables'

import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

import { useRouteState } from './useRouteState'

const router = useRouter()
const activeProjectStore = useActiveProjectStore()
const tableStore = useTableStore()
const { isActive } = useRouteState()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)

const handleViewTable = (tableId: string) => {
  router.push(`/projects/${activeProjectId.value}/tables/${tableId}`)
}
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel>
      <span class="text-xs font-semibold uppercase tracking-wider">Tables</span>
    </SidebarGroupLabel>
    <SidebarGroupContent>
      <SidebarMenu>
        <div v-if="tableStore.isLoading" class="px-3 py-2 text-xs text-muted-foreground">
          Loading tables...
        </div>
        <SidebarMenuItem v-for="table in tableStore.tables" :key="table.id">
          <SidebarMenuButton
            @click="handleViewTable(table.id)"
            :isActive="isActive(`/projects/${activeProjectId}/tables/${table.id}`)"
          >
            <TableIcon class="w-4 h-4 flex-shrink-0" />
            <span class="flex-1 truncate">{{ table.name }}</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem v-if="!tableStore.isLoading && !tableStore.tables.length">
          <div class="px-3 py-2 text-xs text-muted-foreground">
            No tables yet
          </div>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
