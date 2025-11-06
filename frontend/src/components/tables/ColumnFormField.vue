<script setup lang="ts">
import { computed, watch } from 'vue'

import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

export interface ColumnFormData {
  name: string
  type: string
  nullable: boolean
  primary_key: boolean
  default_value: string
  is_identity: boolean
  unique: boolean
}

interface Props {
  modelValue: ColumnFormData
  mode?: 'inline' | 'full'
  disabled?: boolean
  existingColumns?: string[]
  showLabels?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: ColumnFormData): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'inline',
  disabled: false,
  existingColumns: () => [],
  showLabels: true,
})

const emit = defineEmits<Emits>()

const dataTypes = {
  numeric: [
    { value: 'SMALLINT', label: 'Small Integer', description: '2 bytes, -32K to 32K' },
    { value: 'INTEGER', label: 'Integer', description: '4 bytes, -2B to 2B' },
    { value: 'BIGINT', label: 'Big Integer', description: '8 bytes, very large numbers' },
    { value: 'DECIMAL', label: 'Decimal', description: 'Exact numeric with precision' },
    { value: 'NUMERIC', label: 'Numeric', description: 'Same as DECIMAL' },
    { value: 'REAL', label: 'Real', description: '4 bytes, 6 decimal precision' },
    { value: 'DOUBLE PRECISION', label: 'Double', description: '8 bytes, 15 decimal precision' },
  ],
  text: [
    { value: 'TEXT', label: 'Text', description: 'Unlimited length' },
    { value: 'VARCHAR', label: 'Varchar', description: 'Variable length with limit' },
    { value: 'CHAR', label: 'Char', description: 'Fixed length' },
  ],
  datetime: [
    { value: 'TIMESTAMP', label: 'Timestamp', description: 'Date and time' },
    { value: 'TIMESTAMPTZ', label: 'Timestamp TZ', description: 'With timezone' },
    { value: 'DATE', label: 'Date', description: 'Date only' },
    { value: 'TIME', label: 'Time', description: 'Time only' },
  ],
  special: [
    { value: 'UUID', label: 'UUID', description: 'Unique identifier' },
    { value: 'BOOLEAN', label: 'Boolean', description: 'True or false' },
    { value: 'JSON', label: 'JSON', description: 'JSON data' },
    { value: 'JSONB', label: 'JSONB', description: 'Binary JSON (faster)' },
  ],
}

const localValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const updateField = <K extends keyof ColumnFormData>(field: K, value: ColumnFormData[K]) => {
  emit('update:modelValue', { ...props.modelValue, [field]: value })
}

const defaultValueSuggestions = computed(() => {
  const type = localValue.value.type.toUpperCase()

  if (type.includes('TIMESTAMP')) {
    return [
      { value: 'NOW()', label: 'NOW()' },
      { value: 'CURRENT_TIMESTAMP', label: 'CURRENT_TIMESTAMP' },
    ]
  }

  if (type === 'DATE') {
    return [{ value: 'CURRENT_DATE', label: 'CURRENT_DATE' }]
  }

  if (type === 'UUID') {
    return [{ value: 'gen_random_uuid()', label: 'gen_random_uuid()' }]
  }

  if (type === 'BOOLEAN') {
    return [
      { value: 'TRUE', label: 'TRUE' },
      { value: 'FALSE', label: 'FALSE' },
    ]
  }

  if (['SMALLINT', 'INTEGER', 'BIGINT', 'DECIMAL', 'NUMERIC'].includes(type)) {
    return [
      { value: '0', label: '0' },
      { value: '1', label: '1' },
    ]
  }

  return []
})

const isNumericType = computed(() => {
  const type = localValue.value.type.toUpperCase()
  return ['SMALLINT', 'INTEGER', 'BIGINT'].includes(type)
})

const canHaveDefault = computed(() => {
  return !localValue.value.is_identity
})

const primaryKey = computed({
  get: () => localValue.value.primary_key,
  set: (value) => {
    const updates: Partial<ColumnFormData> = { primary_key: value }
    if (value) {
      updates.nullable = false
    }
    emit('update:modelValue', { ...props.modelValue, ...updates })
  },
})

const isNullableDisabled = computed(() => {
  return localValue.value.is_identity || localValue.value.primary_key
})

watch(
  () => localValue.value.is_identity,
  (isIdentity) => {
    if (isIdentity) {
      updateField('default_value', '')
      updateField('nullable', false)
    }
  }
)

const setDefaultValue = (value: string) => {
  updateField('default_value', value)
}
</script>

<template>
  <div :class="mode === 'inline' ? 'space-y-2.5' : 'space-y-4'">
    <div :class="mode === 'inline' ? 'grid grid-cols-2 gap-2' : 'space-y-4'">
      <div class="space-y-1">
        <Label v-if="showLabels" class="text-[11px] font-medium text-muted-foreground">
          Column Name <span class="text-red-500">*</span>
        </Label>
        <Input
          :model-value="localValue.name"
          @update:model-value="(value) => updateField('name', String(value ?? ''))"
          placeholder="e.g., email, created_at"
          :disabled="disabled"
          class="h-7 text-xs"
        />
      </div>

      <div class="space-y-1">
        <Label v-if="showLabels" class="text-[11px] font-medium text-muted-foreground">
          Data Type <span class="text-red-500">*</span>
        </Label>
        <Select
          :model-value="localValue.type"
          @update:model-value="(value) => updateField('type', String(value ?? 'TEXT'))"
        >
          <SelectTrigger class="h-7 text-xs">
            <SelectValue placeholder="Select type">
              <template v-if="localValue.type" #default>
                {{ dataTypes.numeric.concat(dataTypes.text, dataTypes.datetime, dataTypes.special).find(t => t.value === localValue.type)?.label || localValue.type }}
              </template>
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            <SelectGroup v-for="(types, category) in dataTypes" :key="category">
              <SelectLabel class="text-[11px] font-semibold uppercase text-muted-foreground">
                {{ category }}
              </SelectLabel>
              <SelectItem v-for="type in types" :key="type.value" :value="type.value" class="text-xs">
                <div class="flex items-center gap-2">
                  <span>{{ type.label }}</span>
                  <span class="text-[11px] text-muted-foreground">{{ type.description }}</span>
                </div>
              </SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>
    </div>

    <div v-if="canHaveDefault" class="space-y-1">
      <Label class="text-[11px] font-medium text-muted-foreground">Default Value</Label>
      <Input
        :model-value="localValue.default_value"
        @update:model-value="(value) => updateField('default_value', String(value ?? ''))"
        placeholder="e.g., '', 0, gen_random_uuid()"
        class="h-7 text-xs"
      />

      <div v-if="defaultValueSuggestions.length > 0" class="flex flex-wrap gap-1">
        <Button
          v-for="suggestion in defaultValueSuggestions"
          :key="suggestion.value"
          type="button"
          variant="outline"
          size="sm"
          @click="setDefaultValue(suggestion.value)"
          class="h-5 text-[11px] px-1.5"
        >
          {{ suggestion.label }}
        </Button>
      </div>
    </div>

    <div class="flex flex-wrap gap-3">
      <div class="flex items-center space-x-1.5">
        <Checkbox
          :checked="!!localValue.nullable"
          @update:checked="(value: boolean) => updateField('nullable', !!value)"
          id="nullable"
          :disabled="isNullableDisabled"
          class="h-3.5 w-3.5"
        />
        <Label
          for="nullable"
          :class="['text-[11px] font-normal', isNullableDisabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer']"
        >
          Nullable
        </Label>
      </div>

      <div class="flex items-center space-x-1.5">
        <Checkbox
          :checked="!!primaryKey"
          @update:checked="(value: boolean) => primaryKey = !!value"
          id="primary-key"
          class="h-3.5 w-3.5"
        />
        <Label
          for="primary-key"
          class="text-[11px] font-normal cursor-pointer"
        >
          Primary Key
        </Label>
      </div>

      <div class="flex items-center space-x-1.5">
        <Checkbox
          :checked="!!localValue.unique"
          @update:checked="(value: boolean) => updateField('unique', !!value)"
          id="unique"
          class="h-3.5 w-3.5"
        />
        <Label
          for="unique"
          class="text-[11px] font-normal cursor-pointer"
        >
          Unique
        </Label>
      </div>

      <div v-if="isNumericType" class="flex items-center space-x-1.5">
        <Checkbox
          :checked="!!localValue.is_identity"
          @update:checked="(value: boolean) => updateField('is_identity', !!value)"
          id="identity"
          class="h-3.5 w-3.5"
        />
        <Label
          for="identity"
          class="text-[11px] font-normal cursor-pointer"
        >
          Auto-increment
        </Label>
      </div>
    </div>
  </div>
</template>
