<script setup lang="ts">
import { reactive } from 'vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseInput from '@/components/atoms/BaseInput.vue';
import BaseLabel from '@/components/atoms/BaseLabel.vue';
import BaseAlert from '@/components/atoms/BaseAlert.vue';

interface Props {
  loading?: boolean;
  error?: string | null;
}

defineProps<Props>();
const emit = defineEmits(['submit']);

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
});

const handleSubmit = () => {
  emit('submit', { ...form });
};
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <BaseAlert v-if="error" variant="destructive" :description="error" />

    <div class="space-y-2">
      <BaseLabel for="name">Full Name</BaseLabel>
      <BaseInput
        id="name"
        v-model="form.name"
        placeholder="John Doe"
        required
        :disabled="loading"
      />
    </div>
    
    <div class="space-y-2">
      <BaseLabel for="email">Email</BaseLabel>
      <BaseInput
        id="email"
        v-model="form.email"
        type="email"
        placeholder="name@example.com"
        required
        :disabled="loading"
      />
    </div>

    <div class="space-y-2">
      <BaseLabel for="password">Password</BaseLabel>
      <BaseInput
        id="password"
        v-model="form.password"
        type="password"
        required
        :disabled="loading"
      />
    </div>

    <div class="space-y-2">
      <BaseLabel for="confirmPassword">Confirm Password</BaseLabel>
      <BaseInput
        id="confirmPassword"
        v-model="form.confirmPassword"
        type="password"
        required
        :disabled="loading"
      />
    </div>

    <BaseButton type="submit" class="w-full" :loading="loading">
      Create Account
    </BaseButton>
  </form>
</template>
