import { computed, onUnmounted, ref } from 'vue'

import { useWebSocket } from '../external/useWebSocket'

interface CacheOptions<T> {
  fetchFn: () => Promise<T>
  subscribeToUpdates?: boolean
  scope?: 'project' | 'table' | 'workflow' | 'global'
  resourceId?: string
  events?: string[]
  ttl?: number // Time to live in milliseconds
  refetchOnMount?: boolean
  optimisticUpdate?: boolean
}

interface CacheEntry<T> {
  data: T | null
  timestamp: number
  loading: boolean
  error: Error | null
}

// Global cache storage
const globalCache = new Map<string, CacheEntry<any>>()

export function useResourceCache<T>(key: string, options: CacheOptions<T>) {
  const {
    fetchFn,
    subscribeToUpdates = false,
    scope,
    resourceId,
    events,
    ttl = 60000, // Default 1 minute
    refetchOnMount = true,
    optimisticUpdate = false,
  } = options

  const data = ref<T | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)
  const lastFetchTime = ref<number>(0)

  let unsubscribe: (() => void) | null = null

  // Check if cache is still valid
  const isCacheValid = (): boolean => {
    const cached = globalCache.get(key)
    if (!cached) return false

    const age = Date.now() - cached.timestamp
    return age < ttl
  }

  // Load from cache if valid
  const loadFromCache = (): boolean => {
    const cached = globalCache.get(key)
    if (cached && isCacheValid()) {
      data.value = cached.data
      error.value = cached.error
      lastFetchTime.value = cached.timestamp
      return true
    }
    return false
  }

  // Save to cache
  const saveToCache = (newData: T | null, newError: Error | null = null) => {
    globalCache.set(key, {
      data: newData,
      timestamp: Date.now(),
      loading: false,
      error: newError,
    })
  }

  // Fetch data
  const fetch = async (force = false) => {
    // Use cache if valid and not forcing
    if (!force && loadFromCache()) {
      return
    }

    loading.value = true
    error.value = null

    try {
      const result = await fetchFn()
      data.value = result
      lastFetchTime.value = Date.now()
      saveToCache(result)
    } catch (err) {
      error.value = err as Error
      saveToCache(null, err as Error)
      console.error(`[Cache] Failed to fetch ${key}:`, err)
    } finally {
      loading.value = false
    }
  }

  // Refresh data
  const refresh = () => fetch(true)

  // Invalidate cache
  const invalidate = () => {
    globalCache.delete(key)
  }

  // Optimistic update
  const updateOptimistic = (updater: (current: T) => T) => {
    if (!optimisticUpdate || !data.value) return

    const previousValue = data.value
    try {
      data.value = updater(data.value)
      saveToCache(data.value)
    } catch (err) {
      // Rollback on error
      data.value = previousValue
      console.error('[Cache] Optimistic update failed:', err)
    }
  }

  // Subscribe to WebSocket updates if enabled
  if (subscribeToUpdates && scope) {
    try {
      const { subscribe } = useWebSocket()
      unsubscribe = subscribe({
        scope,
        resourceId,
        events,
        onMessage: (message) => {
          handleRealtimeUpdate(message)
        },
        throttle: 500, // Throttle updates to 500ms
      })
    } catch (err) {
      console.warn('[Cache] WebSocket not available, updates will be manual only')
    }
  }

  // Handle real-time updates
  const handleRealtimeUpdate = (message: any) => {
    try {
      const { type, data: updateData } = message

      // Handle different update types
      switch (type) {
        case 'created':
          if (Array.isArray(data.value)) {
            data.value = [...data.value, updateData] as any
            saveToCache(data.value)
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
              saveToCache(data.value)
            }
          } else if ((data.value as any)?.id === updateData?.id) {
            data.value = updateData
            saveToCache(data.value)
          }
          break

        case 'deleted':
          if (Array.isArray(data.value)) {
            data.value = (data.value as any[]).filter((item) => item?.id !== updateData?.id) as any
            saveToCache(data.value)
          } else if ((data.value as any)?.id === updateData?.id) {
            data.value = null
            saveToCache(null)
          }
          break

        case 'bulk_update':
          // For bulk updates, just refresh to avoid overwhelming the UI
          refresh()
          break

        default:
          // Unknown update type, refresh to be safe
          refresh()
      }
    } catch (err) {
      console.error('[Cache] Failed to handle realtime update:', err)
      // On error, refresh to ensure consistency
      refresh()
    }
  }

  // Initial fetch
  if (refetchOnMount) {
    fetch()
  } else {
    loadFromCache()
  }

  // Cleanup
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
    updateOptimistic,
    isCacheValid: computed(() => isCacheValid()),
    lastFetchTime: computed(() => lastFetchTime.value),
  }
}

// Clear all cache
export function clearAllCache() {
  globalCache.clear()
}

// Clear specific cache entries
export function clearCache(key: string) {
  globalCache.delete(key)
}

// Get cache stats
export function getCacheStats() {
  return {
    size: globalCache.size,
    keys: Array.from(globalCache.keys()),
  }
}
