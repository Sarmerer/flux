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
  <div class="h-full">
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
        <div class="border-b bg-muted/30">
          <div class="px-6 py-8">
            <div class="flex items-start justify-between mb-6">
              <div class="flex-1">
                <h1 class="text-3xl font-bold text-foreground mb-2">{{ activeProject.name }}</h1>
                <p class="text-muted-foreground text-sm max-w-2xl">
                  {{ activeProject.description || 'No description provided' }}
                </p>
              </div>
              <Button variant="outline" size="sm">
                <Settings class="w-4 h-4 mr-2" />
                Settings
              </Button>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <StatCard
                v-for="stat in statsCards"
                :key="stat.title"
                :title="stat.title"
                :value="stat.value"
                :description="stat.description"
                :icon="stat.icon"
              />
            </div>
          </div>

          <Tabs :model-value="currentTab" @update:model-value="onTabChange" class="px-6">
            <TabsList class="bg-transparent border-b-0 h-auto p-0">
              <TabsTrigger
                value="overview"
                class="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent"
              >
                Overview
              </TabsTrigger>
              <TabsTrigger
                value="tables"
                class="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent"
              >
                Tables
              </TabsTrigger>
              <TabsTrigger
                value="workflows"
                class="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent"
              >
                Workflows
              </TabsTrigger>
              <TabsTrigger
                value="activity"
                class="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent"
              >
                Activity
              </TabsTrigger>
            </TabsList>
          </Tabs>
        </div>

        <div class="p-6">
          <Tabs :model-value="currentTab" @update:model-value="onTabChange">
            <TabsContent value="overview" class="mt-0">
              <OverviewTab :active-project="activeProject" />
            </TabsContent>

            <TabsContent value="tables" class="mt-0">
              <TablesTab :active-project="activeProject" />
            </TabsContent>

            <TabsContent value="workflows" class="mt-0">
              <WorkflowsTab :active-project="activeProject" />
            </TabsContent>

            <TabsContent value="activity" class="mt-0">
              <ActivityTab />
            </TabsContent>
          </Tabs>
        </div>
      </template>
    </LoadingWrapper>
  </div>
</template>
