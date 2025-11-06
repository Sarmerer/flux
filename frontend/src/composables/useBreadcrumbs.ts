import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useActiveProjectStore } from '@/stores/activeProject'
import { ROUTE_SEGMENTS } from '@/constants/routes'

export interface BreadcrumbItem {
  label: string
  path?: string
  isProjectSelector?: boolean
}

export function useBreadcrumbs() {
  const route = useRoute()
  const activeProjectStore = useActiveProjectStore()
  const activeProject = computed(() => activeProjectStore.activeProject)

  const breadcrumbs = computed<BreadcrumbItem[]>(() => {
    const crumbs: BreadcrumbItem[] = []
    const path = route.path

    if (path === '/') {
      crumbs.push({ label: 'Home' })
      return crumbs
    }

    if (path === '/projects') {
      crumbs.push({ label: 'Projects' })
      return crumbs
    }

    if (path === '/settings') {
      crumbs.push({ label: 'Settings' })
      return crumbs
    }

    if (activeProject.value) {
      crumbs.push({
        label: activeProject.value.name,
        path: `/projects/${activeProject.value.id}`,
        isProjectSelector: true,
      })

      const segment = getProjectRouteSegment(path)
      if (segment) {
        crumbs.push(buildSegmentCrumb(segment, activeProject.value.id))
      }
    }

    return crumbs
  })

  return {
    breadcrumbs,
  }
}

function getProjectRouteSegment(path: string): string | null {
  if (path.includes(`/${ROUTE_SEGMENTS.TABLES}`)) return ROUTE_SEGMENTS.TABLES
  if (path.includes(`/${ROUTE_SEGMENTS.WORKFLOWS}`)) return ROUTE_SEGMENTS.WORKFLOWS
  if (path.includes(`/${ROUTE_SEGMENTS.MEMBERS}`)) return ROUTE_SEGMENTS.MEMBERS
  if (path.includes(`/${ROUTE_SEGMENTS.ACTIVITY}`)) return ROUTE_SEGMENTS.ACTIVITY
  return null
}

function buildSegmentCrumb(segment: string, projectId: string): BreadcrumbItem {
  const labels: Record<string, string> = {
    [ROUTE_SEGMENTS.TABLES]: 'Tables',
    [ROUTE_SEGMENTS.WORKFLOWS]: 'Workflows',
    [ROUTE_SEGMENTS.MEMBERS]: 'Members',
    [ROUTE_SEGMENTS.ACTIVITY]: 'Activity',
  }

  return {
    label: labels[segment] || segment,
    path: `/projects/${projectId}/${segment}`,
  }
}
