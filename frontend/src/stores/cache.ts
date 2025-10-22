import { computed, ref } from 'vue'

import { defineStore } from 'pinia'

interface CacheEntry<T> {
  data: T | null
  timestamp: number
  loading: boolean
  error: Error | null
  tags?: string[]
}

export const useCacheStore = defineStore('cache', () => {
  const cache = ref<Map<string, CacheEntry<any>>>(new Map())

  const get = <T>(key: string): CacheEntry<T> | undefined => {
    return cache.value.get(key)
  }

  const set = <T>(
    key: string,
    data: T | null,
    error: Error | null = null,
    options?: { tags?: string[] }
  ) => {
    cache.value.set(key, {
      data,
      timestamp: Date.now(),
      loading: false,
      error,
      tags: options?.tags,
    })
  }

  const isValid = (key: string, ttlMs: number): boolean => {
    const entry = cache.value.get(key)
    if (!entry) return false

    const age = Date.now() - entry.timestamp
    return age < ttlMs
  }

  const invalidate = (key: string) => {
    cache.value.delete(key)
  }

  const invalidatePattern = (pattern: string | RegExp) => {
    const regex = typeof pattern === 'string' ? new RegExp(pattern.replace('*', '.*')) : pattern

    const keysToInvalidate = Array.from(cache.value.keys()).filter((key) => regex.test(key))

    keysToInvalidate.forEach((key) => invalidate(key))
  }

  const invalidateByTag = (tag: string) => {
    const keysToInvalidate = Array.from(cache.value.entries())
      .filter(([_, entry]) => entry.tags?.includes(tag))
      .map(([key]) => key)

    keysToInvalidate.forEach((key) => invalidate(key))
  }

  const clear = () => {
    cache.value.clear()
  }

  const clearStale = (maxAgeMs: number = 5 * 60 * 1000) => {
    const now = Date.now()
    const keysToDelete: string[] = []

    cache.value.forEach((entry, key) => {
      if (now - entry.timestamp > maxAgeMs) {
        keysToDelete.push(key)
      }
    })

    keysToDelete.forEach((key) => cache.value.delete(key))
  }

  const stats = computed(() => ({
    size: cache.value.size,
    keys: Array.from(cache.value.keys()),
    entries: Array.from(cache.value.entries()).map(([key, entry]) => ({
      key,
      timestamp: entry.timestamp,
      age: Date.now() - entry.timestamp,
      hasError: !!entry.error,
      hasData: entry.data !== null,
      tags: entry.tags,
    })),
  }))

  const hydrate = (storageKey: string = 'flow-cache') => {
    try {
      const stored = localStorage.getItem(storageKey)
      if (stored) {
        const parsed = JSON.parse(stored)
        cache.value = new Map(parsed)
      }
    } catch (err) {
      console.error('[Cache] Failed to hydrate from localStorage:', err)
    }
  }

  const persist = (storageKey: string = 'flow-cache') => {
    try {
      const serialized = JSON.stringify(Array.from(cache.value.entries()))
      localStorage.setItem(storageKey, serialized)
    } catch (err) {
      console.error('[Cache] Failed to persist to localStorage:', err)
    }
  }

  return {
    cache: computed(() => cache.value),
    get,
    set,
    isValid,
    invalidate,
    invalidatePattern,
    invalidateByTag,
    clear,
    clearStale,
    stats,
    hydrate,
    persist,
  }
})
