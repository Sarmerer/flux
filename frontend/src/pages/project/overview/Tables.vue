<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import ResourceListItem from '@/components/project/ResourceListItem.vue'
import { Plus, Table } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

import { useTableStore } from '@/stores/tables'
import type { Project } from '@/types/api'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const router = useRouter()
const props = defineProps<{ activeProject: Project }>()

const tablesStore = useTableStore()

const onTableCreateClick = () => {
  router.push(`/projects/${props.activeProject.id}/tables`)
}
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-center justify-between">
        <div>
          <CardTitle>Tables</CardTitle>
          <CardDescription>Manage your database tables</CardDescription>
        </div>
        <Button @click="onTableCreateClick">
          <Plus class="w-4 h-4 mr-2" />
          New Table
        </Button>
      </div>
    </CardHeader>
    <CardContent>
      <LoadingWrapper :is-loading="tablesStore.isLoading" loading-text="Loading tables...">
        <ErrorState
          v-if="tablesStore.error"
          :error="tablesStore.error"
          title="Failed to load tables"
          :on-retry="() => tablesStore.loadByProjectId(props.activeProject.id)"
        />
        <div v-else class="space-y-4">
          <ResourceListItem
            v-for="table in tablesStore.tables"
            :key="table.id"
            :name="table.name"
            :description="table.description"
            :icon="Table"
            icon-bg-color="bg-blue-100"
            icon-color="text-blue-600"
            @click="router.push(`/projects/${props.activeProject.id}/tables/${table.id}/data`)"
          />
        </div>
      </LoadingWrapper>
    </CardContent>
  </Card>
</template>
