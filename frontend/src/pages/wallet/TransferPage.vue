<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useWallet } from '@/features/wallet/useWallet';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseInput from '@/components/atoms/BaseInput.vue';

const router = useRouter();
const { transfer, isTransferLoading, balance, wallet } = useWallet();

const form = ref({
  toUserId: '',
  amount: 0,
});

const error = ref('');

const currency = computed(() => wallet.value?.currency || 'USD');

const isFormValid = computed(() => {
  return form.value.toUserId.trim() !== '' && form.value.amount > 0 && form.value.amount <= balance.value;
});

const handleTransfer = async () => {
  if (!isFormValid.value) return;
  
  error.value = '';
  try {
    await transfer(form.value.toUserId, form.value.amount);
    router.push('/');
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to complete transfer';
  }
};
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6 py-8">
    <div class="flex items-center space-x-4 mb-8">
      <BaseButton variant="ghost" size="icon" @click="router.back()">
        <BaseIcon name="ChevronLeft" class="h-6 w-6" />
      </BaseButton>
      <div>
        <h2 class="text-3xl font-bold tracking-tight">Transfer Funds</h2>
        <p class="text-muted-foreground">Send money to another user instantly.</p>
      </div>
    </div>

    <div class="rounded-xl border bg-card text-card-foreground shadow-sm overflow-hidden">
      <div class="p-6 bg-primary/5 border-b">
        <div class="flex justify-between items-center">
          <span class="text-sm font-medium text-muted-foreground">Available Balance</span>
          <span class="text-2xl font-bold text-primary">{{ balance }} {{ currency }}</span>
        </div>
      </div>
      
      <form @submit.prevent="handleTransfer" class="p-6 space-y-6">
        <div class="space-y-2">
          <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
            Recipient User ID
          </label>
          <BaseInput 
            v-model="form.toUserId" 
            placeholder="Enter recipient's user ID"
            required
          />
          <p class="text-[0.8rem] text-muted-foreground">
            The unique identifier of the user you want to send money to.
          </p>
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
            Amount ({{ currency }})
          </label>
          <BaseInput 
            v-model.number="form.amount" 
            type="number" 
            min="1"
            :max="balance"
            placeholder="0.00"
            required
          />
          <div v-if="form.amount > balance" class="text-[0.8rem] text-destructive font-medium">
            Insufficient balance.
          </div>
        </div>

        <div v-if="error" class="p-3 rounded-md bg-destructive/10 text-destructive text-sm font-medium border border-destructive/20">
          {{ error }}
        </div>

        <div class="pt-4 flex flex-col gap-3">
          <BaseButton 
            type="submit" 
            class="w-full h-12 text-lg"
            :disabled="!isFormValid || isTransferLoading"
          >
            <BaseIcon v-if="isTransferLoading" name="Loader2" class="mr-2 h-5 w-5 animate-spin" />
            <BaseIcon v-else name="Send" class="mr-2 h-5 w-5" />
            {{ isTransferLoading ? 'Processing...' : 'Send Transfer' }}
          </BaseButton>
          
          <BaseButton 
            type="button" 
            variant="outline" 
            class="w-full h-12"
            @click="router.back()"
          >
            Cancel
          </BaseButton>
        </div>
      </form>
    </div>

    <div class="rounded-xl border border-blue-200 bg-blue-50 p-6 text-blue-900">
      <div class="flex items-start gap-4">
        <BaseIcon name="Info" class="h-6 w-6 text-blue-600 mt-0.5" />
        <div class="space-y-1">
          <h4 class="font-semibold">Security Note</h4>
          <p class="text-sm opacity-90">
            Transfers are processed instantly and cannot be reversed. Please double-check the Recipient User ID before sending.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
