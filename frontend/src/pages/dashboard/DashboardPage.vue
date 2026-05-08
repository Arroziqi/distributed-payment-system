<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import DashboardOverview from '@/components/organisms/DashboardOverview.vue';
import TransactionTable from '@/components/organisms/TransactionTable.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import { useWallet } from '@/features/wallet/useWallet';
import { useTransactions } from '@/features/transactions/useTransactions';

const router = useRouter();
const { balance, lockedBalance, wallet, loading: walletLoading } = useWallet();
const { transactions: allTransactions, loading: transactionsLoading } = useTransactions();

const currency = computed(() => wallet.value?.currency || 'USD');


const stats = computed(() => {
  const userWalletId = wallet.value?.id;

  if (!userWalletId) return { totalIncome: 0, totalExpense: 0, activeTransactions: 0 };

  let totalIncome = 0;
  let totalExpense = 0;

  allTransactions.value.forEach((tx: any) => {
    if (tx.Status === 'completed') {
      if (tx.ToWalletID === userWalletId) {
        totalIncome += tx.Amount;
      }
      if (tx.FromWalletID === userWalletId) {
        totalExpense += tx.Amount;
      }
    }
  });

  return {
    totalIncome,
    totalExpense,
    activeTransactions: lockedBalance.value,
  };
});

const formattedTransactions = computed(() => {
  const userWalletId = wallet.value?.ID;
  if (!userWalletId) return [];

  return allTransactions.value
    .filter((tx: any) => tx.FromWalletID === userWalletId || tx.ToWalletID === userWalletId)
    .map((tx: any) => {
      const isIncome = tx.ToWalletID === userWalletId;
      return {
        id: tx.ID,
        type: tx.Type,
        amount: isIncome ? tx.Amount : -tx.Amount,
        currency: currency.value,
        recipient: isIncome ? (tx.FromWalletID || 'System') : (tx.ToWalletID || 'System'),
        createdAt: new Date(tx.CreatedAt).toLocaleString(),
        status: tx.Status,
      };
    })
    .sort((a: any, b: any) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
});

const loading = computed(() => walletLoading.value || transactionsLoading.value);
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
          <TransactionTable :transactions="formattedTransactions" :loading="loading" />
        </div>
      </div>
      
      <div class="col-span-3 rounded-xl border bg-card text-card-foreground shadow">
        <div class="p-6">
          <h3 class="text-lg font-semibold mb-4">Quick Actions</h3>
          <div class="grid grid-cols-2 gap-4">
            <BaseButton 
              variant="outline" 
              class="h-24 flex-col gap-2"
              @click="router.push('/wallet/transfer')"
            >
              <BaseIcon name="Send" class="h-6 w-6" />
              Transfer
            </BaseButton>
            <BaseButton variant="outline" class="h-24 flex-col gap-2" @click="router.push('/wallet/topup')">
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
