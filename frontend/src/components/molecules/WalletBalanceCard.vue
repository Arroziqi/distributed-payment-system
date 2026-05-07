<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseSkeleton from '@/components/atoms/BaseSkeleton.vue';

interface Props {
  balance?: number;
  currency?: string;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  balance: 0,
  currency: 'USD',
});
</script>

<template>
  <Card class="overflow-hidden">
    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
      <CardTitle class="text-sm font-medium">Total Balance</CardTitle>
      <BaseIcon name="Wallet" class="h-4 w-4 text-muted-foreground" />
    </CardHeader>
    <CardContent>
      <div v-if="props.loading">
        <BaseSkeleton class="h-8 w-[100px]" />
      </div>
      <div v-else class="text-2xl font-bold">
        {{ props.currency }} {{ props.balance.toLocaleString() }}
      </div>
      <p class="text-xs text-muted-foreground">
        Available funds in your wallet
      </p>
    </CardContent>
  </Card>
</template>
