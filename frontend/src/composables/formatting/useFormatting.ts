/**
 * Shared formatting utilities for dates, times, and numbers
 */
export function useFormatting() {
  /**
   * Format a date string to a human-readable relative time
   * @param dateString ISO date string
   * @returns Formatted relative time string (e.g., "2 hours ago", "3 days ago")
   */
  const formatRelativeTime = (dateString: string): string => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffSeconds = Math.floor(diffMs / 1000)
    const diffMinutes = Math.floor(diffSeconds / 60)
    const diffHours = Math.floor(diffMinutes / 60)
    const diffDays = Math.floor(diffHours / 24)

    if (diffSeconds < 60) {
      return 'Just now'
    } else if (diffMinutes < 60) {
      return `${diffMinutes} minute${diffMinutes > 1 ? 's' : ''} ago`
    } else if (diffHours < 24) {
      return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`
    } else if (diffDays < 7) {
      return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`
    } else if (diffDays < 30) {
      const weeks = Math.floor(diffDays / 7)
      return `${weeks} week${weeks > 1 ? 's' : ''} ago`
    } else if (diffDays < 365) {
      const months = Math.floor(diffDays / 30)
      return `${months} month${months > 1 ? 's' : ''} ago`
    } else {
      const years = Math.floor(diffDays / 365)
      return `${years} year${years > 1 ? 's' : ''} ago`
    }
  }

  /**
   * Format a date string to a localized date format
   * @param dateString ISO date string
   * @param options Intl.DateTimeFormatOptions
   * @returns Formatted date string
   */
  const formatDate = (
    dateString: string,
    options: Intl.DateTimeFormatOptions = {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }
  ): string => {
    return new Date(dateString).toLocaleDateString(undefined, options)
  }

  /**
   * Format a date string to include time
   * @param dateString ISO date string
   * @returns Formatted date and time string
   */
  const formatDateTime = (dateString: string): string => {
    return new Date(dateString).toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  /**
   * Format a number with thousand separators
   * @param value Number to format
   * @returns Formatted number string
   */
  const formatNumber = (value: number): string => {
    return value.toLocaleString()
  }

  /**
   * Format bytes to human-readable size
   * @param bytes Number of bytes
   * @returns Formatted size string (e.g., "1.5 MB")
   */
  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'

    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  /**
   * Format duration in milliseconds to human-readable string
   * @param ms Duration in milliseconds
   * @returns Formatted duration string (e.g., "2h 30m", "45s")
   */
  const formatDuration = (ms: number): string => {
    const seconds = Math.floor(ms / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)

    if (hours > 0) {
      const remainingMinutes = minutes % 60
      return `${hours}h${remainingMinutes > 0 ? ` ${remainingMinutes}m` : ''}`
    } else if (minutes > 0) {
      const remainingSeconds = seconds % 60
      return `${minutes}m${remainingSeconds > 0 ? ` ${remainingSeconds}s` : ''}`
    } else {
      return `${seconds}s`
    }
  }

  return {
    formatRelativeTime,
    formatDate,
    formatDateTime,
    formatNumber,
    formatBytes,
    formatDuration
  }
}
