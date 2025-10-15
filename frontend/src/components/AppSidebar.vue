<script setup lang="ts">
import { User } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useProjectsStore } from '@/stores/projects'
import { useSidebarItemsStore } from '@/stores/ui/sidebar-items'

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
const { items, projectItems, userItems } = useSidebarItemsStore()
const authStore = useAuthStore()

const currentProjectId = computed(() => {
  const route = router.currentRoute.value
  return route.params.projectId as string
})

const handleItemClick = async (item: any) => {
  if (item.action === 'logout') {
    await authStore.logout()
    router.push('/login')
  } else {
    router.push(item.url)
  }
}

// Set project items when we have a project ID
if (currentProjectId.value) {
  useSidebarItemsStore().setProjectItems(currentProjectId.value)
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
            <SidebarMenuItem v-for="item in items" :key="item.id">
              <SidebarMenuButton @click="handleItemClick(item)">
                <component :is="item.icon" />
                <span>{{ item.title }}</span>
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
            <SidebarMenuItem v-for="item in projectItems" :key="item.id">
              <SidebarMenuButton @click="handleItemClick(item)">
                <component :is="item.icon" />
                <span>{{ item.title }}</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <!-- User Section -->
    <SidebarFooter class="p-4">
      <div class="space-y-2">
        <div class="flex items-center space-x-2 text-sm text-gray-600">
          <User class="w-4 h-4" />
          <span>{{ authStore.user?.name || 'User' }}</span>
        </div>
        <SidebarMenu>
          <SidebarMenuItem v-for="item in userItems" :key="item.id">
            <SidebarMenuButton @click="handleItemClick(item)" class="text-sm">
              <component :is="item.icon" />
              <span>{{ item.title }}</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
