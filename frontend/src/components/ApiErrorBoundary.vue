<script setup lang="ts">
import { AlertCircle, RefreshCw, Server } from 'lucide-vue-next'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  error?: Error | string
  retry?: () => void
}

const props = defineProps<Props>()

const errorMessage =
  typeof props.error === 'string' ? props.error : props.error?.message || 'An error occurred'

const isConnectionError =
  errorMessage.includes('Failed to fetch') ||
  errorMessage.includes('Network') ||
  errorMessage.includes('CORS') ||
  errorMessage.includes('ERR_CONNECTION')

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
</script>

<template>
  <div class="flex items-center justify-center min-h-[400px] p-6">
    <Card class="max-w-md w-full">
      <CardHeader>
        <div class="flex items-center space-x-2">
          <Server v-if="isConnectionError" class="w-5 h-5 text-destructive" />
          <AlertCircle v-else class="w-5 h-5 text-destructive" />
          <CardTitle>{{ isConnectionError ? 'Cannot Connect to API' : 'Error' }}</CardTitle>
        </div>
        <CardDescription>
          {{ isConnectionError ? 'Unable to reach the backend server' : 'Something went wrong' }}
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <Alert variant="destructive">
          <AlertTitle>Error Details</AlertTitle>
          <AlertDescription class="mt-2 text-sm">
            {{ errorMessage }}
          </AlertDescription>
        </Alert>

        <div v-if="isConnectionError" class="space-y-3">
          <div class="text-sm text-muted-foreground space-y-2">
            <p class="font-medium">Please check:</p>
            <ul class="list-disc list-inside space-y-1 ml-2">
              <li>
                Backend server is running on
                <code class="text-xs bg-muted px-1 py-0.5 rounded">{{ apiUrl }}</code>
              </li>
              <li>CORS is properly configured</li>
              <li>Network connection is stable</li>
            </ul>
          </div>

          <div class="bg-muted p-3 rounded-md text-sm">
            <p class="font-medium mb-2">To start the backend:</p>
            <code class="text-xs">cd backend && go run cmd/api/main.go</code>
          </div>
        </div>

        <Button v-if="retry" @click="retry" class="w-full">
          <RefreshCw class="w-4 h-4 mr-2" />
          Retry
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
