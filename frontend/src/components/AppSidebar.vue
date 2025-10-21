<script setup lang="ts">
import { User } from 'lucide-vue-next'
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/projects'
import { useSidebarItemsStore } from '@/stores/ui/sidebar-items'

import { Badge } from '@/components/ui/badge'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

const router = useRouter()
const route = useRoute()
const sidebarStore = useSidebarItemsStore()
const authStore = useAuthStore()
const projectStore = useProjectStore()

const items = computed(() => sidebarStore.items)
const projectItems = computed(() => sidebarStore.projectItems)
const userItems = computed(() => sidebarStore.userItems)

// Get current project ID from route or store
const currentProjectId = computed(() => {
  return (route.params.projectId as string) || projectStore.currentProject?.id
})

// Watch for project changes and update sidebar
watch(
  currentProjectId,
  (newProjectId) => {
    if (!newProjectId) {
      sidebarStore.clearProjectItems()
    }
    // Note: The actual counts will be updated by the individual pages (Tables.vue, ProjectDetail.vue, etc.)
  },
  { immediate: true }
)

const isActive = (itemUrl: string) => {
  return route.path === itemUrl || route.path.startsWith(itemUrl + '/')
}

const handleItemClick = async (item: any) => {
  if (item.action === 'logout') {
    await authStore.logout()
    router.push('/login')
  } else {
    router.push(item.url)
  }
}
</script>

<template>
  <Sidebar collapsible="icon">
    <SidebarHeader class="p-4">
      <div class="flex items-center space-x-2">
        <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
          <span class="text-white font-bold text-sm">F</span>
        </div>
      </div>
    </SidebarHeader>

    <SidebarContent>
      <!-- Main Navigation -->
      <SidebarGroup>
        <SidebarGroupLabel>Application</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in items" :key="item.id" v-show="!item.hidden">
              <SidebarMenuButton
                @click="handleItemClick(item)"
                :isActive="isActive(item.url)"
                :disabled="item.disabled"
              >
                <component :is="item.icon" />
                <span>{{ item.title }}</span>
                <Badge v-if="item.badge" :variant="item.badgeVariant || 'default'" class="ml-auto">
                  {{ item.badge }}
                </Badge>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <!-- Project Navigation (only show when in a project) -->
      <SidebarGroup v-if="currentProjectId && projectItems.length > 0">
        <SidebarGroupLabel>Project</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in projectItems" :key="item.id" v-show="!item.hidden">
              <SidebarMenuButton
                @click="handleItemClick(item)"
                :isActive="isActive(item.url)"
                :disabled="item.disabled"
              >
                <component :is="item.icon" />
                <span>{{ item.title }}</span>
                <Badge v-if="item.badge" :variant="item.badgeVariant || 'default'" class="ml-auto">
                  {{ item.badge }}
                </Badge>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <!-- User Section -->
    <SidebarFooter class="p-4">
      <div class="space-y-2">
        <div class="flex items-center space-x-2 text-sm text-muted-foreground">
          <User class="w-4 h-4" />
          <span>{{ authStore.user?.name || 'User' }}</span>
        </div>
        <SidebarMenu>
          <SidebarMenuItem v-for="item in userItems" :key="item.id">
            <SidebarMenuButton
              @click="handleItemClick(item)"
              class="text-sm"
              :isActive="item.id !== 'logout' && isActive(item.url)"
            >
              <component :is="item.icon" />
              <span>{{ item.title }}</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
