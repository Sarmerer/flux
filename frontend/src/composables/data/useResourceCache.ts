import { computed, onUnmounted, ref } from 'vue'

import { useCacheStore } from '@/stores/cache'

import { useWebSocket } from '../external/useWebSocket'

interface CacheOptions<T> {
  fetchFn: () => Promise<T>
  subscribeToUpdates?: boolean
  resourceId?: string
  events?: string[]
  ttlMs?: number
  fetchOnMount?: boolean
  tags?: string[]
}

export function useResourceCache<T>(key: string, options: CacheOptions<T>) {
  const {
    fetchFn,
    subscribeToUpdates = false,
    resourceId,
    events,
    ttlMs = 60000,
    fetchOnMount = true,
    tags,
  } = options

  const cacheStore = useCacheStore()
  const data = ref<T | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  let unsubscribe: (() => void) | null = null

  const syncFromStore = () => {
    const cached = cacheStore.get<T>(key)
    if (cached) {
      data.value = cached.data
      error.value = cached.error
    }
  }

  const fetch = async (force = false) => {
    if (!force && cacheStore.isValid(key, ttlMs)) {
      syncFromStore()
      return
    }

    loading.value = true
    error.value = null

    try {
      const result = await fetchFn()
      data.value = result
      error.value = null
      cacheStore.set<T>(key, result, null, { tags })
    } catch (err) {
      error.value = err as Error
      data.value = null
      cacheStore.set<T>(key, null, err as Error, { tags })
      console.error(`[Cache] Failed to fetch ${key}:`, err)
    } finally {
      loading.value = false
    }
  }

  const refresh = () => fetch(true)

  const invalidate = () => cacheStore.invalidate(key)

  const handleRealtimeUpdate = (message: any) => {
    try {
      const { type, data: updateData } = message

      switch (type) {
        case 'created':
          if (Array.isArray(data.value)) {
            data.value = [...data.value, updateData] as any
            cacheStore.set<T>(key, data.value, null, { tags })
          }
          break

        case 'updated':
          if (Array.isArray(data.value)) {
            const index = (data.value as any[]).findIndex((item) => item?.id === updateData?.id)
            if (index !== -1) {
              data.value = [
                ...(data.value as any[]).slice(0, index),
                updateData,
                ...(data.value as any[]).slice(index + 1),
              ] as any
              cacheStore.set<T>(key, data.value, null, { tags })
            }
          } else if ((data.value as any)?.id === updateData?.id) {
            data.value = updateData
            cacheStore.set<T>(key, data.value, null, { tags })
          }
          break

        case 'deleted':
          if (Array.isArray(data.value)) {
            data.value = (data.value as any[]).filter((item) => item?.id !== updateData?.id) as any
            cacheStore.set<T>(key, data.value, null, { tags })
          } else if ((data.value as any)?.id === updateData?.id) {
            data.value = null
            cacheStore.set<T>(key, null, null, { tags })
          }
          break

        case 'bulk_update':
          refresh()
          break

        default:
          refresh()
      }
    } catch (err) {
      console.error('[Cache] Failed to handle realtime update:', err)
      refresh()
    }
  }

  if (subscribeToUpdates && events?.length) {
    try {
      const { subscribe } = useWebSocket()
      unsubscribe = subscribe({
        resourceId,
        events,
        onMessage: handleRealtimeUpdate,
        throttle: 500,
      })
    } catch (err) {
      console.warn('[Cache] WebSocket not available, updates will be manual only')
    }
  }

  if (fetchOnMount) {
    fetch()
  } else {
    syncFromStore()
  }

  onUnmounted(() => {
    if (unsubscribe) {
      unsubscribe()
    }
  })

  return {
    data: computed(() => data.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    refresh,
    fetch,
    invalidate,
  }
}
