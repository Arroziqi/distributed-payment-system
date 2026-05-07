<script setup lang="ts">
import { ref } from 'vue';
import DashboardOverview from '@/components/organisms/DashboardOverview.vue';
import TransactionTable from '@/components/organisms/TransactionTable.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';

const balance = ref(12500);
const currency = ref('USD');
const stats = ref({
  totalIncome: 15000,
  totalExpense: 2500,
  activeTransactions: 3
});

const transactions = ref([
  { id: '1', type: 'incoming', amount: 500, currency: 'USD', recipient: 'John Doe', createdAt: '2026-05-07 10:00', status: 'completed' },
  { id: '2', type: 'outgoing', amount: -150, currency: 'USD', recipient: 'Amazon', createdAt: '2026-05-07 09:30', status: 'completed' },
  { id: '3', type: 'outgoing', amount: -200, currency: 'USD', recipient: 'Starbucks', createdAt: '2026-05-07 08:45', status: 'pending' },
]);

const loading = ref(false);
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Dashboard</h2>
      <div class="flex items-center space-x-2">
        <BaseButton>
          <BaseIcon name="Download" class="mr-2 h-4 w-4" />
          Download Report
        </BaseButton>
      </div>
    </div>
    
    <DashboardOverview
      :balance="balance"
      :currency="currency"
      :stats="stats"
      :loading="loading"
    />

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <div class="col-span-4 rounded-xl border bg-card text-card-foreground shadow">
        <div class="p-6">
          <h3 class="text-lg font-semibold mb-4">Recent Transactions</h3>
          <TransactionTable :transactions="transactions" :loading="loading" />
        </div>
      </div>
      
      <div class="col-span-3 rounded-xl border bg-card text-card-foreground shadow">
        <div class="p-6">
          <h3 class="text-lg font-semibold mb-4">Quick Actions</h3>
          <div class="grid grid-cols-2 gap-4">
            <BaseButton variant="outline" class="h-24 flex-col gap-2">
              <BaseIcon name="Send" class="h-6 w-6" />
              Transfer
            </BaseButton>
            <BaseButton variant="outline" class="h-24 flex-col gap-2">
              <BaseIcon name="PlusCircle" class="h-6 w-6" />
              Top Up
            </BaseButton>
            <BaseButton variant="outline" class="h-24 flex-col gap-2">
              <BaseIcon name="ArrowDownCircle" class="h-6 w-6" />
              Withdraw
            </BaseButton>
            <BaseButton variant="outline" class="h-24 flex-col gap-2">
              <BaseIcon name="CreditCard" class="h-6 w-6" />
              Cards
            </BaseButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
