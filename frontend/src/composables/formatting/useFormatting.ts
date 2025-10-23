import { useI18n } from 'vue-i18n'

import { type Locale, formatDistanceToNowStrict, parseISO } from 'date-fns'
import { de, enUS, es, fr } from 'date-fns/locale'

type SupportedLocale = 'en-US' | 'fr-FR' | 'de-DE' | 'es-ES'

const dateFnsLocales: Record<SupportedLocale, Locale> = {
  'en-US': enUS,
  'fr-FR': fr,
  'de-DE': de,
  'es-ES': es,
}

export function useFormatting() {
  const { locale } = useI18n()

  const formatDate = (
    date: string | number | Date,
    options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' }
  ): string => {
    const d = new Date(date)
    if (isNaN(d.getTime())) return ''
    return d.toLocaleDateString(locale.value, options)
  }

  const formatDateTime = (
    date: string | number | Date,
    options: Intl.DateTimeFormatOptions = {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }
  ): string => {
    const d = new Date(date)
    if (isNaN(d.getTime())) return ''
    return d.toLocaleString(locale.value, options)
  }

  const formatTime = (
    date: string | number | Date,
    options: Intl.DateTimeFormatOptions = { hour: '2-digit', minute: '2-digit' }
  ): string => {
    const d = new Date(date)
    if (isNaN(d.getTime())) return ''
    return d.toLocaleTimeString(locale.value, options)
  }

  const formatRelativeTime = (date: string | number | Date): string => {
    try {
      const d = typeof date === 'string' ? parseISO(date) : new Date(date)
      const localeKey = Object.keys(dateFnsLocales).includes(locale.value)
        ? (locale.value as SupportedLocale)
        : 'en-US'
      return formatDistanceToNowStrict(d, {
        addSuffix: true,
        locale: dateFnsLocales[localeKey],
      })
    } catch {
      return ''
    }
  }

  const formatNumber = (value: number, options?: Intl.NumberFormatOptions): string => {
    return new Intl.NumberFormat(locale.value, options).format(value)
  }

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    const value = parseFloat((bytes / Math.pow(k, i)).toFixed(2))
    return `${value} ${sizes[i]}`
  }

  const formatDuration = (ms: number): string => {
    const sec = Math.floor(ms / 1000)
    const min = Math.floor(sec / 60)
    const hr = Math.floor(min / 60)
    if (hr > 0) return `${hr}h ${min % 60}m`
    if (min > 0) return `${min}m ${sec % 60}s`
    return `${sec}s`
  }

  return {
    formatDate,
    formatDateTime,
    formatTime,
    formatRelativeTime,
    formatNumber,
    formatBytes,
    formatDuration,
  }
}
