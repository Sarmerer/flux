import { ref } from 'vue'

interface ConfirmOptions {
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  type?: 'danger' | 'warning' | 'info'
}

const isOpen = ref(false)
const options = ref<ConfirmOptions | null>(null)
let resolveCallback: ((value: boolean) => void) | null = null

export function useConfirm() {
  const confirm = (opts: ConfirmOptions): Promise<boolean> => {
    return new Promise((resolve) => {
      options.value = {
        confirmText: 'Confirm',
        cancelText: 'Cancel',
        type: 'info',
        ...opts,
      }
      isOpen.value = true
      resolveCallback = resolve
    })
  }

  const handleConfirm = () => {
    isOpen.value = false
    if (resolveCallback) {
      resolveCallback(true)
      resolveCallback = null
    }
  }

  const handleCancel = () => {
    isOpen.value = false
    if (resolveCallback) {
      resolveCallback(false)
      resolveCallback = null
    }
  }

  return {
    confirm,
    isOpen,
    options,
    handleConfirm,
    handleCancel,
  }
}
