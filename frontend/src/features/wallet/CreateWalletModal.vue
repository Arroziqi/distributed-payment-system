<script setup lang="ts">
import { ref } from 'vue';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseInput from '@/components/atoms/BaseInput.vue';
import BaseLabel from '@/components/atoms/BaseLabel.vue';
import { useWallet } from '@/features/wallet/useWallet';
import { toast } from 'vue-sonner';

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  'update:open': [value: boolean];
  'success': [];
}>();

const { createWallet, isCreateWalletLoading } = useWallet();
const currency = ref('USD');
const error = ref('');

const handleOpenChange = (value: boolean) => {
  emit('update:open', value);
};

const handleCreate = async () => {
  error.value = '';
  try {
    await createWallet(currency.value);
    emit('success');
    handleOpenChange(false);
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Failed to create wallet';
  }
};
</script>

<template>
  <Dialog :open="props.open" @update:open="handleOpenChange">
    <DialogContent class="sm:max-w-[425px] bg-background/95 backdrop-blur-sm border-primary/20">
      <DialogHeader>
        <DialogTitle class="text-2xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
          Create Your Wallet
        </DialogTitle>
        <DialogDescription class="text-muted-foreground">
          You don't have a wallet yet. Create one now to start managing your funds and making transactions.
        </DialogDescription>
      </DialogHeader>
      
      <div class="grid gap-6 py-6">
        <div class="space-y-2">
          <BaseLabel for="currency" class="text-sm font-semibold tracking-wide uppercase">
            Default Currency
          </BaseLabel>
          <div class="relative">
            <BaseInput 
              id="currency" 
              v-model="currency" 
              placeholder="e.g. USD, EUR, IDR"
              class="pl-10 h-12 border-primary/10 focus:border-primary/30 transition-all"
            />
            <div class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
              💰
            </div>
          </div>
          <p class="text-[10px] text-muted-foreground italic">
            This will be the primary currency for your transactions.
          </p>
        </div>

        <div v-if="error" class="p-4 rounded-xl bg-destructive/10 text-destructive text-xs font-medium border border-destructive/20 animate-in fade-in slide-in-from-top-1">
          {{ error }}
        </div>
      </div>

      <DialogFooter>
        <BaseButton 
          type="button" 
          class="w-full h-12 text-md font-bold shadow-lg shadow-primary/20 hover:shadow-primary/30 transition-all transform hover:-translate-y-0.5 active:translate-y-0"
          :disabled="isCreateWalletLoading"
          @click="handleCreate"
        >
          <span v-if="isCreateWalletLoading" class="animate-spin mr-2">⏳</span>
          {{ isCreateWalletLoading ? 'Initializing...' : 'Create Wallet Now' }}
        </BaseButton>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
