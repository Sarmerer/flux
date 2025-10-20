import { ref } from 'vue'

export interface Toast {
  id?: string
  title?: string
  description?: string
  variant?: 'default' | 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

const toasts = ref<Toast[]>([])
let toastId = 0

export function useToast() {
  const show = (options: Toast | string) => {
    const id = `toast-${++toastId}`

    let toast: Toast
    if (typeof options === 'string') {
      toast = { id, description: options, variant: 'default', duration: 3000 }
    } else {
      toast = {
        id,
        variant: 'default',
        duration: 3000,
        ...options,
      }
    }

    toasts.value.push(toast)

    if (toast.duration && toast.duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, toast.duration)
    }

    return id
  }

  const success = (title: string, description?: string, duration?: number) =>
    show({ title, description, variant: 'success', duration })

  const error = (title: string, description?: string, duration?: number) =>
    show({ title, description, variant: 'error', duration })

  const warning = (title: string, description?: string, duration?: number) =>
    show({ title, description, variant: 'warning', duration })

  const info = (title: string, description?: string, duration?: number) =>
    show({ title, description, variant: 'info', duration })

  const removeToast = (id: string) => {
    const index = toasts.value.findIndex((t) => t.id === id)
    if (index !== -1) {
      toasts.value.splice(index, 1)
    }
  }

  const clear = () => {
    toasts.value = []
  }

  return {
    toasts,
    show,
    success,
    error,
    warning,
    info,
    removeToast,
    clear,
  }
}
