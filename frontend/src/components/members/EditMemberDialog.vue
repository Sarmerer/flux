<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

import { useProjectMemberStore } from '@/stores/projectMember'
import type { Permission, ProjectMemberWithUser, Role } from '@/types/auth'

import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

interface Props {
  open: boolean
  projectId: string
  member: ProjectMemberWithUser | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  'member-updated': []
}>()

const projectMemberStore = useProjectMemberStore()

const selectedRole = ref<Role>('viewer')
const selectedPermissions = ref<Permission[]>([])
const isSubmitting = ref(false)
const error = ref<string | null>(null)

const permissionGroups = [
  {
    name: 'Projects',
    permissions: [
      { label: 'View', value: 'projects.view' as Permission },
      { label: 'Edit', value: 'projects.edit' as Permission },
      { label: 'Delete', value: 'projects.delete' as Permission },
    ],
  },
  {
    name: 'Tables',
    permissions: [
      { label: 'Create', value: 'tables.create' as Permission },
      { label: 'Edit', value: 'tables.edit' as Permission },
      { label: 'Delete', value: 'tables.delete' as Permission },
    ],
  },
  {
    name: 'Workflows',
    permissions: [
      { label: 'Create', value: 'workflows.create' as Permission },
      { label: 'Edit', value: 'workflows.edit' as Permission },
      { label: 'Delete', value: 'workflows.delete' as Permission },
    ],
  },
  {
    name: 'Team',
    permissions: [
      { label: 'Manage Members', value: 'users.manage' as Permission },
      { label: 'Manage Settings', value: 'settings.manage' as Permission },
    ],
  },
]

const hasChanges = computed(() => {
  if (!props.member) return false

  const roleChanged = selectedRole.value !== props.member.role
  const permissionsChanged =
    JSON.stringify([...selectedPermissions.value].sort()) !==
    JSON.stringify([...props.member.permissions].sort())

  return roleChanged || permissionsChanged
})

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.member) {
      selectedRole.value = props.member.role
      selectedPermissions.value = [...props.member.permissions]
      error.value = null
    }
  }
)

const togglePermission = (permission: Permission) => {
  const index = selectedPermissions.value.indexOf(permission)
  if (index > -1) {
    selectedPermissions.value.splice(index, 1)
  } else {
    selectedPermissions.value.push(permission)
  }
}

const getInitials = (name: string): string => {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const handleSubmit = async () => {
  if (!props.member || !hasChanges.value) return

  isSubmitting.value = true
  error.value = null

  try {
    await projectMemberStore.updateProjectMember(
      props.projectId,
      props.member.id,
      selectedRole.value,
      selectedPermissions.value
    )
    emit('member-updated')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to update member'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="$emit('update:open', $event)">
    <DialogContent class="sm:max-w-[600px]" v-if="member">
      <DialogHeader>
        <DialogTitle>Edit Member Role</DialogTitle>
        <DialogDescription>
          Update role and permissions for {{ member.user.name }}
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="rounded-lg border p-4 bg-muted/50">
          <div class="flex items-center gap-3">
            <Avatar>
              <AvatarFallback>{{ getInitials(member.user.name) }}</AvatarFallback>
            </Avatar>
            <div>
              <div class="font-medium">{{ member.user.name }}</div>
              <div class="text-sm text-muted-foreground">{{ member.user.email }}</div>
            </div>
          </div>
        </div>

        <div class="space-y-2">
          <Label for="edit-role">Role</Label>
          <Select v-model="selectedRole">
            <SelectTrigger id="edit-role">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="admin">
                <div class="flex flex-col items-start">
                  <div class="font-medium">Admin</div>
                  <div class="text-xs text-muted-foreground">Full access to all features</div>
                </div>
              </SelectItem>
              <SelectItem value="manager">
                <div class="flex flex-col items-start">
                  <div class="font-medium">Manager</div>
                  <div class="text-xs text-muted-foreground">Can manage workflows and tables</div>
                </div>
              </SelectItem>
              <SelectItem value="developer">
                <div class="flex flex-col items-start">
                  <div class="font-medium">Developer</div>
                  <div class="text-xs text-muted-foreground">Can create and edit resources</div>
                </div>
              </SelectItem>
              <SelectItem value="viewer">
                <div class="flex flex-col items-start">
                  <div class="font-medium">Viewer</div>
                  <div class="text-xs text-muted-foreground">Read-only access</div>
                </div>
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="space-y-3" v-if="selectedRole !== 'admin'">
          <Label>Permissions</Label>
          <div class="space-y-2 rounded-lg border p-4">
            <div v-for="group in permissionGroups" :key="group.name" class="space-y-2">
              <div class="font-medium text-sm">{{ group.name }}</div>
              <div class="grid grid-cols-2 gap-2 pl-4">
                <div
                  v-for="permission in group.permissions"
                  :key="permission.value"
                  class="flex items-center space-x-2"
                >
                  <Checkbox
                    :id="`edit-${permission.value}`"
                    :checked="selectedPermissions.includes(permission.value)"
                    @update:checked="togglePermission(permission.value)"
                  />
                  <Label
                    :for="`edit-${permission.value}`"
                    class="text-sm font-normal cursor-pointer"
                  >
                    {{ permission.label }}
                  </Label>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="error" class="rounded-lg border border-destructive bg-destructive/10 p-3">
          <p class="text-sm text-destructive">{{ error }}</p>
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            @click="$emit('update:open', false)"
            :disabled="isSubmitting"
          >
            Cancel
          </Button>
          <Button type="submit" :disabled="isSubmitting || !hasChanges">
            <Loader2 v-if="isSubmitting" class="mr-2 h-4 w-4 animate-spin" />
            Save Changes
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
