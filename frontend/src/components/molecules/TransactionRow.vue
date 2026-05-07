<script setup lang="ts">
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseBadge from '@/components/atoms/BaseBadge.vue';

interface Props {
  type: 'incoming' | 'outgoing' | 'topup' | 'withdraw';
  amount: number;
  currency: string;
  recipient: string;
  date: string;
  status: 'completed' | 'pending' | 'failed';
}

defineProps<Props>();

const typeConfig = {
  incoming: { icon: 'ArrowDownLeft', color: 'text-green-500', label: 'Received' },
  outgoing: { icon: 'ArrowUpRight', color: 'text-red-500', label: 'Sent' },
  topup: { icon: 'PlusCircle', color: 'text-blue-500', label: 'Top Up' },
  withdraw: { icon: 'MinusCircle', color: 'text-orange-500', label: 'Withdraw' },
};

const statusConfig = {
  completed: { variant: 'default' as const, label: 'Completed' },
  pending: { variant: 'secondary' as const, label: 'Pending' },
  failed: { variant: 'destructive' as const, label: 'Failed' },
};
</script>

<template>
  <div class="flex items-center justify-between p-4 hover:bg-muted/50 rounded-lg transition-colors">
    <div class="flex items-center space-x-4">
      <div :class="['p-2 rounded-full bg-muted', typeConfig[type].color]">
        <BaseIcon :name="typeConfig[type].icon" class="h-5 w-5" />
      </div>
      <div>
        <p class="text-sm font-medium leading-none">{{ typeConfig[type].label }} - {{ recipient }}</p>
        <p class="text-xs text-muted-foreground">{{ date }}</p>
      </div>
    </div>
    <div class="text-right space-y-1">
      <p :class="['text-sm font-bold', typeConfig[type].color]">
        {{ type === 'incoming' || type === 'topup' ? '+' : '-' }}
        {{ currency }} {{ amount.toLocaleString() }}
      </p>
      <BaseBadge :variant="statusConfig[status].variant" class="text-[10px] h-4 px-1">
        {{ statusConfig[status].label }}
      </BaseBadge>
    </div>
  </div>
</template>
