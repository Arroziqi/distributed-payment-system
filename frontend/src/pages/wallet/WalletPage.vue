<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useWallet } from '@/features/wallet/useWallet';
import WalletBalanceCard from '@/components/molecules/WalletBalanceCard.vue';
import CreateWalletModal from '@/features/wallet/CreateWalletModal.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { useRouter } from 'vue-router';

const router = useRouter();
const { balance, lockedBalance, currency, loading, isError, error, refetch } = useWallet();
const isCreateModalOpen = ref(false);

watch([isError, error, loading], () => {
  if (isError.value && !loading.value) {
    const err = error.value as any;
    if (err?.response?.status === 404 || err?.status === 404) {
      isCreateModalOpen.value = true;
    }
  }
}, { immediate: true });

const handleSuccess = () => {
  isCreateModalOpen.value = false;
  refetch();
};
</script>

<template>
  <div class="space-y-8 animate-in fade-in duration-500">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-4xl font-extrabold tracking-tight bg-gradient-to-br from-foreground to-foreground/60 bg-clip-text text-transparent">
          Wallet Overview
        </h2>
        <p class="text-muted-foreground mt-1 text-lg">
          Manage your assets and monitor your financial health.
        </p>
      </div>
      <div class="flex items-center gap-3">
        <BaseButton 
          variant="outline" 
          @click="router.push('/wallet/topup')"
          class="h-11 px-6 border-primary/20 hover:bg-primary/5 transition-all shadow-sm"
        >
          <BaseIcon name="Plus" class="mr-2 h-4 w-4" />
          Top Up
        </BaseButton>
        <BaseButton 
          @click="router.push('/wallet/transfer')"
          class="h-11 px-6 shadow-lg shadow-primary/20 transition-all transform hover:-translate-y-0.5"
        >
          <BaseIcon name="Send" class="mr-2 h-4 w-4" />
          Transfer
        </BaseButton>
      </div>
    </div>

    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <WalletBalanceCard 
        :balance="balance" 
        :currency="currency" 
        :loading="loading" 
        class="border-primary/10 shadow-md hover:shadow-lg transition-shadow"
      />
      
      <Card class="border-primary/10 shadow-md hover:shadow-lg transition-shadow overflow-hidden group">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Locked Balance</CardTitle>
          <BaseIcon name="Lock" class="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ currency }} {{ lockedBalance.toLocaleString() }}
          </div>
          <p class="text-xs text-muted-foreground mt-1">
            Funds currently in escrow or pending.
          </p>
        </CardContent>
      </Card>

      <Card class="border-primary/10 shadow-md hover:shadow-lg transition-shadow overflow-hidden group">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Quick Actions</CardTitle>
          <BaseIcon name="Zap" class="h-4 w-4 text-primary animate-pulse" />
        </CardHeader>
        <CardContent class="grid grid-cols-2 gap-2 mt-2">
          <BaseButton variant="secondary" size="sm" class="w-full text-[10px] h-8" @click="router.push('/transactions')">
            History
          </BaseButton>
          <BaseButton variant="secondary" size="sm" class="w-full text-[10px] h-8" @click="router.push('/profile')">
            Settings
          </BaseButton>
        </CardContent>
      </Card>
    </div>

    <Card class="border-primary/5 bg-primary/[0.02] backdrop-blur-sm overflow-hidden">
      <CardHeader>
        <CardTitle class="text-xl">Wallet Security</CardTitle>
        <CardDescription>
          Your funds are protected by our multi-layered security system.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-start gap-4 p-4 rounded-lg bg-background/50 border border-primary/5">
          <div class="p-2 rounded-full bg-primary/10 text-primary">
            <BaseIcon name="ShieldCheck" class="h-5 w-5" />
          </div>
          <div>
            <h4 class="font-semibold text-sm">Two-Factor Authentication</h4>
            <p class="text-xs text-muted-foreground">Ensure your account is secured with 2FA enabled in settings.</p>
          </div>
        </div>
      </CardContent>
    </Card>

    <CreateWalletModal 
      v-model:open="isCreateModalOpen" 
      @success="handleSuccess" 
    />
  </div>
</template>
