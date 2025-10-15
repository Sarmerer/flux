<script setup lang="ts">
import { ref, computed } from 'vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  Activity, 
  Search, 
  Calendar,
  User,
  Database,
  Table,
  Workflow,
  CheckCircle,
  XCircle,
  AlertCircle,
  Clock,
  RefreshCw
} from 'lucide-vue-next'


const searchQuery = ref('')
const filterType = ref('all')
const filterStatus = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)

const activities = ref([
  {
    id: '1',
    type: 'table_created',
    message: 'Table "users" was created',
    entity_type: 'table',
    entity_id: '1',
    entity_name: 'users',
    details: {
      columns: 5,
      schema: 'public'
    },
    user_id: 'user1',
    user_name: 'John Doe',
    created_at: '2024-01-21T10:30:00Z',
    status: 'completed'
  },
  {
    id: '2',
    type: 'workflow_triggered',
    message: 'Workflow "New User Welcome" was triggered',
    entity_type: 'workflow',
    entity_id: '1',
    entity_name: 'New User Welcome',
    details: {
      trigger_data: { user_id: '123' },
      execution_time: '2.3s'
    },
    user_id: 'system',
    user_name: 'System',
    created_at: '2024-01-21T10:25:00Z',
    status: 'completed'
  },
  {
    id: '3',
    type: 'row_created',
    message: 'New row created in table "users"',
    entity_type: 'row',
    entity_id: '123',
    entity_name: 'users',
    details: {
      row_data: { email: 'newuser@example.com', name: 'New User' }
    },
    user_id: 'user1',
    user_name: 'John Doe',
    created_at: '2024-01-21T10:20:00Z',
    status: 'completed'
  },
  {
    id: '4',
    type: 'workflow_failed',
    message: 'Workflow "Order Processing" failed to execute',
    entity_type: 'workflow',
    entity_id: '2',
    entity_name: 'Order Processing',
    details: {
      error: 'Database connection timeout',
      retry_count: 3
    },
    user_id: 'system',
    user_name: 'System',
    created_at: '2024-01-21T09:45:00Z',
    status: 'failed'
  },
  {
    id: '5',
    type: 'database_connected',
    message: 'Database connection established',
    entity_type: 'database',
    entity_id: '1',
    entity_name: 'main_db',
    details: {
      host: 'localhost:5432',
      database: 'flowdb'
    },
    user_id: 'user1',
    user_name: 'John Doe',
    created_at: '2024-01-21T09:30:00Z',
    status: 'completed'
  }
])

const activityTypes = [
  { value: 'all', label: 'All Activities' },
  { value: 'table_created', label: 'Table Created' },
  { value: 'table_updated', label: 'Table Updated' },
  { value: 'table_deleted', label: 'Table Deleted' },
  { value: 'row_created', label: 'Row Created' },
  { value: 'row_updated', label: 'Row Updated' },
  { value: 'row_deleted', label: 'Row Deleted' },
  { value: 'workflow_triggered', label: 'Workflow Triggered' },
  { value: 'workflow_completed', label: 'Workflow Completed' },
  { value: 'workflow_failed', label: 'Workflow Failed' },
  { value: 'database_connected', label: 'Database Connected' },
  { value: 'database_disconnected', label: 'Database Disconnected' }
]

const statusOptions = [
  { value: 'all', label: 'All Statuses' },
  { value: 'completed', label: 'Completed' },
  { value: 'failed', label: 'Failed' },
  { value: 'pending', label: 'Pending' },
  { value: 'in_progress', label: 'In Progress' }
]

const filteredActivities = computed(() => {
  let filtered = activities.value
  
  if (searchQuery.value) {
    filtered = filtered.filter(activity => 
      activity.message.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      activity.entity_name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (activity.user_name && activity.user_name.toLowerCase().includes(searchQuery.value.toLowerCase()))
    )
  }
  
  if (filterType.value !== 'all') {
    filtered = filtered.filter(activity => activity.type === filterType.value)
  }
  
  if (filterStatus.value !== 'all') {
    filtered = filtered.filter(activity => activity.status === filterStatus.value)
  }
  
  return filtered.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

const paginatedActivities = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredActivities.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(filteredActivities.value.length / pageSize.value)
})

const getActivityIcon = (type: string) => {
  switch (type) {
    case 'table_created':
    case 'table_updated':
    case 'table_deleted':
      return Table
    case 'row_created':
    case 'row_updated':
    case 'row_deleted':
      return Database
    case 'workflow_triggered':
    case 'workflow_completed':
    case 'workflow_failed':
      return Workflow
    case 'database_connected':
    case 'database_disconnected':
      return Database
    default:
      return Activity
  }
}

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'completed':
      return CheckCircle
    case 'failed':
      return XCircle
    case 'pending':
    case 'in_progress':
      return Clock
    default:
      return AlertCircle
  }
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'completed':
      return 'text-green-600'
    case 'failed':
      return 'text-red-600'
    case 'pending':
    case 'in_progress':
      return 'text-yellow-600'
    default:
      return 'text-gray-600'
  }
}

const getStatusBadgeVariant = (status: string) => {
  switch (status) {
    case 'completed':
      return 'default'
    case 'failed':
      return 'destructive'
    case 'pending':
    case 'in_progress':
      return 'secondary'
    default:
      return 'outline'
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatRelativeTime = (dateString: string) => {
  const now = new Date()
  const date = new Date(dateString)
  const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60))
  
  if (diffInMinutes < 1) return 'Just now'
  if (diffInMinutes < 60) return `${diffInMinutes}m ago`
  
  const diffInHours = Math.floor(diffInMinutes / 60)
  if (diffInHours < 24) return `${diffInHours}h ago`
  
  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 7) return `${diffInDays}d ago`
  
  return formatDate(dateString)
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Activity Log</h1>
        <p class="text-gray-600">Monitor all activities and events in your project</p>
      </div>
      <Button variant="outline" size="sm">
        <RefreshCw class="w-4 h-4 mr-2" />
        Refresh
      </Button>
    </div>

    <!-- Filters -->
    <Card>
      <CardContent class="p-4">
        <div class="flex items-center space-x-4">
          <div class="relative flex-1">
            <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
            <Input
              v-model="searchQuery"
              placeholder="Search activities..."
              class="pl-10"
            />
          </div>
          <Select v-model="filterType">
            <SelectTrigger class="w-48">
              <SelectValue placeholder="Filter by type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="type in activityTypes" :key="type.value" :value="type.value">
                {{ type.label }}
              </SelectItem>
            </SelectContent>
          </Select>
          <Select v-model="filterStatus">
            <SelectTrigger class="w-48">
              <SelectValue placeholder="Filter by status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="status in statusOptions" :key="status.value" :value="status.value">
                {{ status.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </CardContent>
    </Card>

    <!-- Activity List -->
    <Card>
      <CardContent class="p-0">
        <div v-if="paginatedActivities.length === 0" class="text-center py-12">
          <Activity class="mx-auto h-12 w-12 text-gray-400" />
          <h3 class="mt-2 text-sm font-medium text-gray-900">No activities found</h3>
          <p class="mt-1 text-sm text-gray-500">
            {{ searchQuery || filterType !== 'all' || filterStatus !== 'all' 
              ? 'Try adjusting your filters.' 
              : 'Activity will appear here as you work on your project.' 
            }}
          </p>
        </div>
        
        <div v-else class="divide-y divide-gray-200">
          <div 
            v-for="activity in paginatedActivities" 
            :key="activity.id"
            class="p-6 hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex items-start space-x-4">
              <div class="flex-shrink-0">
                <div class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                  <component :is="getActivityIcon(activity.type)" class="w-5 h-5 text-blue-600" />
                </div>
              </div>
              
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-2">
                    <p class="text-sm font-medium text-gray-900">{{ activity.message }}</p>
                    <Badge :variant="getStatusBadgeVariant(activity.status)">
                      {{ activity.status }}
                    </Badge>
                  </div>
                  <div class="flex items-center space-x-2 text-sm text-gray-500">
                    <component :is="getStatusIcon(activity.status)" :class="getStatusColor(activity.status)" class="w-4 h-4" />
                    <span>{{ formatRelativeTime(activity.created_at) }}</span>
                  </div>
                </div>
                
                <div class="mt-2 flex items-center space-x-4 text-sm text-gray-500">
                  <div class="flex items-center space-x-1">
                    <component :is="getActivityIcon(activity.entity_type)" class="w-4 h-4" />
                    <span>{{ activity.entity_name }}</span>
                  </div>
                  <div class="flex items-center space-x-1">
                    <User class="w-4 h-4" />
                    <span>{{ activity.user_name }}</span>
                  </div>
                  <div class="flex items-center space-x-1">
                    <Calendar class="w-4 h-4" />
                    <span>{{ formatDate(activity.created_at) }}</span>
                  </div>
                </div>
                
                <div v-if="activity.details && Object.keys(activity.details).length > 0" class="mt-3">
                  <details class="text-sm">
                    <summary class="cursor-pointer text-gray-600 hover:text-gray-900">
                      View Details
                    </summary>
                    <pre class="mt-2 p-3 bg-gray-100 rounded text-xs overflow-x-auto">{{ JSON.stringify(activity.details, null, 2) }}</pre>
                  </details>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-between">
      <div class="text-sm text-gray-700">
        Showing {{ (currentPage - 1) * pageSize + 1 }} to {{ Math.min(currentPage * pageSize, filteredActivities.length) }} of {{ filteredActivities.length }} results
      </div>
      <div class="flex items-center space-x-2">
        <Button
          variant="outline"
          size="sm"
          :disabled="currentPage === 1"
          @click="currentPage--"
        >
          Previous
        </Button>
        <span class="text-sm text-gray-700">
          Page {{ currentPage }} of {{ totalPages }}
        </span>
        <Button
          variant="outline"
          size="sm"
          :disabled="currentPage === totalPages"
          @click="currentPage++"
        >
          Next
        </Button>
      </div>
    </div>
  </div>
</template>
