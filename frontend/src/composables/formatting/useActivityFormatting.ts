import {
  CheckCircle,
  Database,
  Edit,
  Play,
  Plus,
  Table,
  Trash,
  Workflow,
  XCircle,
} from 'lucide-vue-next'
import type { Component } from 'vue'

import type { ActivityLog } from '@/types'

export function useActivityFormatting() {
  const getActivityIcon = (activity: ActivityLog): Component => {
    switch (activity.type) {
      case 'table_created':
        return Plus
      case 'table_updated':
        return Edit
      case 'table_deleted':
        return Trash
      case 'row_created':
        return Plus
      case 'row_updated':
        return Edit
      case 'row_deleted':
        return Trash
      case 'workflow_triggered':
        return Play
      case 'workflow_completed':
        return CheckCircle
      case 'workflow_failed':
        return XCircle
      default:
        return Table
    }
  }

  const getEntityIcon = (entityType: string): Component => {
    const iconMap: Record<string, Component> = {
      table: Table,
      row: Table,
      workflow: Workflow,
      database: Database,
    }

    return iconMap[entityType] || Table
  }

  const getActivityColor = (type: ActivityLog['type']): string => {
    const colorMap: Record<ActivityLog['type'], string> = {
      table_created: 'text-green-600',
      table_updated: 'text-blue-600',
      table_deleted: 'text-red-600',
      row_created: 'text-green-600',
      row_updated: 'text-blue-600',
      row_deleted: 'text-red-600',
      workflow_triggered: 'text-purple-600',
      workflow_completed: 'text-green-600',
      workflow_failed: 'text-red-600',
    }

    return colorMap[type] || 'text-gray-600'
  }

  const getActivityBgColor = (type: ActivityLog['type']): string => {
    const colorMap: Record<ActivityLog['type'], string> = {
      table_created: 'bg-green-50',
      table_updated: 'bg-blue-50',
      table_deleted: 'bg-red-50',
      row_created: 'bg-green-50',
      row_updated: 'bg-blue-50',
      row_deleted: 'bg-red-50',
      workflow_triggered: 'bg-purple-50',
      workflow_completed: 'bg-green-50',
      workflow_failed: 'bg-red-50',
    }

    return colorMap[type] || 'bg-gray-50'
  }

  const formatActivityType = (type: ActivityLog['type']): string => {
    const typeMap: Record<ActivityLog['type'], string> = {
      table_created: 'Table Created',
      table_updated: 'Table Updated',
      table_deleted: 'Table Deleted',
      row_created: 'Row Created',
      row_updated: 'Row Updated',
      row_deleted: 'Row Deleted',
      workflow_triggered: 'Workflow Triggered',
      workflow_completed: 'Workflow Completed',
      workflow_failed: 'Workflow Failed',
    }

    return typeMap[type] || type
  }

  const formatEntityType = (entityType: string): string => {
    return entityType.charAt(0).toUpperCase() + entityType.slice(1)
  }

  return {
    getActivityIcon,
    getEntityIcon,
    getActivityColor,
    getActivityBgColor,
    formatActivityType,
    formatEntityType,
  }
}
