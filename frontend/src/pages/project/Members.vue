<script setup lang="ts">
import AddMemberDialog from '@/components/members/AddMemberDialog.vue'
import EditMemberDialog from '@/components/members/EditMemberDialog.vue'
import { Edit, Loader2, MoreVertical, Trash2, UserPlus, Users } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'

import { useAuthStore } from '@/stores/auth'
import { useProjectMemberStore } from '@/stores/projectMember'
import type { ProjectMemberWithUser } from '@/types'
import { formatDistanceToNow } from 'date-fns'

import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import { useRouteContext } from '@/composables/routing'

const projectMemberStore = useProjectMemberStore()
const authStore = useAuthStore()

const { projectId } = useRouteContext()
const members = computed(() => projectMemberStore.getProjectMembers(projectId.value))
const isLoading = ref(false)
const error = ref<string | null>(null)

const isAddDialogOpen = ref(false)
const isEditDialogOpen = ref(false)
const selectedMember = ref<ProjectMemberWithUser | null>(null)

const canManageMembers = computed(() => {
  return projectMemberStore.hasPermission('users.manage')
})

onMounted(async () => {
  isLoading.value = true
  error.value = null
  try {
    await Promise.all([
      projectMemberStore.loadProjectMembers(projectId.value),
      projectMemberStore.loadMyProjectRole(projectId.value),
    ])
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load members'
  } finally {
    isLoading.value = false
  }
})

const openAddDialog = () => {
  isAddDialogOpen.value = true
}

const openEditDialog = (member: ProjectMemberWithUser) => {
  selectedMember.value = member
  isEditDialogOpen.value = true
}

const removeMember = async (member: ProjectMemberWithUser) => {
  if (member.user.id === authStore.user?.id) {
    return
  }

  if (!confirm(`Are you sure you want to remove ${member.user.name} from this project?`)) {
    return
  }

  try {
    await projectMemberStore.removeProjectMember(projectId.value, member.id)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to remove member'
  }
}

const handleMemberAdded = () => {
  isAddDialogOpen.value = false
}

const handleMemberUpdated = () => {
  isEditDialogOpen.value = false
}

const getInitials = (name: string): string => {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const getRoleVariant = (role: string) => {
  const variants: Record<string, 'default' | 'secondary' | 'outline' | 'destructive'> = {
    admin: 'default',
    manager: 'secondary',
    developer: 'outline',
    viewer: 'outline',
  }
  return variants[role] || 'outline'
}

const getRoleLabel = (role: string): string => {
  const labels: Record<string, string> = {
    admin: 'Admin',
    manager: 'Manager',
    developer: 'Developer',
    viewer: 'Viewer',
  }
  return labels[role] || role
}

const getPermissionLabel = (permission: string): string => {
  return permission.split('.').pop() || permission
}

const formatDate = (date: string): string => {
  return formatDistanceToNow(new Date(date), { addSuffix: true })
}
</script>

<template>
  <div class="container mx-auto py-8 px-4">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold">Team Members</h1>
        <p class="text-muted-foreground mt-1">Manage roles and permissions for your team</p>
      </div>
      <Button @click="openAddDialog" v-if="canManageMembers">
        <UserPlus class="mr-2 h-4 w-4" />
        Add Member
      </Button>
    </div>

    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div v-else-if="error" class="rounded-lg border border-destructive bg-destructive/10 p-4">
      <p class="text-sm text-destructive">{{ error }}</p>
    </div>

    <div v-else-if="members.length === 0" class="rounded-lg border border-dashed p-12 text-center">
      <Users class="mx-auto h-12 w-12 text-muted-foreground mb-4" />
      <h3 class="text-lg font-semibold mb-2">No team members yet</h3>
      <p class="text-sm text-muted-foreground mb-4">Get started by adding your first team member</p>
      <Button @click="openAddDialog" v-if="canManageMembers">
        <UserPlus class="mr-2 h-4 w-4" />
        Add Member
      </Button>
    </div>

    <div v-else class="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Member</TableHead>
            <TableHead>Role</TableHead>
            <TableHead>Permissions</TableHead>
            <TableHead>Joined</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="member in members" :key="member.id">
            <TableCell>
              <div class="flex items-center gap-3">
                <Avatar>
                  <AvatarFallback>{{ getInitials(member.user.name) }}</AvatarFallback>
                </Avatar>
                <div>
                  <div class="font-medium">{{ member.user.name }}</div>
                  <div class="text-sm text-muted-foreground">{{ member.user.email }}</div>
                </div>
              </div>
            </TableCell>
            <TableCell>
              <Badge :variant="getRoleVariant(member.role)">
                {{ getRoleLabel(member.role) }}
              </Badge>
            </TableCell>
            <TableCell>
              <div class="flex flex-wrap gap-1">
                <Badge v-if="member.role === 'admin'" variant="outline" class="text-xs">
                  All Permissions
                </Badge>
                <template v-else>
                  <Badge
                    v-for="permission in member.permissions.slice(0, 2)"
                    :key="permission"
                    variant="outline"
                    class="text-xs"
                  >
                    {{ getPermissionLabel(permission) }}
                  </Badge>
                  <Badge v-if="member.permissions.length > 2" variant="outline" class="text-xs">
                    +{{ member.permissions.length - 2 }} more
                  </Badge>
                </template>
              </div>
            </TableCell>
            <TableCell class="text-muted-foreground">
              {{ formatDate(member.created_at) }}
            </TableCell>
            <TableCell class="text-right">
              <DropdownMenu>
                <DropdownMenuTrigger as-child>
                  <Button variant="ghost" size="icon" :disabled="!canManageMembers">
                    <MoreVertical class="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem @click="openEditDialog(member)">
                    <Edit class="mr-2 h-4 w-4" />
                    Edit Role & Permissions
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    @click="removeMember(member)"
                    class="text-destructive"
                    :disabled="member.user.id === authStore.user?.id"
                  >
                    <Trash2 class="mr-2 h-4 w-4" />
                    Remove Member
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <AddMemberDialog
      v-model:open="isAddDialogOpen"
      :project-id="projectId"
      @member:add="handleMemberAdded"
    />

    <EditMemberDialog
      v-model:open="isEditDialogOpen"
      :project-id="projectId"
      :member="selectedMember"
      @member:update="handleMemberUpdated"
    />
  </div>
</template>
