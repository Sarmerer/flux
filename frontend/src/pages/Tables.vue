<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  Table, 
  Plus, 
  Search, 
  Edit, 
  Trash2, 
  Database,
  MoreHorizontal,
  Eye,
  Settings,
  Calendar
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const projectId = computed(() => route.params.projectId as string)
const searchQuery = ref('')
const isCreateDialogOpen = ref(false)
const isEditDialogOpen = ref(false)
const selectedTable = ref(null)

const tables = ref([
  {
    id: '1',
    name: 'users',
    description: 'User accounts and profiles',
    columns: 5,
    rows: 1250,
    created_at: '2024-01-15T10:30:00Z',
    updated_at: '2024-01-20T14:22:00Z'
  },
  {
    id: '2',
    name: 'products',
    description: 'Product catalog and inventory',
    columns: 8,
    rows: 3400,
    created_at: '2024-01-16T09:15:00Z',
    updated_at: '2024-01-19T16:45:00Z'
  },
  {
    id: '3',
    name: 'orders',
    description: 'Customer orders and transactions',
    columns: 6,
    rows: 890,
    created_at: '2024-01-17T11:20:00Z',
    updated_at: '2024-01-21T08:30:00Z'
  }
])

const newTable = ref({
  name: '',
  description: '',
  columns: []
})

const columnTypes = [
  'VARCHAR',
  'INTEGER',
  'BIGINT',
  'DECIMAL',
  'BOOLEAN',
  'DATE',
  'TIMESTAMP',
  'TEXT',
  'JSON',
  'UUID'
]

const filteredTables = computed(() => {
  if (!searchQuery.value) return tables.value
  return tables.value.filter(table => 
    table.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    (table.description && table.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )
})

const handleCreateTable = async () => {
  if (!newTable.value.name.trim()) return
  
  try {
    // TODO: Call API to create table
    const table = {
      id: Date.now().toString(),
      name: newTable.value.name,
      description: newTable.value.description,
      columns: newTable.value.columns.length,
      rows: 0,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    
    tables.value.push(table)
    isCreateDialogOpen.value = false
    newTable.value = { name: '', description: '', columns: [] }
    
    // Navigate to table builder
    router.push(`/projects/${projectId.value}/tables/${table.id}/builder`)
  } catch (error) {
    console.error('Failed to create table:', error)
  }
}

const handleEditTable = (table: any) => {
  selectedTable.value = table
  isEditDialogOpen.value = true
}

const handleDeleteTable = async (tableId: string) => {
  if (confirm('Are you sure you want to delete this table? This action cannot be undone.')) {
    try {
      // TODO: Call API to delete table
      tables.value = tables.value.filter(t => t.id !== tableId)
    } catch (error) {
      console.error('Failed to delete table:', error)
    }
  }
}

const handleViewTable = (tableId: string) => {
  router.push(`/projects/${projectId.value}/tables/${tableId}/data`)
}

const handleEditSchema = (tableId: string) => {
  router.push(`/projects/${projectId.value}/tables/${tableId}/builder`)
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const addColumn = () => {
  newTable.value.columns.push({
    name: '',
    type: 'VARCHAR',
    nullable: true,
    primary_key: false,
    default_value: ''
  })
}

const removeColumn = (index: number) => {
  newTable.value.columns.splice(index, 1)
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Tables</h1>
        <p class="text-gray-600">Manage your database tables and schemas</p>
      </div>
      <Dialog v-model:open="isCreateDialogOpen">
        <DialogTrigger asChild>
          <Button class="flex items-center space-x-2">
            <Plus class="w-4 h-4" />
            <span>New Table</span>
          </Button>
        </DialogTrigger>
        <DialogContent class="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Create New Table</DialogTitle>
            <DialogDescription>
              Create a new table with custom columns and schema.
            </DialogDescription>
          </DialogHeader>
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label for="table-name">Table Name</Label>
                <Input
                  id="table-name"
                  v-model="newTable.name"
                  placeholder="e.g., users, products"
                  required
                />
              </div>
              <div class="space-y-2">
                <Label for="table-description">Description</Label>
                <Input
                  id="table-description"
                  v-model="newTable.description"
                  placeholder="Brief description"
                />
              </div>
            </div>
            
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <Label>Columns</Label>
                <Button type="button" variant="outline" size="sm" @click="addColumn">
                  <Plus class="w-4 h-4 mr-1" />
                  Add Column
                </Button>
              </div>
              
              <div class="space-y-3 max-h-60 overflow-y-auto">
                <div 
                  v-for="(column, index) in newTable.columns" 
                  :key="index"
                  class="flex items-center space-x-2 p-3 border rounded-lg"
                >
                  <div class="flex-1 grid grid-cols-4 gap-2">
                    <Input
                      v-model="column.name"
                      placeholder="Column name"
                    />
                    <Select v-model="column.type">
                      <SelectTrigger>
                        <SelectValue placeholder="Type" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="type in columnTypes" :key="type" :value="type">
                          {{ type }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <Input
                      v-model="column.default_value"
                      placeholder="Default value"
                    />
                    <div class="flex items-center space-x-2">
                      <input
                        v-model="column.nullable"
                        type="checkbox"
                        class="rounded"
                      />
                      <span class="text-sm">Nullable</span>
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    @click="removeColumn(index)"
                  >
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" @click="isCreateDialogOpen = false">
              Cancel
            </Button>
            <Button @click="handleCreateTable" :disabled="!newTable.name.trim()">
              Create Table
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>

    <!-- Search -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
      <Input
        v-model="searchQuery"
        placeholder="Search tables..."
        class="pl-10"
      />
    </div>

    <!-- Tables Grid -->
    <div v-if="filteredTables.length === 0" class="text-center py-12">
      <Table class="mx-auto h-12 w-12 text-gray-400" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">
        {{ searchQuery ? 'No tables found' : 'No tables yet' }}
      </h3>
      <p class="mt-1 text-sm text-gray-500">
        {{ searchQuery ? 'Try adjusting your search terms.' : 'Get started by creating your first table.' }}
      </p>
      <div v-if="!searchQuery" class="mt-6">
        <Button @click="isCreateDialogOpen = true">Create Table</Button>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card 
        v-for="table in filteredTables" 
        :key="table.id"
        class="hover:shadow-lg transition-shadow duration-200"
      >
        <CardHeader class="pb-3">
          <div class="flex items-start justify-between">
            <div class="flex items-center space-x-3">
              <div class="p-2 bg-blue-100 rounded-lg">
                <Table class="h-5 w-5 text-blue-600" />
              </div>
              <div>
                <CardTitle class="text-lg">{{ table.name }}</CardTitle>
                <CardDescription class="mt-1">
                  {{ table.description || 'No description' }}
                </CardDescription>
              </div>
            </div>
            <div class="flex items-center space-x-1">
              <Button
                variant="ghost"
                size="sm"
                @click="handleViewTable(table.id)"
                title="View Data"
              >
                <Eye class="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="handleEditSchema(table.id)"
                title="Edit Schema"
              >
                <Settings class="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="handleDeleteTable(table.id)"
                title="Delete Table"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent class="pt-0">
          <div class="space-y-3">
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Columns</span>
              <span class="font-medium">{{ table.columns }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Rows</span>
              <span class="font-medium">{{ table.rows.toLocaleString() }}</span>
            </div>
            <div class="flex items-center justify-between text-sm text-gray-500">
              <div class="flex items-center space-x-1">
                <Calendar class="h-4 w-4" />
                <span>Updated {{ formatDate(table.updated_at) }}</span>
              </div>
              <Badge variant="secondary">Active</Badge>
            </div>
          </div>
          <div class="mt-4 flex space-x-2">
            <Button 
              variant="outline" 
              size="sm" 
              class="flex-1"
              @click="handleViewTable(table.id)"
            >
              <Eye class="w-4 h-4 mr-1" />
              View Data
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              class="flex-1"
              @click="handleEditSchema(table.id)"
            >
              <Settings class="w-4 h-4 mr-1" />
              Edit Schema
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
