<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useWallet } from '@/features/wallet/useWallet';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseInput from '@/components/atoms/BaseInput.vue';

const router = useRouter();
const { topup, isTopupLoading, balance, wallet } = useWallet();

const amount = ref<number>(0);
const error = ref('');

const currency = computed(() => wallet.value?.currency || 'USD');

const PRESET_AMOUNTS = [50, 100, 250, 500];

const isFormValid = computed(() => amount.value > 0);

const selectPreset = (preset: number) => {
  amount.value = preset;
};

const handleTopUp = async () => {
  if (!isFormValid.value) return;

  error.value = '';
  try {
    await topup(amount.value);
    router.push('/');
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to complete top-up';
  }
};
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6 py-8">
    <!-- Header -->
    <div class="flex items-center space-x-4 mb-8">
      <BaseButton variant="ghost" size="icon" @click="router.back()">
        <BaseIcon name="ChevronLeft" class="h-6 w-6" />
      </BaseButton>
      <div>
        <h2 class="text-3xl font-bold tracking-tight">Top Up Wallet</h2>
        <p class="text-muted-foreground">Add funds to your wallet instantly.</p>
      </div>
    </div>

    <!-- Current Balance Banner -->
    <div class="rounded-xl border bg-card text-card-foreground shadow-sm overflow-hidden">
      <div class="p-6 bg-primary/5 border-b">
        <div class="flex justify-between items-center">
          <span class="text-sm font-medium text-muted-foreground">Current Balance</span>
          <span class="text-2xl font-bold text-primary">{{ balance }} {{ currency }}</span>
        </div>
      </div>

      <form @submit.prevent="handleTopUp" class="p-6 space-y-6">
        <!-- Preset Quick-Select -->
        <div class="space-y-2">
          <label class="text-sm font-medium leading-none">Quick Select Amount</label>
          <div class="grid grid-cols-4 gap-2">
            <button
              v-for="preset in PRESET_AMOUNTS"
              :key="preset"
              type="button"
              :class="[
                'rounded-lg border px-4 py-2 text-sm font-semibold transition-colors',
                amount === preset
                  ? 'bg-primary text-primary-foreground border-primary'
                  : 'bg-background text-foreground border-border hover:bg-accent hover:text-accent-foreground'
              ]"
              @click="selectPreset(preset)"
            >
              {{ preset }}
            </button>
          </div>
        </div>

        <!-- Custom Amount Input -->
        <div class="space-y-2">
          <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
            Custom Amount ({{ currency }})
          </label>
          <BaseInput
            v-model.number="amount"
            type="number"
            min="1"
            placeholder="0.00"
            required
          />
          <p class="text-[0.8rem] text-muted-foreground">
            Enter a custom amount or use the quick-select options above.
          </p>
        </div>

        <!-- Error Message -->
        <div
          v-if="error"
          class="p-3 rounded-md bg-destructive/10 text-destructive text-sm font-medium border border-destructive/20"
        >
          {{ error }}
        </div>

        <!-- Actions -->
        <div class="pt-4 flex flex-col gap-3">
          <BaseButton
            type="submit"
            class="w-full h-12 text-lg"
            :disabled="!isFormValid || isTopupLoading"
          >
            <BaseIcon v-if="isTopupLoading" name="Loader2" class="mr-2 h-5 w-5 animate-spin" />
            <BaseIcon v-else name="PlusCircle" class="mr-2 h-5 w-5" />
            {{ isTopupLoading ? 'Processing...' : `Top Up ${amount > 0 ? amount + ' ' + currency : ''}` }}
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

    <!-- Info Note -->
    <div class="rounded-xl border border-emerald-200 bg-emerald-50 p-6 text-emerald-900">
      <div class="flex items-start gap-4">
        <BaseIcon name="Info" class="h-6 w-6 text-emerald-600 mt-0.5" />
        <div class="space-y-1">
          <h4 class="font-semibold">How Top-Up Works</h4>
          <p class="text-sm opacity-90">
            Funds are added to your available balance immediately and can be used for transfers right away.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
