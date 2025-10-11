<script setup lang="ts">
import DataTable from '@/components/DataTable/index.vue'
import { onMounted, ref } from 'vue'

import { type Payment, columns } from '@/components/DataTable/columns'

const data = ref<Payment[]>([])

async function getData(): Promise<Payment[]> {
  // Fetch data from your API here.
  return [
    {
      id: '728ed52f',
      amount: 100,
      status: 'pending',
      email: 'm@example.com',
    },
  ]
}

onMounted(async () => {
  data.value = await getData()
})
</script>

<template>
  <div class="container py-10 px-6 mx-auto">
    <DataTable :columns="columns" :data="data" />
  </div>
</template>
