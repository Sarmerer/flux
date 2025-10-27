<script setup lang="ts">
import { ChevronRight, Plus, Search } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

import { useToast } from '@/composables/ui'

const router = useRouter()
const route = useRoute()
const activeProjectStore = useActiveProjectStore()
const toast = useToast()

const activeProject = computed(() => activeProjectStore.activeProject)

interface BreadcrumbItem {
  label: string
  path?: string
}

const breadcrumbs = computed<BreadcrumbItem[]>(() => {
  const crumbs: BreadcrumbItem[] = []

  if (route.path === '/') {
    crumbs.push({ label: 'Home' })
  } else if (route.path === '/projects') {
    crumbs.push({ label: 'Projects' })
  } else if (route.path === '/settings') {
    crumbs.push({ label: 'Settings' })
  } else if (activeProject.value) {
    crumbs.push({ label: 'Projects', path: '/projects' })
    crumbs.push({ label: activeProject.value.name, path: `/projects/${activeProject.value.id}` })

    if (route.path.includes('/tables')) {
      crumbs.push({ label: 'Tables' })
    } else if (route.path.includes('/workflows')) {
      crumbs.push({ label: 'Workflows' })
    } else if (route.path.includes('/members')) {
      crumbs.push({ label: 'Members' })
    } else if (route.path.includes('/activity')) {
      crumbs.push({ label: 'Activity' })
    }
  }

  return crumbs
})

const showQuickActions = computed(() => !!activeProject.value)

const onQuickAction = (action: string) => {
  if (!activeProject.value) {
    toast.warning('No Project Selected', 'Please select a project first')
    return
  }

  switch (action) {
    case 'table':
      router.push(`/projects/${activeProject.value.id}/tables?action=create`)
      break
    case 'workflow':
      router.push(`/projects/${activeProject.value.id}/workflows?action=create`)
      break
  }
}
</script>

<template>
  <header class="sticky top-0 z-10 flex h-14 items-center gap-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 px-6">
    <nav class="flex items-center gap-2 text-sm flex-1">
      <template v-for="(crumb, index) in breadcrumbs" :key="index">
        <button
          v-if="crumb.path"
          @click="router.push(crumb.path)"
          class="text-muted-foreground hover:text-foreground transition-colors"
        >
          {{ crumb.label }}
        </button>
        <span v-else class="text-foreground font-medium">{{ crumb.label }}</span>
        <ChevronRight v-if="index < breadcrumbs.length - 1" class="w-4 h-4 text-muted-foreground" />
      </template>
    </nav>

    <div class="flex items-center gap-3">
      <Button
        variant="outline"
        size="sm"
        class="gap-2 w-[220px] justify-start text-muted-foreground hover:text-foreground transition-colors"
      >
        <Search class="w-3.5 h-3.5 flex-shrink-0" />
        <span class="flex-1 text-left text-sm">Search...</span>
        <kbd
          class="inline-flex h-5 select-none items-center gap-0.5 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium"
        >
          <span class="text-xs">⌘</span>K
        </kbd>
      </Button>

      <DropdownMenu v-if="showQuickActions">
        <DropdownMenuTrigger asChild>
          <Button variant="default" size="sm" class="gap-2 shadow-sm">
            <Plus class="w-4 h-4" />
            <span>Create</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>Quick Actions</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="onQuickAction('table')">
            <Plus class="w-4 h-4 mr-2" />
            New Table
          </DropdownMenuItem>
          <DropdownMenuItem @click="onQuickAction('workflow')">
            <Plus class="w-4 h-4 mr-2" />
            New Workflow
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </header>
</template>
