<script setup lang="ts">
import WalletBalanceCard from '@/components/molecules/WalletBalanceCard.vue';
import StatsCard from '@/components/molecules/StatsCard.vue';

interface Props {
  balance: number;
  currency: string;
  stats: {
    totalIncome: number;
    totalExpense: number;
    activeTransactions: number;
  };
  loading?: boolean;
}

defineProps<Props>();
</script>

<template>
  <div class="space-y-6">
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <WalletBalanceCard
        :balance="balance"
        :currency="currency"
        :loading="loading"
      />
      <StatsCard
        title="Income"
        :value="`+${currency} ${stats.totalIncome.toLocaleString()}`"
        description="This month"
        icon="ArrowDownLeft"
        :trend="{ value: 12, isPositive: true }"
      />
      <StatsCard
        title="Expense"
        :value="`-${currency} ${stats.totalExpense.toLocaleString()}`"
        description="This month"
        icon="ArrowUpRight"
        :trend="{ value: 4, isPositive: false }"
      />
      <StatsCard
        title="Active Escrows"
        :value="stats.activeTransactions"
        description="Pending completion"
        icon="Clock"
      />
    </div>
  </div>
</template>
