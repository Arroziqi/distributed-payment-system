<script setup lang="ts">
import BaseTable from '@/components/atoms/BaseTable.vue';
import BaseBadge from '@/components/atoms/BaseBadge.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';

interface Transaction {
  id: string;
  type: string;
  amount: number;
  currency: string;
  recipient: string;
  createdAt: string;
  status: string;
}

interface Props {
  transactions: Transaction[];
  loading?: boolean;
}

defineProps<Props>();

const columns = [
  { key: 'type', label: 'Type' },
  { key: 'recipient', label: 'Recipient/Source' },
  { key: 'amount', label: 'Amount' },
  { key: 'createdAt', label: 'Date' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: '', class: 'text-right' },
];

const getStatusVariant = (status: string) => {
  switch (status.toLowerCase()) {
    case 'completed': return 'default';
    case 'pending': return 'secondary';
    case 'failed': return 'destructive';
    default: return 'outline';
  }
};
</script>

<template>
  <BaseTable :columns="columns" :data="transactions" :loading="loading">
    <template #cell(type)="{ row }">
      <div class="flex items-center">
        <div :class="['mr-2 p-1 rounded-full bg-muted', row.amount > 0 ? 'text-green-500' : 'text-red-500']">
          <BaseIcon :name="row.amount > 0 ? 'ArrowDownLeft' : 'ArrowUpRight'" size="14" />
        </div>
        <span class="capitalize">{{ row.type }}</span>
      </div>
    </template>
    
    <template #cell(amount)="{ row }">
      <span :class="['font-medium', row.amount > 0 ? 'text-green-600' : 'text-red-600']">
        {{ row.amount > 0 ? '+' : '' }}{{ row.currency }} {{ row.amount.toLocaleString() }}
      </span>
    </template>

    <template #cell(status)="{ row }">
      <BaseBadge :variant="getStatusVariant(row.status)" class="capitalize">
        {{ row.status }}
      </BaseBadge>
    </template>

    <template #cell(actions)>
      <BaseButton variant="ghost" size="icon">
        <BaseIcon name="MoreHorizontal" class="h-4 w-4" />
      </BaseButton>
    </template>
  </BaseTable>
</template>
