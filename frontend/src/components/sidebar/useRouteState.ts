import { useRoute } from 'vue-router'

export const useRouteState = () => {
  const route = useRoute()

  const isActive = (itemUrl: string) => {
    return route.path === itemUrl || route.path.startsWith(itemUrl + '/')
  }

  const isExactActive = (itemUrl: string) => {
    return route.path === itemUrl
  }

  return {
    isActive,
    isExactActive,
  }
}
