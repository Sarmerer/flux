<script setup lang="ts">
import { Activity } from 'lucide-vue-next'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Activity Log</CardTitle>
      <CardDescription>Recent activity in this project</CardDescription>
    </CardHeader>
    <CardContent>
      <div class="text-center py-8">
        <Activity class="mx-auto h-12 w-12 text-muted-foreground" />
        <h3 class="mt-2 text-sm font-medium text-foreground">No activity yet</h3>
        <p class="mt-1 text-sm text-muted-foreground">
          Activity will appear here as you work on your project.
        </p>
      </div>
    </CardContent>
  </Card>
</template>
