export interface TypeConversionOptions {
  allowNull?: boolean
  strict?: boolean
}

export interface ColumnValidation {
  type: string
  isNullable: boolean
  name: string
}

export function useColumnTypeUtils() {
  const convertValueByType = (
    value: any,
    type: string,
    options: TypeConversionOptions = {}
  ): any => {
    if (value === null || value === '') {
      return options.allowNull !== false ? null : value
    }

    const upperType = type.toUpperCase()

    if (upperType.includes('INT') || upperType.includes('SERIAL')) {
      const parsed = parseInt(value, 10)
      if (options.strict && isNaN(parsed)) {
        throw new Error(`Cannot convert "${value}" to integer`)
      }
      return parsed
    }

    if (
      upperType.includes('DECIMAL') ||
      upperType.includes('NUMERIC') ||
      upperType.includes('REAL') ||
      upperType.includes('DOUBLE')
    ) {
      const parsed = parseFloat(value)
      if (options.strict && isNaN(parsed)) {
        throw new Error(`Cannot convert "${value}" to number`)
      }
      return parsed
    }

    if (upperType === 'BOOLEAN' || upperType === 'BOOL') {
      if (typeof value === 'boolean') return value
      if (typeof value === 'string') {
        const lower = value.toLowerCase()
        return lower === 'true' || lower === '1' || lower === 'yes'
      }
      return Boolean(value)
    }

    if (upperType === 'JSON' || upperType === 'JSONB') {
      if (typeof value === 'string') {
        try {
          return JSON.parse(value)
        } catch (error) {
          if (options.strict) {
            throw new Error(`Invalid JSON: ${error}`)
          }
          return value
        }
      }
      return value
    }

    return value
  }

  const getInputType = (columnType: string): string => {
    const upperType = columnType.toUpperCase()

    if (upperType.includes('INT') || upperType.includes('SERIAL')) {
      return 'number'
    }

    if (
      upperType.includes('DECIMAL') ||
      upperType.includes('NUMERIC') ||
      upperType.includes('REAL') ||
      upperType.includes('DOUBLE')
    ) {
      return 'number'
    }

    if (upperType === 'BOOLEAN' || upperType === 'BOOL') {
      return 'checkbox'
    }

    if (upperType === 'DATE') {
      return 'date'
    }

    if (upperType.includes('TIMESTAMP') || upperType.includes('TIME')) {
      return 'datetime-local'
    }

    return 'text'
  }

  const isTextArea = (columnType: string): boolean => {
    const upperType = columnType.toUpperCase()
    return upperType === 'TEXT' || upperType === 'JSON' || upperType === 'JSONB'
  }

  const validateColumnValue = (value: any, column: ColumnValidation): string | null => {
    if (!column.isNullable && (value === '' || value === null || value === undefined)) {
      return `${column.name} is required`
    }

    if (value === '' || value === null) {
      return null
    }

    const upperType = column.type.toUpperCase()

    if (upperType.includes('INT') && isNaN(parseInt(value, 10))) {
      return `${column.name} must be a valid integer`
    }

    if (
      (upperType.includes('DECIMAL') ||
        upperType.includes('NUMERIC') ||
        upperType.includes('REAL') ||
        upperType.includes('DOUBLE')) &&
      isNaN(parseFloat(value))
    ) {
      return `${column.name} must be a valid number`
    }

    if ((upperType === 'JSON' || upperType === 'JSONB') && typeof value === 'string') {
      try {
        JSON.parse(value)
      } catch {
        return `${column.name} must be valid JSON`
      }
    }

    return null
  }

  return {
    convertValueByType,
    getInputType,
    isTextArea,
    validateColumnValue,
  }
}
