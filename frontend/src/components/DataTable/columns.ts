import DropdownAction from '@/components/DataTable/DropDown.vue'
import type { ColumnDef } from '@tanstack/vue-table'
import { h } from 'vue'

export interface Payment {
  id: string
  amount: number
  status: 'pending' | 'processing' | 'success' | 'failed'
  email: string
}

export const columns: ColumnDef<Payment>[] = [
  {
    id: 'actions',
    enableHiding: false,
    cell: ({ row }) => {
      const payment = row.original

      return h(
        'div',
        { class: 'relative' },
        h(DropdownAction, {
          payment,
        })
      )
    },
  },
]
