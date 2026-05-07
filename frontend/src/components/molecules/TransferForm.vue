<script setup lang="ts">
import { reactive } from 'vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseInput from '@/components/atoms/BaseInput.vue';
import BaseLabel from '@/components/atoms/BaseLabel.vue';
import BaseAlert from '@/components/atoms/BaseAlert.vue';

interface Props {
  loading?: boolean;
  error?: string | null;
  maxAmount?: number;
}

defineProps<Props>();
const emit = defineEmits(['submit']);

const form = reactive({
  to_user_id: '',
  amount: 0,
  description: '',
});

const handleSubmit = () => {
  emit('submit', { ...form });
};
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <BaseAlert v-if="error" variant="destructive" :description="error" />

    <div class="space-y-2">
      <BaseLabel for="to_user_id">Recipient User ID</BaseLabel>
      <BaseInput
        id="to_user_id"
        v-model="form.to_user_id"
        placeholder="User ID or Email"
        required
        :disabled="loading"
      />
    </div>
    
    <div class="space-y-2">
      <BaseLabel for="amount">Amount</BaseLabel>
      <BaseInput
        id="amount"
        v-model.number="form.amount"
        type="number"
        placeholder="0.00"
        required
        :disabled="loading"
        :max="maxAmount"
      />
      <p v-if="maxAmount" class="text-xs text-muted-foreground">
        Available: {{ maxAmount.toLocaleString() }}
      </p>
    </div>

    <div class="space-y-2">
      <BaseLabel for="description">Description (Optional)</BaseLabel>
      <BaseInput
        id="description"
        v-model="form.description"
        placeholder="Lunch, etc."
        :disabled="loading"
      />
    </div>

    <BaseButton type="submit" class="w-full" :loading="loading">
      Send Money
    </BaseButton>
  </form>
</template>
