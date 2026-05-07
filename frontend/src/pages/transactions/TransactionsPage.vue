<script setup lang="ts">
import TransactionTable from '@/components/organisms/TransactionTable.vue';
import { useTransactions } from '@/features/transactions/useTransactions';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import SearchToolbar from '@/components/molecules/SearchToolbar.vue';
import { ref } from 'vue';

const search = ref('');
const { transactions, loading } = useTransactions();
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">Transactions</h2>
        <p class="text-muted-foreground">View and manage all your payments.</p>
      </div>
      <div class="flex items-center space-x-2">
        <BaseButton variant="outline">
          <BaseIcon name="Download" class="mr-2 h-4 w-4" />
          Export
        </BaseButton>
      </div>
    </div>

    <div class="flex items-center justify-between space-x-4">
      <SearchToolbar v-model="search" placeholder="Search transactions..." />
    </div>

    <TransactionTable :transactions="transactions" :loading="loading" />
  </div>
</template>
