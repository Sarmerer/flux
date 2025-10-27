<script setup lang="ts">
import { ChevronRight, LogOut, Settings as SettingsIcon, User } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { SidebarFooter } from '@/components/ui/sidebar'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <SidebarFooter class="p-3 border-t">
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" class="w-full justify-start group-data-[collapsible=icon]:justify-center">
          <div class="w-7 h-7 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
            <User class="w-4 h-4 text-primary" />
          </div>
          <span class="flex-1 truncate text-left ml-2 group-data-[collapsible=icon]:hidden">{{ authStore.user?.name || 'User' }}</span>
          <ChevronRight class="w-3 h-3 opacity-50 rotate-90 group-data-[collapsible=icon]:hidden" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" class="w-56">
        <DropdownMenuItem @click="router.push('/settings')">
          <SettingsIcon class="w-4 h-4 mr-2" />
          Account Settings
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem @click="handleLogout">
          <LogOut class="w-4 h-4 mr-2" />
          Logout
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </SidebarFooter>
</template>
