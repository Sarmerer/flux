<script setup lang="ts">
import { Edit, Eye, MoreVertical, Plus, Table as TableIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useTableStore } from '@/stores/tables'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

import { useToast } from '@/composables/ui'

import { useRouteState } from './useRouteState'

const router = useRouter()
const activeProjectStore = useActiveProjectStore()
const tableStore = useTableStore()
const toast = useToast()
const { isActive } = useRouteState()

const activeProjectId = computed(() => activeProjectStore.activeProject?.id)

const handleCreateTable = () => {
  if (!activeProjectId.value) {
    toast.warning('No Project Selected', 'Please select a project first')
    return
  }
  router.push(`/projects/${activeProjectId.value}/tables?action=create`)
}

const handleViewTable = (tableId: string) => {
  router.push(`/projects/${activeProjectId.value}/tables/${tableId}/data`)
}

const handleEditTable = (tableId: string) => {
  router.push(`/projects/${activeProjectId.value}/tables/${tableId}/builder`)
}
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel class="flex items-center justify-between">
      <span class="text-xs font-semibold uppercase tracking-wider">Tables</span>
      <Button
        variant="ghost"
        size="sm"
        class="h-5 w-5 p-0 opacity-60 hover:opacity-100 transition-opacity group-data-[collapsible=icon]:hidden"
        @click="handleCreateTable"
      >
        <Plus class="w-3.5 h-3.5" />
      </Button>
    </SidebarGroupLabel>
    <SidebarGroupContent>
      <SidebarMenu>
        <div v-if="tableStore.isLoading" class="px-3 py-2 text-xs text-muted-foreground">
          Loading tables...
        </div>
        <SidebarMenuItem v-for="table in tableStore.tables" :key="table.id">
          <div class="group relative w-full">
            <SidebarMenuButton
              @click="handleViewTable(table.id)"
              :isActive="isActive(`/projects/${activeProjectId}/tables/${table.id}`)"
              class="pr-8"
            >
              <TableIcon class="w-4 h-4 flex-shrink-0" />
              <span class="flex-1 truncate">{{ table.name }}</span>
            </SidebarMenuButton>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  class="absolute right-1 top-1 h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity group-data-[collapsible=icon]:hidden"
                >
                  <MoreVertical class="w-3.5 h-3.5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem @click="handleViewTable(table.id)">
                  <Eye class="w-4 h-4 mr-2" />
                  View Data
                </DropdownMenuItem>
                <DropdownMenuItem @click="handleEditTable(table.id)">
                  <Edit class="w-4 h-4 mr-2" />
                  Edit Schema
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </SidebarMenuItem>
        <SidebarMenuItem v-if="!tableStore.isLoading && !tableStore.tables.length">
          <div class="px-3 py-2">
            <p class="text-xs text-muted-foreground mb-2">No tables yet</p>
            <Button
              variant="ghost"
              size="sm"
              class="w-full justify-start h-8 text-xs"
              @click="handleCreateTable"
            >
              <Plus class="w-3.5 h-3.5 mr-2" />
              Create your first table
            </Button>
          </div>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
