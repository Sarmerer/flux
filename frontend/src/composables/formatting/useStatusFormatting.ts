import { AlertCircle, CheckCircle, Clock, XCircle } from 'lucide-vue-next'
import type { Component } from 'vue'

export function useStatusFormatting() {
  const getStatusIcon = (isActive: boolean): Component => {
    return isActive ? CheckCircle : XCircle
  }

  const getStatusColor = (isActive: boolean): string => {
    return isActive ? 'text-green-600' : 'text-gray-400'
  }

  const getStatusBadgeVariant = (
    status: string
  ): 'default' | 'secondary' | 'destructive' | 'outline' => {
    const statusMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
      completed: 'default',
      success: 'default',
      active: 'default',
      pending: 'secondary',
      in_progress: 'secondary',
      failed: 'destructive',
      error: 'destructive',
      inactive: 'outline',
      cancelled: 'outline',
    }

    return statusMap[status.toLowerCase()] || 'outline'
  }

  const getMutationStatusIcon = (
    status: 'pending' | 'in_progress' | 'completed' | 'failed'
  ): Component => {
    const iconMap: Record<string, Component> = {
      pending: Clock,
      in_progress: Clock,
      completed: CheckCircle,
      failed: AlertCircle,
    }

    return iconMap[status] || Clock
  }

  const getMutationStatusColor = (
    status: 'pending' | 'in_progress' | 'completed' | 'failed'
  ): string => {
    const colorMap: Record<string, string> = {
      pending: 'text-yellow-600',
      in_progress: 'text-blue-600',
      completed: 'text-green-600',
      failed: 'text-red-600',
    }

    return colorMap[status] || 'text-gray-600'
  }

  const getMutationStatusBgColor = (
    status: 'pending' | 'in_progress' | 'completed' | 'failed'
  ): string => {
    const colorMap: Record<string, string> = {
      pending: 'bg-yellow-50',
      in_progress: 'bg-blue-50',
      completed: 'bg-green-50',
      failed: 'bg-red-50',
    }

    return colorMap[status] || 'bg-gray-50'
  }

  return {
    getStatusIcon,
    getStatusColor,
    getStatusBadgeVariant,
    getMutationStatusIcon,
    getMutationStatusColor,
    getMutationStatusBgColor,
  }
}
