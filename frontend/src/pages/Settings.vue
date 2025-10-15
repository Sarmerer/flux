<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  User, 
  Bell, 
  Shield, 
  Database,
  Globe,
  Key,
  Save,
  Eye,
  EyeOff
} from 'lucide-vue-next'

const authStore = useAuthStore()

const user = ref({
  name: authStore.user?.name || '',
  email: authStore.user?.email || '',
  avatar: ''
})

const preferences = ref({
  theme: 'light',
  language: 'en',
  timezone: 'UTC',
  notifications: {
    email: true,
    push: true,
    workflow: true,
    security: true
  }
})

const security = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
  twoFactorEnabled: false,
  apiKey: 'sk-1234567890abcdef',
  showApiKey: false
})

const database = ref({
  defaultHost: 'localhost',
  defaultPort: 5432,
  connectionTimeout: 30,
  maxConnections: 10,
  sslEnabled: true
})

const isSaving = ref(false)
const showApiKey = ref(false)

const themes = [
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' },
  { value: 'system', label: 'System' }
]

const languages = [
  { value: 'en', label: 'English' },
  { value: 'es', label: 'Spanish' },
  { value: 'fr', label: 'French' },
  { value: 'de', label: 'German' }
]

const timezones = [
  { value: 'UTC', label: 'UTC' },
  { value: 'America/New_York', label: 'Eastern Time' },
  { value: 'America/Chicago', label: 'Central Time' },
  { value: 'America/Denver', label: 'Mountain Time' },
  { value: 'America/Los_Angeles', label: 'Pacific Time' },
  { value: 'Europe/London', label: 'London' },
  { value: 'Europe/Paris', label: 'Paris' },
  { value: 'Asia/Tokyo', label: 'Tokyo' }
]

const saveProfile = async () => {
  isSaving.value = true
  try {
    // TODO: Call API to update profile
    console.log('Saving profile:', user.value)
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('Profile updated successfully!')
  } catch (error) {
    console.error('Failed to update profile:', error)
    alert('Failed to update profile')
  } finally {
    isSaving.value = false
  }
}

const savePreferences = async () => {
  isSaving.value = true
  try {
    // TODO: Call API to update preferences
    console.log('Saving preferences:', preferences.value)
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('Preferences updated successfully!')
  } catch (error) {
    console.error('Failed to update preferences:', error)
    alert('Failed to update preferences')
  } finally {
    isSaving.value = false
  }
}

const saveSecurity = async () => {
  isSaving.value = true
  try {
    // TODO: Call API to update security settings
    console.log('Saving security settings:', security.value)
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('Security settings updated successfully!')
    security.value.currentPassword = ''
    security.value.newPassword = ''
    security.value.confirmPassword = ''
  } catch (error) {
    console.error('Failed to update security settings:', error)
    alert('Failed to update security settings')
  } finally {
    isSaving.value = false
  }
}

const saveDatabase = async () => {
  isSaving.value = true
  try {
    // TODO: Call API to update database settings
    console.log('Saving database settings:', database.value)
    await new Promise(resolve => setTimeout(resolve, 1000))
    alert('Database settings updated successfully!')
  } catch (error) {
    console.error('Failed to update database settings:', error)
    alert('Failed to update database settings')
  } finally {
    isSaving.value = false
  }
}

const generateNewApiKey = async () => {
  if (confirm('Are you sure you want to generate a new API key? The current key will be invalidated.')) {
    try {
      // TODO: Call API to generate new key
      security.value.apiKey = 'sk-' + Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
      alert('New API key generated successfully!')
    } catch (error) {
      console.error('Failed to generate API key:', error)
      alert('Failed to generate API key')
    }
  }
}

const toggleTwoFactor = async () => {
  try {
    // TODO: Call API to toggle 2FA
    security.value.twoFactorEnabled = !security.value.twoFactorEnabled
    alert(security.value.twoFactorEnabled ? 'Two-factor authentication enabled!' : 'Two-factor authentication disabled!')
  } catch (error) {
    console.error('Failed to toggle 2FA:', error)
    alert('Failed to update two-factor authentication')
  }
}
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold text-gray-900">Settings</h1>
      <p class="text-gray-600">Manage your account settings and preferences</p>
    </div>

    <Tabs default-value="profile" class="space-y-6">
      <TabsList>
        <TabsTrigger value="profile">Profile</TabsTrigger>
        <TabsTrigger value="preferences">Preferences</TabsTrigger>
        <TabsTrigger value="security">Security</TabsTrigger>
        <TabsTrigger value="database">Database</TabsTrigger>
      </TabsList>

      <!-- Profile Tab -->
      <TabsContent value="profile">
        <Card>
          <CardHeader>
            <CardTitle class="flex items-center space-x-2">
              <User class="w-5 h-5" />
              <span>Profile Information</span>
            </CardTitle>
            <CardDescription>
              Update your personal information and profile details
            </CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label for="name">Full Name</Label>
                <Input
                  id="name"
                  v-model="user.name"
                  placeholder="Enter your full name"
                />
              </div>
              <div class="space-y-2">
                <Label for="email">Email Address</Label>
                <Input
                  id="email"
                  v-model="user.email"
                  type="email"
                  placeholder="Enter your email"
                />
              </div>
            </div>
            <div class="space-y-2">
              <Label for="avatar">Avatar URL</Label>
              <Input
                id="avatar"
                v-model="user.avatar"
                placeholder="https://example.com/avatar.jpg"
              />
            </div>
            <div class="flex justify-end">
              <Button @click="saveProfile" :disabled="isSaving">
                <Save class="w-4 h-4 mr-2" />
                {{ isSaving ? 'Saving...' : 'Save Profile' }}
              </Button>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <!-- Preferences Tab -->
      <TabsContent value="preferences">
        <div class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle class="flex items-center space-x-2">
                <Globe class="w-5 h-5" />
                <span>General Preferences</span>
              </CardTitle>
              <CardDescription>
                Configure your general application preferences
              </CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label for="theme">Theme</Label>
                  <Select v-model="preferences.theme">
                    <SelectTrigger>
                      <SelectValue placeholder="Select theme" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="theme in themes" :key="theme.value" :value="theme.value">
                        {{ theme.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div class="space-y-2">
                  <Label for="language">Language</Label>
                  <Select v-model="preferences.language">
                    <SelectTrigger>
                      <SelectValue placeholder="Select language" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="lang in languages" :key="lang.value" :value="lang.value">
                        {{ lang.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div class="space-y-2">
                <Label for="timezone">Timezone</Label>
                <Select v-model="preferences.timezone">
                  <SelectTrigger>
                    <SelectValue placeholder="Select timezone" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem v-for="tz in timezones" :key="tz.value" :value="tz.value">
                      {{ tz.label }}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle class="flex items-center space-x-2">
                <Bell class="w-5 h-5" />
                <span>Notifications</span>
              </CardTitle>
              <CardDescription>
                Configure your notification preferences
              </CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="flex items-center justify-between">
                <div>
                  <Label>Email Notifications</Label>
                  <p class="text-sm text-gray-500">Receive notifications via email</p>
                </div>
                <Switch v-model="preferences.notifications.email" />
              </div>
              <div class="flex items-center justify-between">
                <div>
                  <Label>Push Notifications</Label>
                  <p class="text-sm text-gray-500">Receive push notifications in browser</p>
                </div>
                <Switch v-model="preferences.notifications.push" />
              </div>
              <div class="flex items-center justify-between">
                <div>
                  <Label>Workflow Notifications</Label>
                  <p class="text-sm text-gray-500">Get notified about workflow executions</p>
                </div>
                <Switch v-model="preferences.notifications.workflow" />
              </div>
              <div class="flex items-center justify-between">
                <div>
                  <Label>Security Notifications</Label>
                  <p class="text-sm text-gray-500">Get notified about security events</p>
                </div>
                <Switch v-model="preferences.notifications.security" />
              </div>
            </CardContent>
          </Card>

          <div class="flex justify-end">
            <Button @click="savePreferences" :disabled="isSaving">
              <Save class="w-4 h-4 mr-2" />
              {{ isSaving ? 'Saving...' : 'Save Preferences' }}
            </Button>
          </div>
        </div>
      </TabsContent>

      <!-- Security Tab -->
      <TabsContent value="security">
        <div class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle class="flex items-center space-x-2">
                <Shield class="w-5 h-5" />
                <span>Password & Security</span>
              </CardTitle>
              <CardDescription>
                Update your password and security settings
              </CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="space-y-2">
                <Label for="current-password">Current Password</Label>
                <Input
                  id="current-password"
                  v-model="security.currentPassword"
                  type="password"
                  placeholder="Enter current password"
                />
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label for="new-password">New Password</Label>
                  <Input
                    id="new-password"
                    v-model="security.newPassword"
                    type="password"
                    placeholder="Enter new password"
                  />
                </div>
                <div class="space-y-2">
                  <Label for="confirm-password">Confirm Password</Label>
                  <Input
                    id="confirm-password"
                    v-model="security.confirmPassword"
                    type="password"
                    placeholder="Confirm new password"
                  />
                </div>
              </div>
              <div class="flex items-center justify-between">
                <div>
                  <Label>Two-Factor Authentication</Label>
                  <p class="text-sm text-gray-500">Add an extra layer of security to your account</p>
                </div>
                <Switch v-model="security.twoFactorEnabled" @click="toggleTwoFactor" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle class="flex items-center space-x-2">
                <Key class="w-5 h-5" />
                <span>API Access</span>
              </CardTitle>
              <CardDescription>
                Manage your API keys for programmatic access
              </CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="space-y-2">
                <Label>API Key</Label>
                <div class="flex items-center space-x-2">
                  <Input
                    v-model="security.apiKey"
                    :type="showApiKey ? 'text' : 'password'"
                    readonly
                    class="font-mono"
                  />
                  <Button
                    variant="outline"
                    size="sm"
                    @click="showApiKey = !showApiKey"
                  >
                    <Eye v-if="!showApiKey" class="w-4 h-4" />
                    <EyeOff v-else class="w-4 h-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    @click="generateNewApiKey"
                  >
                    Generate New
                  </Button>
                </div>
                <p class="text-sm text-gray-500">
                  Keep your API key secure and never share it publicly
                </p>
              </div>
            </CardContent>
          </Card>

          <div class="flex justify-end">
            <Button @click="saveSecurity" :disabled="isSaving">
              <Save class="w-4 h-4 mr-2" />
              {{ isSaving ? 'Saving...' : 'Save Security Settings' }}
            </Button>
          </div>
        </div>
      </TabsContent>

      <!-- Database Tab -->
      <TabsContent value="database">
        <Card>
          <CardHeader>
            <CardTitle class="flex items-center space-x-2">
              <Database class="w-5 h-5" />
              <span>Database Settings</span>
            </CardTitle>
            <CardDescription>
              Configure default database connection settings
            </CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label for="default-host">Default Host</Label>
                <Input
                  id="default-host"
                  v-model="database.defaultHost"
                  placeholder="localhost"
                />
              </div>
              <div class="space-y-2">
                <Label for="default-port">Default Port</Label>
                <Input
                  id="default-port"
                  v-model="database.defaultPort"
                  type="number"
                  placeholder="5432"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label for="connection-timeout">Connection Timeout (seconds)</Label>
                <Input
                  id="connection-timeout"
                  v-model="database.connectionTimeout"
                  type="number"
                  placeholder="30"
                />
              </div>
              <div class="space-y-2">
                <Label for="max-connections">Max Connections</Label>
                <Input
                  id="max-connections"
                  v-model="database.maxConnections"
                  type="number"
                  placeholder="10"
                />
              </div>
            </div>
            <div class="flex items-center justify-between">
              <div>
                <Label>SSL Enabled</Label>
                <p class="text-sm text-gray-500">Use SSL for database connections</p>
              </div>
              <Switch v-model="database.sslEnabled" />
            </div>
            <div class="flex justify-end">
              <Button @click="saveDatabase" :disabled="isSaving">
                <Save class="w-4 h-4 mr-2" />
                {{ isSaving ? 'Saving...' : 'Save Database Settings' }}
              </Button>
            </div>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
