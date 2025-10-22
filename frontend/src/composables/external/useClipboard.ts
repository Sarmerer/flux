import { ref } from 'vue'

import { useToast } from './useToast'

export function useClipboard() {
  const { success, error } = useToast()
  const copied = ref(false)

  const copy = async (text: string, successMessage = 'Copied to clipboard') => {
    try {
      await navigator.clipboard.writeText(text)
      copied.value = true
      success(successMessage)

      setTimeout(() => {
        copied.value = false
      }, 2000)
    } catch (err) {
      error('Failed to copy to clipboard')
      console.error('Copy failed:', err)
    }
  }

  return {
    copy,
    copied,
  }
}
