import { computed, type ComputedRef } from 'vue'
import { useRoute } from 'vue-router'

interface RouteContext {
  projectId: ComputedRef<string>
  tableId: ComputedRef<string | undefined>
  workflowId: ComputedRef<string | undefined>
  hasProject: ComputedRef<boolean>
  hasTable: ComputedRef<boolean>
  hasWorkflow: ComputedRef<boolean>
}

export function useRouteContext(): RouteContext {
  const route = useRoute()

  const projectId = computed(() => (route.params.projectId as string) || '')
  const tableId = computed(() => route.params.tableId as string | undefined)
  const workflowId = computed(() => route.params.workflowId as string | undefined)

  const hasProject = computed(() => !!projectId.value)
  const hasTable = computed(() => !!tableId.value)
  const hasWorkflow = computed(() => !!workflowId.value)

  return {
    projectId,
    tableId,
    workflowId,
    hasProject,
    hasTable,
    hasWorkflow,
  }
}
