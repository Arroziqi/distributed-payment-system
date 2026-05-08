<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useWallet } from '@/features/wallet/useWallet';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseInput from '@/components/atoms/BaseInput.vue';
import { toast } from 'vue-sonner';

const router = useRouter();
const { createWallet, isCreateWalletLoading } = useWallet();

const currency = ref('');
const error = ref('');

const handleCreate = async () => {
  error.value = '';
  try {
    await createWallet(currency.value || undefined);
    router.push('/wallet');
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Failed to create wallet';
    toast.error(error.value);
  }
};
</script>

<template>
  <div class="max-w-lg mx-auto py-8 space-y-6">
    <h2 class="text-3xl font-bold">Create Wallet</h2>
    <p class="text-muted-foreground">Create a new wallet for your account. Currency is optional.</p>
    <div class="space-y-2">
      <label class="text-sm font-medium">Currency (e.g., USD)</label>
      <BaseInput v-model="currency" placeholder="Enter currency code" />
    </div>
    <div v-if="error" class="p-3 rounded-md bg-destructive/10 text-destructive text-sm font-medium border border-destructive/20">
      {{ error }}
    </div>
    <BaseButton :disabled="isCreateWalletLoading" @click="handleCreate" class="w-full">
      <span v-if="isCreateWalletLoading" class="animate-spin mr-2">⏳</span>
      {{ isCreateWalletLoading ? 'Creating...' : 'Create Wallet' }}
    </BaseButton>
  </div>
</template>
