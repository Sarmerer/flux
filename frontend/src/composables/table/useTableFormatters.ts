export function useTableFormatters() {
  const formatCellValue = (value: any, columnType: string): string => {
    if (value === null || value === undefined) return '-'

    const upperType = columnType.toUpperCase()

    if (upperType === 'JSON' || upperType === 'JSONB') {
      if (typeof value === 'object') {
        const str = JSON.stringify(value)
        return str.length > 50 ? str.substring(0, 50) + '...' : str
      }
    }

    if (upperType === 'BOOLEAN' || upperType === 'BOOL') {
      return value ? 'true' : 'false'
    }

    if (upperType.includes('TIMESTAMP') || upperType === 'DATE') {
      try {
        const date = new Date(value)
        return date.toLocaleString()
      } catch {
        return String(value)
      }
    }

    const str = String(value)
    return str.length > 100 ? str.substring(0, 100) + '...' : str
  }

  const truncateText = (text: string, maxLength: number): string => {
    if (text.length <= maxLength) return text
    return text.substring(0, maxLength) + '...'
  }

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'

    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  const formatNumber = (num: number): string => {
    return new Intl.NumberFormat().format(num)
  }

  return {
    formatCellValue,
    truncateText,
    formatBytes,
    formatNumber,
  }
}
