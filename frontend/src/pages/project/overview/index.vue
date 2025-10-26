<script setup lang="ts">
import ActivityTab from './Activity.vue'
import OverviewTab from './Overview.vue'
import TablesTab from './Tables.vue'
import WorkflowsTab from './Workflows.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import StatCard from '@/components/project/StatCard.vue'
import { Activity, Settings, Table, Workflow } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useActiveProjectStore } from '@/stores/activeProject'
import { useTableStore } from '@/stores/tables'
import { useWorkflowStore } from '@/stores/workflows'

import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const route = useRoute()
const router = useRouter()

const activeProjectStore = useActiveProjectStore()
const activeProject = computed(() => activeProjectStore.activeProject)
const activeProjectId = computed(() => activeProject.value?.id)

const tablesStore = useTableStore()
const workflowsStore = useWorkflowStore()

const currentTab = computed(() => (route.meta.tab as string) || 'overview')

const statsCards = computed(() => [
  {
    title: 'Tables',
    value: tablesStore.length,
    description: 'Database tables',
    icon: Table,
  },
  {
    title: 'Workflows',
    value: workflowsStore.length,
    description: 'Active workflows',
    icon: Workflow,
  },
  {
    title: 'Last Activity',
    value: '2 hours ago',
    description: 'Time since last update',
    icon: Activity,
  },
])

const onTabChange = (tab: string | number) => {
  const tabString = String(tab)
  if (tabString === 'overview') {
    router.push(`/projects/${activeProjectId.value}/overview`)
  } else {
    router.push(`/projects/${activeProjectId.value}/overview/${tabString}`)
  }
}
</script>

<template>
  <div class="p-6 space-y-6">
    <LoadingWrapper :is-loading="activeProjectStore.isLoading" loading-text="Loading project...">
      <ErrorState
        v-if="activeProjectStore.error"
        :error="activeProjectStore.error"
        title="Project not found"
        description="The project you're looking for doesn't exist."
      >
        <Button size="sm" @click="router.push('/projects')">Back to Projects</Button>
      </ErrorState>

      <template v-else-if="activeProject">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-foreground">{{ activeProject.name }}</h1>
            <p class="text-muted-foreground">{{ activeProject.description || 'No description' }}</p>
          </div>
          <div class="flex items-center space-x-2">
            <Button variant="outline">
              <Settings class="w-4 h-4 mr-2" />
              Settings
            </Button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
          <StatCard
            v-for="stat in statsCards"
            :key="stat.title"
            :title="stat.title"
            :value="stat.value"
            :description="stat.description"
            :icon="stat.icon"
          />
        </div>

        <Tabs :model-value="currentTab" @update:model-value="onTabChange" class="space-y-6">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="tables">Tables</TabsTrigger>
            <TabsTrigger value="workflows">Workflows</TabsTrigger>
            <TabsTrigger value="activity">Activity</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" class="space-y-6">
            <OverviewTab :active-project="activeProject" />
          </TabsContent>

          <TabsContent value="tables">
            <TablesTab :active-project="activeProject" />
          </TabsContent>

          <TabsContent value="workflows">
            <WorkflowsTab :active-project="activeProject" />
          </TabsContent>

          <TabsContent value="activity">
            <ActivityTab />
          </TabsContent>
        </Tabs>
      </template>
    </LoadingWrapper>
  </div>
</template>
