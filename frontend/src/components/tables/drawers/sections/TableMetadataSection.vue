<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface TableMetadata {
  name: string
  description: string
}

interface Props {
  modelValue: TableMetadata
  validationErrors?: Record<string, string>
  disabled?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: TableMetadata): void
}

const props = withDefaults(defineProps<Props>(), {
  validationErrors: () => ({}),
  disabled: false,
})

const emit = defineEmits<Emits>()

const updateField = <K extends keyof TableMetadata>(field: K, value: TableMetadata[K]) => {
  emit('update:modelValue', { ...props.modelValue, [field]: value })
}
</script>

<template>
  <div class="space-y-3">
    <div class="space-y-1.5">
      <Label for="table-name" class="text-xs font-medium">
        Name <span class="text-red-500">*</span>
      </Label>
      <Input
        id="table-name"
        :model-value="modelValue.name"
        @update:model-value="(value) => updateField('name', String(value ?? ''))"
        placeholder="e.g., users, products, orders"
        :class="validationErrors.name && 'border-red-500'"
        :disabled="disabled"
        class="h-8 text-sm"
      />
      <p v-if="validationErrors.name" class="text-xs text-red-500">
        {{ validationErrors.name }}
      </p>
      <p v-else class="text-xs text-muted-foreground">
        Use lowercase letters, numbers, and underscores only
      </p>
    </div>

    <div class="space-y-1.5">
      <Label for="table-description" class="text-xs font-medium">Description</Label>
      <Input
        id="table-description"
        :model-value="modelValue.description"
        @update:model-value="(value) => updateField('description', String(value ?? ''))"
        placeholder="Optional"
        :disabled="disabled"
        class="h-8 text-sm"
      />
    </div>
  </div>
</template>
