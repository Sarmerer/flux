import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import type { Ref } from 'vue'
import type { Updater } from '@tanstack/vue-table'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export async function safeJson<T = unknown>(res: Response): Promise<T | null> {
  try {
    return (await res.json()) as T
  } catch {
    return null
  }
}

export function valueUpdater<T>(updaterOrValue: Updater<T>, ref: Ref<T>) {
  ref.value = typeof updaterOrValue === 'function'
    ? (updaterOrValue as (old: T) => T)(ref.value)
    : updaterOrValue
}

