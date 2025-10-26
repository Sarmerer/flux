<script setup lang="ts">
import ErrorState from '@/components/ui/ErrorState.vue'
import LoadingWrapper from '@/components/ui/LoadingWrapper.vue'
import { ExternalLink, Plus, Table } from 'lucide-vue-next'
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
          <div
            v-for="table in tablesStore.tables"
            :key="table.id"
            class="flex items-center justify-between p-4 border rounded-lg hover:bg-accent cursor-pointer"
            @click="router.push(`/projects/${props.activeProject.id}/tables/${table.id}/data`)"
          >
            <div class="flex items-center space-x-3">
              <div class="p-2 bg-blue-100 rounded-lg">
                <Table class="h-5 w-5 text-blue-600" />
              </div>
              <div>
                <h4 class="font-medium">{{ table.name }}</h4>
                <p class="text-sm text-muted-foreground">
                  {{ table.description || 'No description' }}
                </p>
              </div>
            </div>
            <ExternalLink class="h-4 w-4 text-muted-foreground" />
          </div>
        </div>
      </LoadingWrapper>
    </CardContent>
  </Card>
</template>
