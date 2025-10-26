<script setup lang="ts">
import { Calendar, User } from 'lucide-vue-next'
import { Table, Workflow } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

import type { Project } from '@/types/api'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { useFormatting } from '@/composables/formatting'

const router = useRouter()
const { formatDate } = useFormatting()

const props = defineProps<{ activeProject: Project }>()

const onTableCreateClick = () => {
  router.push(`/projects/${props.activeProject.id}/tables`)
}

const onWorkflowCreateClick = () => {
  router.push(`/projects/${props.activeProject.id}/workflows`)
}
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
        <CardDescription>Get started with your project</CardDescription>
      </CardHeader>
      <CardContent class="space-y-3">
        <Button variant="outline" class="w-full justify-start" @click="onTableCreateClick">
          <Table class="w-4 h-4 mr-2" />
          Create New Table
        </Button>
        <Button variant="outline" class="w-full justify-start" @click="onWorkflowCreateClick">
          <Workflow class="w-4 h-4 mr-2" />
          Create Workflow
        </Button>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>Project Information</CardTitle>
        <CardDescription>Details about this project</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">Created</span>
          <div class="flex items-center space-x-1 text-sm text-muted-foreground">
            <Calendar class="w-4 h-4" />
            <span>{{ formatDate(activeProject.created_at) }}</span>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">Owner</span>
          <div class="flex items-center space-x-1 text-sm text-muted-foreground">
            <User class="w-4 h-4" />
            <span>{{ activeProject.owner_id || 'Unknown' }}</span>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">Status</span>
          <Badge variant="secondary">Active</Badge>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
