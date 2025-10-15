<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { 
  Table, 
  Plus, 
  Save, 
  ArrowLeft,
  Trash2,
  GripVertical,
  Key,
  Link,
  Eye,
  EyeOff
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const projectId = computed(() => route.params.projectId as string)
const tableId = computed(() => route.params.tableId as string)

const tableName = ref('users')
const tableDescription = ref('User accounts and profiles')
const columns = ref([
  {
    id: '1',
    name: 'id',
    type: 'UUID',
    nullable: false,
    primary_key: true,
    unique: false,
    default_value: '',
    foreign_key: null,
    order: 1
  },
  {
    id: '2',
    name: 'email',
    type: 'VARCHAR',
    nullable: false,
    primary_key: false,
    unique: true,
    default_value: '',
    foreign_key: null,
    order: 2
  },
  {
    id: '3',
    name: 'name',
    type: 'VARCHAR',
    nullable: false,
    primary_key: false,
    unique: false,
    default_value: '',
    foreign_key: null,
    order: 3
  },
  {
    id: '4',
    name: 'created_at',
    type: 'TIMESTAMP',
    nullable: false,
    primary_key: false,
    unique: false,
    default_value: 'CURRENT_TIMESTAMP',
    foreign_key: null,
    order: 4
  }
])

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

const isSaving = ref(false)
const showPreview = ref(true)

const addColumn = () => {
  const newColumn = {
    id: Date.now().toString(),
    name: '',
    type: 'VARCHAR',
    nullable: true,
    primary_key: false,
    unique: false,
    default_value: '',
    foreign_key: null,
    order: columns.value.length + 1
  }
  columns.value.push(newColumn)
}

const removeColumn = (columnId: string) => {
  columns.value = columns.value.filter(col => col.id !== columnId)
  // Reorder remaining columns
  columns.value.forEach((col, index) => {
    col.order = index + 1
  })
}

const moveColumn = (fromIndex: number, toIndex: number) => {
  const column = columns.value.splice(fromIndex, 1)[0]
  columns.value.splice(toIndex, 0, column)
  // Update order
  columns.value.forEach((col, index) => {
    col.order = index + 1
  })
}

const saveTable = async () => {
  isSaving.value = true
  try {
    // TODO: Call API to save table schema
    console.log('Saving table:', {
      name: tableName.value,
      description: tableDescription.value,
      columns: columns.value
    })
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Navigate back to tables list
    router.push(`/projects/${projectId.value}/tables`)
  } catch (error) {
    console.error('Failed to save table:', error)
  } finally {
    isSaving.value = false
  }
}

const generateSQL = () => {
  const primaryKeys = columns.value.filter(col => col.primary_key).map(col => col.name)
  const foreignKeys = columns.value.filter(col => col.foreign_key).map(col => 
    `FOREIGN KEY (${col.name}) REFERENCES ${col.foreign_key.table}(${col.foreign_key.column})`
  )
  
  let sql = `CREATE TABLE ${tableName.value} (\n`
  
  columns.value.forEach((col, index) => {
    let columnDef = `  ${col.name} ${col.type}`
    
    if (!col.nullable) columnDef += ' NOT NULL'
    if (col.unique) columnDef += ' UNIQUE'
    if (col.default_value) columnDef += ` DEFAULT ${col.default_value}`
    if (col.primary_key) columnDef += ' PRIMARY KEY'
    
    sql += columnDef
    if (index < columns.value.length - 1 || foreignKeys.length > 0) {
      sql += ','
    }
    sql += '\n'
  })
  
  if (foreignKeys.length > 0) {
    sql += '  ' + foreignKeys.join(',\n  ') + '\n'
  }
  
  sql += ');'
  
  return sql
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <Button variant="ghost" size="sm" @click="router.push(`/projects/${projectId}/tables`)">
          <ArrowLeft class="w-4 h-4 mr-2" />
          Back to Tables
        </Button>
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Table Builder</h1>
          <p class="text-gray-600">Design your table schema with drag-and-drop</p>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <Button variant="outline" @click="showPreview = !showPreview">
          <Eye v-if="!showPreview" class="w-4 h-4 mr-2" />
          <EyeOff v-else class="w-4 h-4 mr-2" />
          {{ showPreview ? 'Hide' : 'Show' }} Preview
        </Button>
        <Button @click="saveTable" :disabled="isSaving">
          <Save class="w-4 h-4 mr-2" />
          {{ isSaving ? 'Saving...' : 'Save Table' }}
        </Button>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Table Configuration -->
      <Card>
        <CardHeader>
          <CardTitle>Table Configuration</CardTitle>
          <CardDescription>Basic table settings and metadata</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="table-name">Table Name</Label>
            <Input
              id="table-name"
              v-model="tableName"
              placeholder="e.g., users, products"
            />
          </div>
          <div class="space-y-2">
            <Label for="table-description">Description</Label>
            <Input
              id="table-description"
              v-model="tableDescription"
              placeholder="Brief description of the table"
            />
          </div>
        </CardContent>
      </Card>

      <!-- Column Management -->
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>Columns</CardTitle>
              <CardDescription>Define your table columns and their properties</CardDescription>
            </div>
            <Button variant="outline" size="sm" @click="addColumn">
              <Plus class="w-4 h-4 mr-1" />
              Add Column
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div class="space-y-3 max-h-96 overflow-y-auto">
            <div 
              v-for="(column, index) in columns" 
              :key="column.id"
              class="flex items-center space-x-2 p-3 border rounded-lg bg-gray-50"
            >
              <GripVertical class="w-4 h-4 text-gray-400 cursor-move" />
              
              <div class="flex-1 grid grid-cols-6 gap-2">
                <Input
                  v-model="column.name"
                  placeholder="Column name"
                  class="col-span-2"
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
                  placeholder="Default"
                  class="col-span-2"
                />
                <div class="flex items-center space-x-1">
                  <Checkbox
                    v-model="column.nullable"
                    id="nullable"
                  />
                  <Label for="nullable" class="text-xs">Nullable</Label>
                </div>
              </div>
              
              <div class="flex items-center space-x-1">
                <Button
                  variant="ghost"
                  size="sm"
                  :class="{ 'bg-blue-100 text-blue-600': column.primary_key }"
                  @click="column.primary_key = !column.primary_key"
                  title="Primary Key"
                >
                  <Key class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  :class="{ 'bg-green-100 text-green-600': column.unique }"
                  @click="column.unique = !column.unique"
                  title="Unique"
                >
                  <Link class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  @click="removeColumn(column.id)"
                  title="Remove Column"
                >
                  <Trash2 class="w-4 h-4" />
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- SQL Preview -->
    <Card v-if="showPreview">
      <CardHeader>
        <CardTitle>SQL Preview</CardTitle>
        <CardDescription>Generated SQL for your table schema</CardDescription>
      </CardHeader>
      <CardContent>
        <pre class="bg-gray-900 text-green-400 p-4 rounded-lg overflow-x-auto text-sm"><code>{{ generateSQL() }}</code></pre>
      </CardContent>
    </Card>

    <!-- Table Preview -->
    <Card v-if="showPreview">
      <CardHeader>
        <CardTitle>Table Preview</CardTitle>
        <CardDescription>Visual representation of your table structure</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="overflow-x-auto">
          <table class="w-full border-collapse border border-gray-300">
            <thead>
              <tr class="bg-gray-100">
                <th class="border border-gray-300 px-4 py-2 text-left font-medium">Column</th>
                <th class="border border-gray-300 px-4 py-2 text-left font-medium">Type</th>
                <th class="border border-gray-300 px-4 py-2 text-left font-medium">Nullable</th>
                <th class="border border-gray-300 px-4 py-2 text-left font-medium">Default</th>
                <th class="border border-gray-300 px-4 py-2 text-left font-medium">Constraints</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="column in columns" :key="column.id" class="hover:bg-gray-50">
                <td class="border border-gray-300 px-4 py-2">
                  <div class="flex items-center space-x-2">
                    <span class="font-medium">{{ column.name || 'unnamed' }}</span>
                    <div class="flex space-x-1">
                      <Badge v-if="column.primary_key" variant="default" class="text-xs">PK</Badge>
                      <Badge v-if="column.unique" variant="secondary" class="text-xs">UQ</Badge>
                    </div>
                  </div>
                </td>
                <td class="border border-gray-300 px-4 py-2 text-sm text-gray-600">{{ column.type }}</td>
                <td class="border border-gray-300 px-4 py-2 text-sm text-gray-600">
                  {{ column.nullable ? 'Yes' : 'No' }}
                </td>
                <td class="border border-gray-300 px-4 py-2 text-sm text-gray-600">
                  {{ column.default_value || '-' }}
                </td>
                <td class="border border-gray-300 px-4 py-2 text-sm text-gray-600">
                  <div class="flex space-x-1">
                    <Badge v-if="column.primary_key" variant="default" class="text-xs">Primary Key</Badge>
                    <Badge v-if="column.unique" variant="secondary" class="text-xs">Unique</Badge>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
